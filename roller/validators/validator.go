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
	reInvalid := regexp.MustCompile(`[^0-9d\+\-\*/! ]`)
	if reInvalid.MatchString(rollString) {
		return errors.New("Invalid characters")
	}

	//check for trailing operands
	reTrailing := regexp.MustCompile(`[\+\-\*/]$`)
	if reTrailing.MatchString(rollString) {
		return errors.New("Trailing operand")
	}

	//check for preceding operands
	rePreceding := regexp.MustCompile(`^[\+\-\*/!]`)
	if rePreceding.MatchString(rollString) {
		return errors.New("Preceding operand")
	}

	//check for invalid consecutive operands
	reConsecutive := regexp.MustCompile(`[\+\-\*/]{2,}`)
	if reConsecutive.MatchString(rollString) {
		return errors.New("Consecutive operands")
	}

	//check for consecutive d's
	reConsecutiveD := regexp.MustCompile(`dd`)
	if reConsecutiveD.MatchString(rollString) {
		return errors.New("Consecutive d's")
	}

	//check for trailing d's
	reTrailingD := regexp.MustCompile(`d$`)
	if reTrailingD.MatchString(rollString) {
		return errors.New("Trailing d")
	}

	//check for preceding d's
	rePrecedingD := regexp.MustCompile(`^d`)
	if rePrecedingD.MatchString(rollString) {
		return errors.New("Preceding d")
	}

	//check for d's with following or preceding operands
	reDOperand := regexp.MustCompile(`[\+\-\*/!]d[\+\-\*/!]`)
	if reDOperand.MatchString(rollString) {
		return errors.New("d with following or preceding operand")
	}

	return nil
}
