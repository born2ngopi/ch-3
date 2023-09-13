package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

type (
	Config struct {
		Database Database
		Redis    Redis
	}

	// configurasi database prostgresql
	Database struct {
		Host     string
		Port     string
		User     string
		DBName   string
		SSLMode  string
		Password string
	}

	Redis struct {
		Addr     string
		Db       int
		Password string
	}
)

func Init() *Config {
	dbRedis, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbRedis = 0
	}
	return &Config{
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Redis: Redis{
			Addr:     os.Getenv("REDIS_ADDR"),
			Db:       dbRedis,
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
}

func (c *Config) ConnectDatabase() *gorm.DB {

	dsn := "host=" + c.Database.Host + " user=" + c.Database.User + " password=" + c.Database.Password + " dbname=" + c.Database.DBName + " port=" + c.Database.Port + " sslmode=" + c.Database.SSLMode + " TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func (c *Config) ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	})

	// check connection redis
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}

	return rdb
}
