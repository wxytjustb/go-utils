package monitor

type SentryConfig struct {
	Dsn   string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	Debug bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
}
