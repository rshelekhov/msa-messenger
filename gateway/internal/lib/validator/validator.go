package validator

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/rshelekhov/msa-messenger/gateway/internal/config"
)

// Password must contain at least 8 characters and can only include Latin letters (A-Z, a-z),
// digits (0-9), and special characters (@, $, !, %, *, ?, &)
const defaultPasswordRegex = "^[A-Za-z\\d@$!%*?&]{8,}$"

func New(cfg *config.Validator) (*validator.Validate, error) {
	const op = "validator.New"

	passwordRegex := cfg.PasswordRegexp
	if passwordRegex == "" {
		passwordRegex = defaultPasswordRegex
	}

	validate := validator.New()

	if err := validate.RegisterValidation("password_regexp", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(passwordRegex)
		return re.MatchString(fl.Field().String())
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return validate, nil
}
