package config

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Env  string
	Port string
	DB   MySQLConfig
}
