package config

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Env    string       `json:"env"`
	Port   string       `json:"port"`
	DB     MySQLConfig  `json:"db"`
	Logger LoggerConfig `json:"logger"`
}
