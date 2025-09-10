package utils

import (
	"regexp"
	"strings"
)

// ConvertMarkdownToTelegram converts Markdown formatting to Telegram's format
func ConvertMarkdownToTelegram(text string) string {
	// Convert headers first
	text = convertHeaders(text)
	
	// Convert code blocks (before other formatting)
	text = convertCodeBlocks(text)
	
	// Convert lists (before italic to avoid conflicts)
	text = convertLists(text)
	
	// Convert bold text
	text = convertBold(text)
	
	// Convert italic text
	text = convertItalic(text)
	
	// Convert links (basic support)
	text = convertLinks(text)
	
	// Clean up extra whitespace
	text = strings.TrimSpace(text)
	
	return text
}

// convertHeaders converts Markdown headers to Telegram format
func convertHeaders(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// H1: # Header -> **Header**
		if strings.HasPrefix(trimmed, "# ") {
			header := strings.TrimPrefix(trimmed, "# ")
			result = append(result, "**"+header+"**")
			continue
		}
		
		// H2: ## Header -> **Header**
		if strings.HasPrefix(trimmed, "## ") {
			header := strings.TrimPrefix(trimmed, "## ")
			result = append(result, "**"+header+"**")
			continue
		}
		
		// H3: ### Header -> **Header**
		if strings.HasPrefix(trimmed, "### ") {
			header := strings.TrimPrefix(trimmed, "### ")
			result = append(result, "**"+header+"**")
			continue
		}
		
		// H4-H6: #### Header -> **Header**
		if strings.HasPrefix(trimmed, "#### ") || strings.HasPrefix(trimmed, "##### ") || strings.HasPrefix(trimmed, "###### ") {
			header := strings.TrimPrefix(trimmed, "#")
			header = strings.TrimPrefix(header, "#")
			header = strings.TrimPrefix(header, "#")
			header = strings.TrimPrefix(header, "#")
			header = strings.TrimPrefix(header, "#")
			header = strings.TrimSpace(header)
			result = append(result, "**"+header+"**")
			continue
		}
		
		result = append(result, line)
	}
	
	return strings.Join(result, "\n")
}

// convertBold converts Markdown bold to Telegram format
func convertBold(text string) string {
	// Convert __text__ to **text**
	text = strings.ReplaceAll(text, "__", "**")
	return text
}

// convertItalic converts Markdown italic to Telegram format
func convertItalic(text string) string {
	// Convert *text* to _text_ (Telegram format)
	// Simple approach: find single asterisks that are not part of double asterisks
	// We'll use a simple string replacement approach
	
	// First, protect existing **bold** by replacing with a temporary marker
	text = strings.ReplaceAll(text, "**", "___BOLD_MARKER___")
	
	// Now convert single asterisks to underscores
	text = strings.ReplaceAll(text, "*", "_")
	
	// Restore bold markers
	text = strings.ReplaceAll(text, "___BOLD_MARKER___", "**")
	
	return text
}

// convertCodeBlocks converts Markdown code blocks to Telegram format
func convertCodeBlocks(text string) string {
	// Convert ```language\ncode\n``` to ```\ncode\n```
	// Use a more flexible regex that handles multiline content
	text = regexp.MustCompile("```\\w*\\n([\\s\\S]*?)\\n```").ReplaceAllString(text, "```\n$1\n```")
	return text
}

// convertInlineCode converts Markdown inline code to Telegram format
func convertInlineCode(text string) string {
	// Convert `code` to `code` (already Telegram format)
	return text
}

// convertLists converts Markdown lists to Telegram format
func convertLists(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Convert unordered lists: - item -> • item
		if strings.HasPrefix(trimmed, "- ") {
			item := strings.TrimPrefix(trimmed, "- ")
			result = append(result, "• "+item)
			continue
		}
		
		// Convert unordered lists: * item -> • item
		if strings.HasPrefix(trimmed, "* ") {
			item := strings.TrimPrefix(trimmed, "* ")
			result = append(result, "• "+item)
			continue
		}
		
		// Convert ordered lists: 1. item -> 1. item (keep as is)
		if regexp.MustCompile(`^\d+\.\s`).MatchString(trimmed) {
			result = append(result, line)
			continue
		}
		
		result = append(result, line)
	}
	
	return strings.Join(result, "\n")
}

// convertLinks converts Markdown links to Telegram format
func convertLinks(text string) string {
	// Convert [text](url) to [text](url) (already Telegram format)
	return text
}

// EscapeTelegramMarkdown escapes special characters for Telegram
func EscapeTelegramMarkdown(text string) string {
	// Escape special characters that have meaning in Telegram Markdown
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "*", "\\*")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	text = strings.ReplaceAll(text, "~", "\\~")
	text = strings.ReplaceAll(text, "`", "\\`")
	text = strings.ReplaceAll(text, ">", "\\>")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, "+", "\\+")
	text = strings.ReplaceAll(text, "-", "\\-")
	text = strings.ReplaceAll(text, "=", "\\=")
	text = strings.ReplaceAll(text, "|", "\\|")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, ".", "\\.")
	text = strings.ReplaceAll(text, "!", "\\!")
	
	return text
}