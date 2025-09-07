package datasources

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGitHubClient_ValidateAuth(t *testing.T) {
    // Mock tests
    config := GitHubConfig{
        Token:      "test-token",
        Owner:      "test-owner",
        Repository: "test-repo",
    }
    
    client := NewGitHubClient(config)
    
    // Test with mock server
    // Implementation needed
    assert.NotNil(t, client)
}

func TestGitHubClient_ParseMarkdown(t *testing.T) {
    client := NewGitHubClient(GitHubConfig{})
    
    markdownContent := `
# JavaScript Quiz

#### Q1. What is JavaScript?
- [ ] A markup language
- [x] A programming language
- [ ] A database
`
    
    file := GitHubFile{
        Name:    "test.md",
        Content: markdownContent,
        SHA:     "test-sha",
    }
    
    quiz, err := client.parseMarkdownQuiz(file)
    require.NoError(t, err)
    assert.Equal(t, "JavaScript Quiz", quiz.Title)
    assert.Len(t, quiz.Questions, 1)
}