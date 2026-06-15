package fetcher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateFSAConfigString validates a config token like "[A-Z]" or "[0-9]".
func ValidateFSAConfigString(configStr string) error {
	regexPattern := `^\[[A-Za-z0-9\-]+\]$`
	if !regexp.MustCompile(regexPattern).MatchString(configStr) {
		return fmt.Errorf("Error: Invalid FSA config string: %s", configStr)
	}
	return nil
}

// ValidatePostalCodeFSA validates the first 3 characters (FSA) of a Canadian postal code.
// Returns (isValid, regexString, error)
func ValidatePostalCodeFSA(postalCodeFSA string) (bool, string, error) {
	if len(postalCodeFSA) == 0 {
		return false, "", errors.New("Error: Postal Code FSA is Empty")
	}
	if len(postalCodeFSA) != 3 {
		return false, "", fmt.Errorf("Error: Postal Code FSA: %s is invalid, should be 3 characters.", postalCodeFSA)
	}

	// Configurable patterns for FSA.
	FSAAllowedFirstChar := "[ABCEGHJKLMNPRSTVXY]"
	FSAAllowedSecondChar := "[0-9]"
	FSAAllowedThirdChar := "[A-Z]"

	if err := ValidateFSAConfigString(FSAAllowedFirstChar); err != nil {
		return false, "", err
	}
	if err := ValidateFSAConfigString(FSAAllowedSecondChar); err != nil {
		return false, "", err
	}
	if err := ValidateFSAConfigString(FSAAllowedThirdChar); err != nil {
		return false, "", err
	}

	postalCodeFSARegex := regexp.MustCompile(fmt.Sprintf(`^%s%s%s$`, FSAAllowedFirstChar, FSAAllowedSecondChar, FSAAllowedThirdChar))
	result := postalCodeFSARegex.MatchString(strings.ToUpper(postalCodeFSA))

	return result, postalCodeFSARegex.String(), nil
}
