package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type Validator struct {
	Validate *validator.Validate
}

type ErrorValidation struct {
	Field string      `json:"field"`
	Tag   string      `json:"tag"`
	Param string      `json:"param"`
	Value interface{} `json:"value"`
}

// NewValidator initializes and returns a CustomValidator with production-ready validations
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validations
	err := v.RegisterValidation("is-valid-email", validateEmail)
	if err != nil {
		panic(fmt.Errorf("failed register custom validation: %v", err))
	}
	err = v.RegisterValidation("is-strong-password", validatePassword)
	if err != nil {
		panic(fmt.Errorf("failed register custom validation: %v", err))
	}
	err = v.RegisterValidation("no-special-chars", validateNoSpecialChars)
	if err != nil {
		panic(fmt.Errorf("failed register custom validation: %v", err))
	}
	err = v.RegisterValidation("required-strict", validateStrictRequired)
	if err != nil {
		panic(fmt.Errorf("failed register custom validation: %v", err))
	}
	err = v.RegisterValidation("is-valid-phone-number", validateE164PhoneNumber)
	if err != nil {
		panic(fmt.Errorf("failed register custom validation: %v", err))
	}
	err = v.RegisterValidation("is-valid-date", validateDate)
	if err != nil {
		panic(fmt.Errorf("failed register date custom validation: %v", err))
	}
	err = v.RegisterValidation("is-valid-datetime", validateDateTime)
	if err != nil {
		panic(fmt.Errorf("failed register datetime custom validation: %v", err))
	}

	return &Validator{Validate: v}
}

// validateEmail is a custom validation function for validating email format and domain restriction
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Complex regex for email validation
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if !re.MatchString(email) {
		return false
	}

	return true
}

// validatePassword checks if the password meets complexity requirements dynamically from config
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 25 {
		return false
	}

	var hasUpper = strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var hasLower = strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	var hasNumber = strings.ContainsAny(password, "0123456789")
	var hasSpecial = strings.ContainsAny(password, "!@#$%^&*()-_+={}[]|:;<>,.?/~`")

	// Example rule: must contain upper, lower, number, and special character
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateNoSpecialChars ensures the field does not contain special characters
func validateNoSpecialChars(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Check if there are any special characters using a regular expression
	re := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	return re.MatchString(value)
}

// validateStrictRequired forces certain fields to be non-empty in strict mode
func validateStrictRequired(fl validator.FieldLevel) bool {
	return fl.Field().String() != ""
}

func validateE164PhoneNumber(fl validator.FieldLevel) bool {
	phone := strings.ReplaceAll(fl.Field().String(), " ", "")

	// Define the regex for E.164 phone numbers
	e164Regex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	// Check if phone matches the E.164 format
	return e164Regex.MatchString(phone)
}

func validateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// validateDateTime validates if the field is a valid RFC3339 datetime
func validateDateTime(fl validator.FieldLevel) bool {
	datetime := fl.Field().String()
	if datetime == "" {
		return true
	}

	// Format: 2006-01-02T15:04:05Z07:00
	_, err := time.Parse(time.RFC3339, datetime)
	if err == nil {
		return true
	}

	// Format: 2006-01-02T15:04:05
	_, err = time.Parse("2006-01-02T15:04:05", datetime)
	if err == nil {
		return true
	}

	return false
}

func (v *Validator) ValidateStruct(obj interface{}) []*ErrorValidation {
	err := v.Validate.Struct(obj)
	if err != nil {
		var errors []*ErrorValidation
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorValidation
			element.Field = cutAfterFirstDot(err.Namespace())
			element.Tag = err.Tag()
			element.Value = err.Value()
			element.Param = err.Param()
			errors = append(errors, &element)
		}

		return errors
	}

	return nil
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func cutAfterFirstDot(namespace string) string {
	index := strings.Index(namespace, ".")
	if index == -1 {
		return namespace
	}

	errField := strings.Split(namespace[index+1:], ".")
	for i := range errField {
		errField[i] = toSnakeCase(errField[i])
	}

	return strings.Join(errField, ".")
}
