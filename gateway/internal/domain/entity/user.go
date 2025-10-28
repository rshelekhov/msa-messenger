package entity

import "time"

type (
	User struct {
		ID              string
		Name            string
		Email           string
		CurrentPassword string
		UpdatedPassword string
	}

	UserCredentials struct {
		Email    string
		Password string
		Name     string
	}

	UserDevice struct {
		UserAgent string
		IP        string
	}

	UserTokens struct {
		AccessToken      string
		RefreshToken     string
		Domain           string
		Path             string
		ExpiresAt        time.Time
		HTTPOnly         bool
		AdditionalFields map[string]string
	}
)
