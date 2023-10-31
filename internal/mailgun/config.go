package mailgun

type Config struct {
	Domain string `koanf:"domain"`
	APIKey string `koanf:"api_key"`
}
