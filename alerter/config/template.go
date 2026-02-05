package config

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
)

//go:embed templates
var embeddedTemplates embed.FS

// MailMatter contient les métadonnées du mail (frontmatter)
type MailMatter struct {
	Subject string `yaml:"subject"`
}

// GetStringFromEmbeddedTemplate parse un template avec les données fournies
// Retourne le contenu du mail et les métadonnées (subject)
func GetStringFromEmbeddedTemplate(templatePath string, data interface{}) (content string, matter MailMatter, err error) {
	var temp *template.Template

	// Récupérer le fichier template depuis embeddedTemplates
	temp, err = template.ParseFS(embeddedTemplates, templatePath)
	if err != nil {
		return
	}

	// Exécuter le template avec les données
	var tpl bytes.Buffer
	if err = temp.Execute(&tpl, data); err != nil {
		return
	}

	// Parser le frontmatter pour extraire le subject
	var mailContent []byte
	mailContent, err = frontmatter.Parse(strings.NewReader(tpl.String()), &matter)
	if err == nil {
		content = string(mailContent)
	}

	return
}

// GetHTMLTemplate retourne le contenu HTML formaté
func GetHTMLTemplate(data interface{}) (content string, subject string, err error) {
	content, matter, err := GetStringFromEmbeddedTemplate("templates/event_modified.html", data)
	if err != nil {
		return "", "", err
	}
	return content, matter.Subject, nil
}

// GetTextTemplate retourne le contenu texte formaté
func GetTextTemplate(data interface{}) (content string, subject string, err error) {
	content, matter, err := GetStringFromEmbeddedTemplate("templates/event_modified.txt", data)
	if err != nil {
		return "", "", err
	}
	return content, matter.Subject, nil
}
