package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidateSlug checks if a string is a valid slug format
func ValidateSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > 100 {
		return false
	}

	// Slug should contain only lowercase letters, numbers, and hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	return matched
}

// ValidateEmail checks if a string is a valid email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	clean := htmlRegex.ReplaceAllString(input, "")

	// Trim whitespace
	clean = strings.TrimSpace(clean)

	return clean
}

// ValidatePassword checks password strength
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// NormalizeString normalizes strings for comparison
func NormalizeString(input string) string {
	// Convert to lowercase and trim
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized = spaceRegex.ReplaceAllString(normalized, " ")

	return normalized
}

// ValidateQuizTitle checks if quiz title is valid
func ValidateQuizTitle(title string) bool {
	title = strings.TrimSpace(title)
	return len(title) >= 3 && len(title) <= 200
}

// ValidateQuestionText checks if question text is valid
func ValidateQuestionText(question string) bool {
	question = strings.TrimSpace(question)
	return len(question) >= 10 && len(question) <= 1000
}

// ValidateChoiceText checks if choice text is valid
func ValidateChoiceText(choice string) bool {
	choice = strings.TrimSpace(choice)
	return len(choice) >= 1 && len(choice) <= 500
}
