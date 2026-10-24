package envs

type PostgresConfig struct {
	PostgresDB       string `env:"POSTGRES_DB,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresSSLMode  string `env:"POSTGRES_SSL_MODE,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
}
