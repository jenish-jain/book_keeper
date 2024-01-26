package mongo

type Config struct {
	host         string
	rwUsername   string
	rwPassword   string
	databaseName string
}

func NewConfig(host string, rwUsername string, rwPassword string, databaseName string) *Config {
	return &Config{host: host,
		rwPassword:   rwPassword,
		rwUsername:   rwUsername,
		databaseName: databaseName,
	}
}
