package validator

import "regexp"

// EmailRegex represents an email regular expression.
var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// PasswordRegex represents a password regular expression.
var PasswordRegex = regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})")

// Validator encapsulates a set of validation rules and methods.
type Validator struct {
	Errors map[string]string
}

// New returns a configured instance of *Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// AddError adds new errors to the errors map.
func (v *Validator) AddError(errKey, errMsg string) {
	// Check for duplicate keys to avoid over writing the previous value.
	if _, exists := v.Errors[errKey]; !exists {
		v.Errors[errKey] = errMsg
	}
}

// IsValid checks if a particular value is valid according to the applied rules.
func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

// Check performs a boolean check on the value and appends the errMsg to the errors map.
func (v *Validator) Check(condition bool, errKey, errMsg string) {
	if !condition {
		v.AddError(errKey, errMsg)
	}
}

// In returns true if a specific value is in a list of strings.
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx regexp.Regexp) bool {
	return rx.MatchString(value)
}

// IsUnique returns true if all string values in a slice are unique.
func IsUnique(list ...string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range list {
		uniqueValues[value] = true
	}

	return len(list) == len(uniqueValues)
}
