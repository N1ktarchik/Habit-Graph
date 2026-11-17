package redis

type redisConfig struct {
	Host     string
	Port     int
	Password string
	db       int
}

func NewRedisConfig(host string, port int, password string, db int) *redisConfig {
	return &redisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		db:       db,
	}
}
