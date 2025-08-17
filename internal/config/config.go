package config

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type AuthConfig struct {
	Enabled          bool              `json:"enabled"`
	Issuer           string            `json:"issuer"`
	Audience         string            `json:"audience"`
	Algorithm        string            `json:"algorithm"`
	AccessTTLSeconds int               `json:"access_ttl_seconds"`
	CurrentKID       string            `json:"current_kid"`
	Keys             map[string]string `json:"keys"`
}

type Config struct {
	Env      string       `json:"env"`
	Port     string       `json:"port"`
	DB       MySQLConfig  `json:"db"`
	DBReader MySQLConfig  `json:"db_reader"`
	Redis    RedisConfig  `json:"redis"`
	Logger   LoggerConfig `json:"logger"`
	Auth     AuthConfig   `json:"auth"`
}
