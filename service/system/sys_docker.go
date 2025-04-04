package system

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spectacleCase/ci-cd-engine/global"
	system "github.com/spectacleCase/ci-cd-engine/models/system"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"strings"
)

// InitDockerCli 初始化 Docker 客户端
func InitDockerCli() {
	var cli *client.Client
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("无法创建 Docker 客户端: %v", err)
	}
	global.DockerCli = cli
}

// AssemblyLineProject 组装项目流水线
func AssemblyLineProject(build system.Stage, deploy system.Stage) {
	ctx := Inspect()

	// 1. 拉取 Python 3.9 镜像
	imageName := build.Image
	Pull(imageName, ctx)
	parts := strings.Split(build.Ports[0], ":")
	localPort := nat.Port(parts[0] + "/tcp")
	reflectionPort := parts[1]
	// 2. 创建容器（但不启动容器）
	conf := &container.Config{
		Image:      imageName,
		Cmd:        deploy.Commands,
		WorkingDir: build.WorkingDirectory,
		ExposedPorts: map[nat.Port]struct{}{
			localPort: {},
		},
	}

	// 配置主机绑定和端口映射
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			localPort: {
				{HostIP: "0.0.0.0", HostPort: reflectionPort},
			},
		},
	}

	// 3. 创建容器（但还不启动）
	resp := CreateContainer(conf, ctx, hostConfig, nil, nil, deploy.Image)

	// 4. 复制文件到容器
	CopyFilesToContainer(resp.ID, build.Volumes[0], build.Volumes[1])

	// 5. 启动容器
	go StartContainerWithFiles(conf, ctx, hostConfig, resp.ID)

	return

	// 6. 清理容器
	//ClearContainer(resp, ctx)
}

// AssemblyLinePythonProject 构建并运行 Python Web 服务
func AssemblyLinePythonProject() {
	ctx := Inspect()

	// 1. 拉取 Python 3.9 镜像
	imageName := "python:3.9"
	Pull(imageName, ctx)

	// 2. 创建容器（但不启动容器）
	conf := &container.Config{
		Image:      imageName,
		Cmd:        []string{"sh", "-c", "pip install --no-cache-dir -r /app/requirements.txt && python /app/main.py"},
		WorkingDir: "/app",
		ExposedPorts: map[nat.Port]struct{}{
			"5000/tcp": {},
		},
	}

	// 配置主机绑定和端口映射
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			"5000/tcp": {
				{HostIP: "0.0.0.0", HostPort: "5000"},
			},
		},
	}

	// 3. 创建容器（但还不启动）
	resp := CreateContainer(conf, ctx, hostConfig, nil, nil, "python-app")

	// 4. 复制文件到容器
	CopyFilesToContainer(resp.ID, "C:/Users/lizhentao/Desktop/python", "/app")

	// 5. 启动容器
	go StartContainerWithFiles(conf, ctx, hostConfig, resp.ID)

	return

	// 6. 清理容器
	//ClearContainer(resp, ctx)
}

// CreateContainer 创建并启动容器
func CreateContainer(conf *container.Config, ctx context.Context, hostConfig *container.HostConfig,
	networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) container.CreateResponse {
	log.Println("正在创建容器...")

	// 检查容器是否存在
	_, err := global.DockerCli.ContainerInspect(ctx, containerName)
	if err == nil {
		// 如果容器已存在，删除它
		log.Printf("容器 %s 已存在，正在删除...\n", containerName)
		err := global.DockerCli.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
		if err != nil {
			log.Fatalf("删除容器失败: %v", err)
		}
		log.Printf("容器 %s 删除成功\n", containerName)
	} else if client.IsErrNotFound(err) {
		// 如果容器不存在，忽略错误
		log.Printf("容器 %s 不存在，准备创建...\n", containerName)
	} else {
		log.Fatalf("检查容器时出错: %v", err)
	}

	// 创建容器
	resp, err := global.DockerCli.ContainerCreate(ctx, conf, hostConfig, networkingConfig, platform, containerName)
	if err != nil {
		log.Fatalf("创建容器失败: %v", err)
	}

	log.Printf("容器 %s 创建成功，准备启动...\n", containerName)

	return resp
}

// CopyFilesToContainer 复制本地文件到容器
func CopyFilesToContainer(containerID, srcPath, destPath string) {
	ctx := context.Background()
	cli := global.DockerCli

	// 读取本地文件并打包成 tar
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	files, err := os.ReadDir(srcPath)
	if err != nil {
		global.CLog.Error("读取目录失败", zap.Error(err))
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", srcPath, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			global.CLog.Error("无法读取文件信息", zap.Error(err))
			continue
		}

		hdr := &tar.Header{
			Name: file.Name(),
			Mode: 0600,
			Size: fileInfo.Size(),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			global.CLog.Error("写入 tar 头失败", zap.Error(err))
		}

		f, err := os.Open(filePath)
		if err != nil {
			global.CLog.Error("打开文件失败", zap.Error(err))
		}

		if _, err := io.Copy(tw, f); err != nil {
			global.CLog.Error("复制文件失败", zap.Error(err))
		}
		//f.Close()
		Close(f)
	}

	//tw.Close()
	Close(tw)

	// 复制到容器
	err = cli.CopyToContainer(ctx, containerID, destPath, &buf, container.CopyToContainerOptions{})
	if err != nil {
		global.CLog.Error("文件复制失败", zap.Error(err))
	}
	global.CLog.Info("文件成功复制到容器")
}

// StartContainerWithFiles 启动容器
func StartContainerWithFiles(conf *container.Config, ctx context.Context, hostConfig *container.HostConfig, containerID string) {
	// 启动容器
	if err := global.DockerCli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		global.CLog.Error("启动容器失败:", zap.Error(err))
	}
	global.CLog.Info("启动容器成功:", zap.String("Image", conf.Image))
	return
}

// Inspect 检查 Docker 初始化
func Inspect() context.Context {
	if global.DockerCli == nil {
		InitDockerCli()
	}
	ctx := context.Background()
	if _, err := global.DockerCli.Ping(ctx); err != nil {
		global.CLog.Error("无法连接 Docker 守护进程: %v\n请确保 Docker 服务正在运行", zap.Error(err))
	}
	return ctx
}

// Pull 拉取镜像
func Pull(imageName string, ctx context.Context) {
	reader, err := global.DockerCli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		global.CLog.Error("拉取镜像失败", zap.Error(err))
	}
	if _, err := io.Copy(os.Stdout, reader); err != nil {
		global.CLog.Error("读取镜像拉取输出时出错:", zap.Error(err))
	}
	Close(reader)
}

// ClearContainer 清理容器
func ClearContainer(resp container.CreateResponse, ctx context.Context) {
	if err := global.DockerCli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
		global.CLog.Error("删除容器失败:", zap.Error(err))
	} else {
		global.CLog.Info("容器已清除")
	}
}

func Close(cls interface{}) {
	switch cls := cls.(type) {
	case io.ReadCloser:
		err := cls.Close()
		if err != nil {
			return
		}
	case tar.Writer:
		err := cls.Close()
		if err != nil {
			return
		}
	case os.File:
		err := cls.Close()
		if err != nil {
			return
		}

	}
}
