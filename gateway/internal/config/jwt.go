package config

type JWT struct {
	JWKSEndpoint string `mapstructure:"JWKS_ENDPOINT"`
}
