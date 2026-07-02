package service

import (
	"regexp"
	"strings"

	"github.com/stackriv/dev-tools/internal/pkg"
)

func MarkdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var out strings.Builder
	inCodeBlock := false
	inList := false
	inOrderedList := false

	for _, line := range lines {
		// Code blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				out.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				lang := strings.TrimPrefix(line, "```")
				if inList {
					out.WriteString("</ul>\n")
					inList = false
				}
				if inOrderedList {
					out.WriteString("</ol>\n")
					inOrderedList = false
				}
				if lang != "" {
					out.WriteString(`<pre><code class="language-` + pkg.EscapeHTML(lang) + `">`)
				} else {
					out.WriteString("<pre><code>")
				}
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			out.WriteString(pkg.EscapeHTML(line) + "\n")
			continue
		}

		// Headings
		if strings.HasPrefix(line, "######") {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h6>" + inlineMarkdown(strings.TrimPrefix(line, "###### ")) + "</h6>\n")
			continue
		}
		if strings.HasPrefix(line, "#####") {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h5>" + inlineMarkdown(strings.TrimPrefix(line, "##### ")) + "</h5>\n")
			continue
		}
		if strings.HasPrefix(line, "####") {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h4>" + inlineMarkdown(strings.TrimPrefix(line, "#### ")) + "</h4>\n")
			continue
		}
		if strings.HasPrefix(line, "###") {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h3>" + inlineMarkdown(strings.TrimPrefix(line, "### ")) + "</h3>\n")
			continue
		}
		if strings.HasPrefix(line, "##") {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h2>" + inlineMarkdown(strings.TrimPrefix(line, "## ")) + "</h2>\n")
			continue
		}
		if strings.HasPrefix(line, "#") {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h1>" + inlineMarkdown(strings.TrimPrefix(line, "# ")) + "</h1>\n")
			continue
		}

		// Horizontal rule
		if line == "---" || line == "***" || line == "___" {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<hr>\n")
			continue
		}

		// Blockquote
		if strings.HasPrefix(line, "> ") {
			out.WriteString("<blockquote>" + inlineMarkdown(strings.TrimPrefix(line, "> ")) + "</blockquote>\n")
			continue
		}

		// Unordered list
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			if inOrderedList {
				out.WriteString("</ol>\n")
				inOrderedList = false
			}
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			content := line[2:]
			out.WriteString("<li>" + inlineMarkdown(content) + "</li>\n")
			continue
		}

		// Ordered list
		olRe := regexp.MustCompile(`^\d+\.\s(.+)`)
		if m := olRe.FindStringSubmatch(line); m != nil {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			if !inOrderedList {
				out.WriteString("<ol>\n")
				inOrderedList = true
			}
			out.WriteString("<li>" + inlineMarkdown(m[1]) + "</li>\n")
			continue
		}

		// Close lists on empty line
		if strings.TrimSpace(line) == "" {
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			if inOrderedList {
				out.WriteString("</ol>\n")
				inOrderedList = false
			}
			out.WriteString("<br>\n")
			continue
		}

		// Table
		if strings.Contains(line, "|") {
			out.WriteString(parseTableLine(line))
			continue
		}

		// Paragraph
		if inList {
			out.WriteString("</ul>\n")
			inList = false
		}
		if inOrderedList {
			out.WriteString("</ol>\n")
			inOrderedList = false
		}
		out.WriteString("<p>" + inlineMarkdown(line) + "</p>\n")
	}

	if inList {
		out.WriteString("</ul>\n")
	}
	if inOrderedList {
		out.WriteString("</ol>\n")
	}
	if inCodeBlock {
		out.WriteString("</code></pre>\n")
	}

	return out.String()
}

func inlineMarkdown(s string) string {
	// Bold italic
	s = regexp.MustCompile(`\*\*\*(.+?)\*\*\*`).ReplaceAllString(s, "<strong><em>$1</em></strong>")
	// Bold
	s = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(s, "<strong>$1</strong>")
	s = regexp.MustCompile(`__(.+?)__`).ReplaceAllString(s, "<strong>$1</strong>")
	// Italic
	s = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(s, "<em>$1</em>")
	s = regexp.MustCompile(`_(.+?)_`).ReplaceAllString(s, "<em>$1</em>")
	// Strikethrough
	s = regexp.MustCompile(`~~(.+?)~~`).ReplaceAllString(s, "<del>$1</del>")
	// Inline code
	s = regexp.MustCompile("`(.+?)`").ReplaceAllString(s, "<code>$1</code>")
	// Links
	s = regexp.MustCompile(`\[(.+?)\]\((.+?)\)`).ReplaceAllString(s, `<a href="$2" target="_blank">$1</a>`)
	// Images
	s = regexp.MustCompile(`!\[(.+?)\]\((.+?)\)`).ReplaceAllString(s, `<img src="$2" alt="$1">`)
	return s
}

func parseTableLine(line string) string {
	if regexp.MustCompile(`^\|?[\s\-:|]+\|`).MatchString(line) {
		return ""
	}
	cells := strings.Split(strings.Trim(line, "|"), "|")
	var sb strings.Builder
	sb.WriteString("<tr>")
	for _, cell := range cells {
		sb.WriteString("<td>" + inlineMarkdown(strings.TrimSpace(cell)) + "</td>")
	}
	sb.WriteString("</tr>\n")
	return sb.String()
}
