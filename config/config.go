package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

var DB *gorm.DB
var RedisCli *redis.Client
var MilvusCli *milvusclient.Client

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectMySQL()
	connectRedis()
	connectMilvus()
}

// MySQL
func connectMySQL() *gorm.DB {
	// 从环境变量读取配置
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDB := os.Getenv("MYSQL_DB")

	// 生成 MySQL DSN 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)

	// 连接 MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}

	fmt.Println("✅ MySQL 连接成功！")

	// 配置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接失败")
	}
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数（默认无限制）
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数（默认 2）
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大存活时间（默认无限制）

	DB = db
	return db
}

// Redis
func connectRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	// 测试 Redis 连接
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil
	}
	RedisCli = client

	return client
}

func connectMilvus() *milvusclient.Client {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	milvusAddr := os.Getenv("MILVUS_HOST")

	cli, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: milvusAddr,
	})

	if err != nil {
		fmt.Println(err)
		// handle error
	}
	MilvusCli = cli
	return cli
}
