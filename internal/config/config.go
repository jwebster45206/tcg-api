package config

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Config struct {
	Env      string       `json:"env"`
	Port     string       `json:"port"`
	DB       MySQLConfig  `json:"db"`
	DBReader MySQLConfig  `json:"db_reader"`
	Redis    RedisConfig  `json:"redis"`
	Logger   LoggerConfig `json:"logger"`
}
