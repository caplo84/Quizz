package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug_Basic(t *testing.T) {
	assert.Equal(t, "hello-world", GenerateSlug("Hello World"))
}

func TestGenerateSlug_Underscores(t *testing.T) {
	assert.Equal(t, "hello-world", GenerateSlug("hello_world"))
}

func TestGenerateSlug_SpecialChars(t *testing.T) {
	assert.Equal(t, "go-programming", GenerateSlug("Go Programming!"))
}

func TestGenerateSlug_MultipleSpaces(t *testing.T) {
	assert.Equal(t, "a-b", GenerateSlug("a   b"))
}

func TestGenerateSlug_LeadingTrailingHyphens(t *testing.T) {
	result := GenerateSlug("  hello  ")
	assert.Equal(t, "hello", result)
}

func TestGenerateSlug_AlreadyLowercase(t *testing.T) {
	assert.Equal(t, "javascript", GenerateSlug("javascript"))
}

func TestGenerateSlug_Numbers(t *testing.T) {
	assert.Equal(t, "quiz-2024", GenerateSlug("Quiz 2024"))
}

func TestGenerateSlug_EmptyString(t *testing.T) {
	assert.Equal(t, "", GenerateSlug(""))
}

func TestGenerateSlug_OnlySpecialChars(t *testing.T) {
	assert.Equal(t, "", GenerateSlug("!@#$%"))
}

// TestGenerateSlug_MixedCaseWithHyphens verifies that dots are removed (not converted to hyphens)
func TestGenerateSlug_MixedCaseWithHyphens(t *testing.T) {
	assert.Equal(t, "reactjs-basics", GenerateSlug("React.js Basics"))
}

func TestGenerateSlug_ConsecutiveHyphens(t *testing.T) {
	result := GenerateSlug("a--b")
	assert.Equal(t, "a-b", result)
}
