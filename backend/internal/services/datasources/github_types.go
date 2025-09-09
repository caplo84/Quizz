package datasources

import "time"

type GitHubConfig struct {
	Token      string
	Owner      string
	Repository string // Changed from Repo
	BaseURL    string
}

type GitHubFile struct { // Fixed typo
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"` // Changed from Sha
	Size        int    `json:"size"`
	URL         string `json:"url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content,omitempty"` // Fixed spacing
	Encoding    string `json:"encoding,omitempty"`
}

type GitHubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type ParsedQuiz struct {
	Title       string
	Description string
	Category    string
	Questions   []ParsedQuestion
	Source      string
	ExternalID  string
	LastCommit  string
}

type ParsedQuestion struct {
	Text    string
	Choices []string
	Correct int
	Type    string
}
