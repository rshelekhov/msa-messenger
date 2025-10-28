package config

type (
	GRPCServices struct {
		SSOService        SSOService        `yaml:"SSOService"`
		ChatService       ChatService       `yaml:"ChatService"`
		SubscriberService SubscriberService `yaml:"SubscriberService"`
	}

	SSOService struct {
		Address            string `yaml:"Address"`
		VerificationURL    string `yaml:"VerificationURL"`
		ConfirmPasswordURL string `yaml:"ConfirmPasswordURL"`
	}

	ChatService struct {
		Address string `yaml:"Address"`
	}

	SubscriberService struct {
		Address string `yaml:"Address"`
	}
)
