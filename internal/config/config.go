package config

import (
	"fmt"
	"strings"

	"github.com/ArminGh02/go-auth-system/internal/imagga"
	"github.com/ArminGh02/go-auth-system/internal/mailgun"
	"github.com/ArminGh02/go-auth-system/internal/s3"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

const Prefix = "AUTH_"

type Config struct {
	ServerPort string         `koanf:"server_port"`
	Database   string         `koanf:"db"`
	Broker     string         `koanf:"broker"`
	S3         s3.Config      `koanf:"s3"`
	Imagga     imagga.Config  `koanf:"imagga"`
	MailGun    mailgun.Config `koanf:"mailgun"`
}

func New() (*Config, error) {
	k := koanf.New(".")

	err := k.Load(
		env.Provider(Prefix, ".", func(s string) string {
			s = strings.ToLower(strings.TrimPrefix(s, Prefix))
			switch {
			case strings.HasPrefix(s, "s3_"),
				strings.HasPrefix(s, "imagga_"),
				strings.HasPrefix(s, "mailgun_"):
				return strings.Replace(s, "_", ".", 1)
			default:
				return s
			}
		}),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return &cfg, nil
}
