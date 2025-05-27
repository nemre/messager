package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config interface {
	GetServer() Server
	GetPostgreSQL() PostgreSQL
	GetRedis() Redis
	GetJob() Job
	GetKafka() Kafka
	GetClient() Client
}

type Server struct {
	Host     string `env:"HOST,required,notEmpty"`
	Port     uint16 `env:"PORT,required,notEmpty"`
	IDHeader string `env:"ID_HEADER,required,notEmpty"`
}

type PostgreSQL struct {
	Host     string `env:"HOST,required,notEmpty"`
	Port     uint16 `env:"PORT,required,notEmpty"`
	User     string `env:"USER,required,notEmpty"`
	Password string `env:"PASSWORD,required,notEmpty"`
	Name     string `env:"NAME,required,notEmpty"`
}

type Redis struct {
	Host     string `env:"HOST,required,notEmpty"`
	Port     uint16 `env:"PORT,required,notEmpty"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	DB       uint16 `env:"DB,required,notEmpty"`
}

type Job struct {
	Interval time.Duration `env:"INTERVAL,required,notEmpty"`
}

type Kafka struct {
	Brokers []string `env:"BROKERS,required,notEmpty"`
	Topic   string   `env:"TOPIC,required,notEmpty"`
	GroupID string   `env:"GROUP_ID,required,notEmpty"`
}

type Client struct {
	URL     string        `env:"URL,required,notEmpty"`
	Token   string        `env:"TOKEN,required,notEmpty"`
	Timeout time.Duration `env:"TIMEOUT,required,notEmpty"`
}

type config struct {
	Server     Server     `envPrefix:"SERVER_"`
	PostgreSQL PostgreSQL `envPrefix:"POSTGRESQL_"`
	Redis      Redis      `envPrefix:"REDIS_"`
	Job        Job        `envPrefix:"JOB_"`
	Kafka      Kafka      `envPrefix:"KAFKA_"`
	Client     Client     `envPrefix:"CLIENT_"`
}

func New() (Config, error) {
	var cfg config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("env.Parse(): %w", err)
	}

	return &cfg, nil
}

func (c *config) GetServer() Server {
	return c.Server
}

func (c *config) GetPostgreSQL() PostgreSQL {
	return c.PostgreSQL
}

func (c *config) GetRedis() Redis {
	return c.Redis
}

func (c *config) GetJob() Job {
	return c.Job
}

func (c *config) GetKafka() Kafka {
	return c.Kafka
}

func (c *config) GetClient() Client {
	return c.Client
}
