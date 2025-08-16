package utils

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "time"
)

// GenerateID generates a random hex ID
func GenerateID(length int) string {
    bytes := make([]byte, length)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// FormatDuration formats duration in a human-readable way
func FormatDuration(d time.Duration) string {
    if d < time.Minute {
        return fmt.Sprintf("%.0fs", d.Seconds())
    }
    if d < time.Hour {
        return fmt.Sprintf("%.0fm", d.Minutes())
    }
    return fmt.Sprintf("%.1fh", d.Hours())
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// RemoveDuplicates removes duplicate strings from slice
func RemoveDuplicates(slice []string) []string {
    keys := make(map[string]bool)
    result := []string{}
    
    for _, item := range slice {
        if !keys[item] {
            keys[item] = true
            result = append(result, item)
        }
    }
    
    return result
}