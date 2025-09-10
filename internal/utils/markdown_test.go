package utils

import (
	"testing"
)

func TestConvertMarkdownToTelegram(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Headers",
			input:    "# Main Title\n## Subtitle\n### Section",
			expected: "**Main Title**\n**Subtitle**\n**Section**",
		},
		{
			name:     "Bold text with underscores",
			input:    "This is __bold__ text",
			expected: "This is **bold** text",
		},
		{
			name:     "Italic text",
			input:    "This is *italic* text",
			expected: "This is _italic_ text",
		},
		{
			name:     "Code blocks",
			input:    "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			expected: "```\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
		},
		{
			name:     "Inline code",
			input:    "Use `fmt.Println()` function",
			expected: "Use `fmt.Println()` function",
		},
		{
			name:     "Unordered lists",
			input:    "- First item\n- Second item\n* Third item",
			expected: "• First item\n• Second item\n• Third item",
		},
		{
			name:     "Ordered lists",
			input:    "1. First item\n2. Second item\n3. Third item",
			expected: "1. First item\n2. Second item\n3. Third item",
		},
		{
			name:     "Simple mixed formatting",
			input:    "# Title\n\nThis is __bold__ and *italic* text.\n\n- Item 1\n- Item 2",
			expected: "**Title**\n\nThis is **bold** and _italic_ text.\n\n• Item 1\n• Item 2",
		},
		{
			name:     "Links",
			input:    "Visit [Google](https://google.com) for search",
			expected: "Visit [Google](https://google.com) for search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertMarkdownToTelegram(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertMarkdownToTelegram() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEscapeTelegramMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic escaping",
			input:    "This has _underscores_ and *asterisks*",
			expected: "This has \\_underscores\\_ and \\*asterisks\\*",
		},
		{
			name:     "Special characters",
			input:    "Text with [brackets] and (parentheses)",
			expected: "Text with \\[brackets\\] and \\(parentheses\\)",
		},
		{
			name:     "Code-like text",
			input:    "Use `backticks` and ```code blocks```",
			expected: "Use \\`backticks\\` and \\`\\`\\`code blocks\\`\\`\\`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeTelegramMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeTelegramMarkdown() = %v, want %v", result, tt.expected)
			}
		})
	}
}
