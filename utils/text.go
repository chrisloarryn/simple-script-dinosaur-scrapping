package utils

import "strings"

// CleanText removes newlines and tabs and trims spaces. If cleaned is empty, returns original input.
func CleanText(input string) string {
	cleaned := strings.TrimSpace(strings.ReplaceAll(input, "\n", ""))
	cleaned = strings.TrimSpace(strings.ReplaceAll(cleaned, "\t", ""))
	if len(cleaned) > 0 {
		return cleaned
	}
	return input
}

// CamelCase converts a phrase into lowerCamelCase.
func CamelCase(input string) string {
	parts := strings.Fields(input)
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.Title(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}
