package server

type Config struct {
	ProviderName   string
	HtpasswdFile   string
	AllowedIPs     []string
	AllowedDomains []string
	AccessLogFile  string
	Port           int
}

func NewDefaultConfig() *Config {
	return &Config{}
}
