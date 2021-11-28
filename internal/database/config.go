package database

// Config provide database configuration
type Config struct {
	Name                   string `env:"DB_NAME" json:",omitempty"`
	User                   string `env:"DB_USER" json:",omitempty"`
	Protocol               string `env:"DB_PROTOCOL, default=tcp" json:",omitempty"`
	Address                string `env:"DB_ADDRESS, default=localhost:3306" json:",omitempty"`
	Password               string `env:"DB_PASSWORD" json:"-"` // ignored by zap's JSON formatter
	EnableSSL              bool   `env:"DB_ENABLE_SSL, default=false" json:",omitempty"`
	SSLCertPath            string `env:"DB_SSLCERT" json:",omitempty"`
	SSLKeyPath             string `env:"DB_SSLKEY" json:",omitempty"`
	SSLRootCertPath        string `env:"DB_SSLROOTCERT" json:",omitempty"`
	ConnectionTimeout      int    `env:"DB_CONNECT_TIMEOUT, default=90" json:",omitempty"`
	PoolMaxIdleConnections int    `env:"DB_POOL_MAX_IDLE_CONNS, default=5" json:",omitempty"`
	PoolMaxConnections     int    `env:"DB_POOL_MAX_CONNS, default=30" json:",omitempty"`
}

// DatabaseConfig get db config
func (c *Config) DatabaseConfig() *Config {
	return c
}
