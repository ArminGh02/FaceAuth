package imagga

type Config struct {
	APIKey    string `koanf:"api_key"`
	APISecret string `koanf:"api_secret"`
}
