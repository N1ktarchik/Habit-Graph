package redis

type redisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisConfig(host string, port int, password string, DB int) *redisConfig {
	return &redisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       DB,
	}
}
