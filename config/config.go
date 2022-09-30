package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env      string `default:"dev" envconfig:"ENV"`
	Port     string `default:"8080" envconfig:"PORT"`
	Version  string `default:"v1" envconfig:"VERSION"`
	LogLevel string `default:"debug" envconfig:"LOG_LEVEL"`

	Database Database
	Redis    Redis
	Kafka    Kafka
}

type Database struct {
	Type     string `default:"mysql" envconfig:"DATABASE_TYPE"`
	Host     string `default:"localhost" envconfig:"DATABASE_HOST"`
	Port     string `default:"3306" envconfig:"DATABASE_PORT"`
	Username string `default:"root" envconfig:"DATABASE_USER"`
	Password string `default:"root" envconfig:"DATABASE_PASSWORD"`
	DBName   string `default:"schema" envconfig:"DATABASE_NAME"`
	SSL      string `default:"DISABLED" envconfig:"DATABASE_SSL"`
}

type Redis struct {
	Host     string `default:"localhost" envconfig:"REDIS_HOST"`
	Port     string `default:"6379" envconfig:"REDIS_PORT"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
}

type Kafka struct {
	Host      string   `default:"localhost" envconfig:"KAFKA_HOST"`
	Port      string   `default:"9092" envconfig:"KAFKA_PORT"`
	TopicName []string `default:"test1,test2" envconfig:"KAFKA_TOPIC_NAME"`
}

func Loadenv(c *Config) error {
	godotenv.Load()
	if err := envconfig.Process("", c); err != nil {
		return err
	}
	return nil
}
