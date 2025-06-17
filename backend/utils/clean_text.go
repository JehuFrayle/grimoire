package utils

import "strings"

// CleanText sanitizes input for safe use in SQL string literals.
// Note: Always use parameterized queries to prevent SQL injection.
func CleanText(text string) string {
	// Remove leading and trailing whitespace
	text = strings.TrimSpace(text)

	// Replace multiple spaces with a single space
	text = strings.Join(strings.Fields(text), " ")

	// Remove any non-printable characters
	var builder strings.Builder
	for _, r := range text {
		if r >= 32 && r <= 126 { // ASCII printable characters
			builder.WriteRune(r)
		}
	}
	cleanedText := builder.String()

	// Escape single quotes for SQL string literals
	cleanedText = strings.ReplaceAll(cleanedText, "'", "''")

	return cleanedText
}
