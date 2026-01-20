package authn

import (
	"fmt"

	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
)

func ValidateUsername(name string) error {
	if name == "" {
		return fmt.Errorf("empty string")
	}

	if !utils.ValidatePrintableAscii(name) {
		return fmt.Errorf("%s contains non-ascii characters", name)
	}

	return nil
}
