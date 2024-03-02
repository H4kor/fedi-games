package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type Config struct {
	Host         string
	Protocol     string
	PublicKeyPem string
	PrivKey      *rsa.PrivateKey
}

var cfg *Config

func (c *Config) FullUrl() string {
	return c.Protocol + "://" + c.Host
}

func getEnv(key string, fb string) string {
	v := os.Getenv(key)
	if v == "" {
		return fb
	}
	return v
}

func GetConfig() Config {
	if cfg == nil {
		privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		pubKey := privKey.Public().(*rsa.PublicKey)

		pubKeyPem := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(pubKey),
			},
		)

		cfg = &Config{
			Host:         getEnv("FEDI_GAMES_HOST", "localhost:4040"),
			Protocol:     getEnv("FEDI_GAMES_PROTOCOL", "http"),
			PublicKeyPem: string(pubKeyPem),
			PrivKey:      privKey,
		}
	}

	return *cfg
}
