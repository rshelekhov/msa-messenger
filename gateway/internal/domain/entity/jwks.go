package entity

import "time"

type JWKS struct {
	Keys []JWK
	TTL  time.Duration
}

type JWK struct {
	Alg string `json:"alg,omitempty"` // The specific cryptographic algorithm used with the key
	Kty string `json:"kty,omitempty"` // The family of cryptographic algorithms used with the key
	Use string `json:"use,omitempty"` // How the key was meant to be used; sig represents the signature
	Kid string `json:"kid,omitempty"` // The unique identifier for the key

	// For RSA keys
	N string `json:"n,omitempty"` // The modulus for the RSA public key
	E string `json:"e,omitempty"` // The exponent for the RSA public key
}
