package s3

type Config struct {
	Endpoint  string `koanf:"endpoint"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
	Region    string `koanf:"region"`
	Bucket    string `koanf:"bucket"`
}
