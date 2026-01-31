package internal

import (
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func LoadScraperConfig(cfg map[string]any, out any) error {
	err := defaults.Set(out)
	if err != nil {
		return err
	}

	err = mapstructure.Decode(cfg, out)
	if err != nil {
		return err
	}

	err = validate.Struct(out)
	if err != nil {
		return err
	}

	return nil
}
