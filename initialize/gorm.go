package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitMySQL() {
	mConfig := config.Config.Mysql

	// 构造 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		mConfig.UserName,
		mConfig.Password,
		mConfig.DbHost,
		mConfig.DbPort,
		mConfig.DbName,
		mConfig.Charset,
	)

	// 设置 GORM 日志级别
	var ormLogger logger.Interface
	if gin.Mode() == gin.DebugMode {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default.LogMode(logger.Warn)
	}

	// 初始化数据库连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("连接 MySQL 失败：" + err.Error())
	}

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接池失败：" + err.Error())
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// 设置建表选项（如果你使用 AutoMigrate）
	db = db.Set("gorm:table_options", "charset=utf8mb4")

	// 保存全局 DB 实例
	global.CDB = db
}
