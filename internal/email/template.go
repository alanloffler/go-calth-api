package email

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

func renderTemplate(name string, data map[string]string) (string, error) {
	content, err := templateFS.ReadFile(fmt.Sprintf("templates/%s.html", name))
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", name, err)
	}

	html := string(content)
	for key, value := range data {
		html = strings.ReplaceAll(html, "{{"+key+"}}", value)
	}

	return html, nil
}
