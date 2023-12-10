package validators

import (
	"errors"
	"regexp"
)

func RollValidator(rollString string) error {
	//check for empty string
	if rollString == "" {
		return errors.New("Empty string")
	}

	//check for invalid characters
	reInvalid := regexp.MustCompile(`[^0-9d+-/*!]`)
	if reInvalid.MatchString(rollString) {
		return errors.New("Invalid characters")
	}

	//check for trailing operands
	reTrailing := regexp.MustCompile(`[+-/*]$`)
	if reTrailing.MatchString(rollString) {
		return errors.New("Trailing operand")
	}

	//check for preceding operands
	rePreceding := regexp.MustCompile(`^[+-/*!]`)
	if rePreceding.MatchString(rollString) {
		return errors.New("Preceding operand")
	}

	//check for invalid consecutive operands
	reConsecutive := regexp.MustCompile(`[+-/*!]{2,}`)
	if reConsecutive.MatchString(rollString) {
		return errors.New("Consecutive operands")
	}

	return nil
}
