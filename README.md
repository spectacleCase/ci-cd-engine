# 轻量级 CI/CD 引擎

## **简介**
本项目是一个轻量级 CI/CD 系统，类似 Drone，但更简单，支持：
- **基于 YAML 定义 Pipeline**
- **Docker/Kubernetes 任务执行**
- **实时日志流式输出**
- **插件系统（如 Slack 通知）**
- **超快启动（相比 Jenkins 更轻量）**

---

## **技术栈**
- **后端**：Go（调度引擎）
- **容器执行**：Docker API / Kubernetes Client
- **数据库**：SQLite / BoltDB（存储任务状态）
- **前端（可选）**：Vue / React（展示任务执行状态）

---

## **快速开始**
1. **安装 Docker**
   https://docs.docker.com/get-docker/
2. **运行 CI/CD 引擎**
```sh
    go run main.go
```
3. **配置 YAML 在仓库根目录创建 .cicd.yaml**:
```shell
version: "1.0"
stages:
  - name: build
    image: golang:1.20
    commands:
      - go build -o myapp
  - name: test
    depends_on: build
    image: golang:1.20
    commands:
      - go test ./...

```