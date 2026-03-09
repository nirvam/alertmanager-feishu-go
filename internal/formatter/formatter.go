package formatter

import (
	"fmt"
	"strings"

	"github.com/nirvam/alertmanager-feishu-go/internal/feishu"
	"github.com/nirvam/alertmanager-feishu-go/pkg/alertmanager"
)

// FormatText formats Alertmanager payload into a plain text Feishu message.
func FormatText(payload alertmanager.Payload) feishu.Message {
	alertname := payload.CommonLabels["alertname"]
	if alertname == "" {
		alertname = "Alert"
	}

	lines := []string{
		fmt.Sprintf("Status: %s", strings.ToUpper(payload.Status)),
		fmt.Sprintf("Alert: %s", alertname),
	}

	for _, alert := range payload.Alerts {
		summary := alert.Annotations["summary"]
		if summary == "" {
			summary = "No summary"
		}
		description := alert.Annotations["description"]
		if description == "" {
			description = "No description"
		}

		lines = append(lines, "\n--- Alert ---")
		lines = append(lines, fmt.Sprintf("Summary: %s", summary))
		lines = append(lines, fmt.Sprintf("Description: %s", description))
		lines = append(lines, fmt.Sprintf("Starts At: %s", alert.StartsAt.Format("2006-01-02 15:04:05")))
	}

	return feishu.Message{
		MsgType: "text",
		Content: map[string]interface{}{
			"text": strings.Join(lines, "\n"),
		},
	}
}

// FormatCard formats Alertmanager payload into a Feishu interactive card (Schema 2.0).
func FormatCard(payload alertmanager.Payload) feishu.Message {
	alertname := payload.CommonLabels["alertname"]
	if alertname == "" {
		alertname = "Alert"
	}

	headerTemplate := "green"
	if payload.Status == "firing" {
		headerTemplate = "red"
	}

	elements := []interface{}{}

	for _, alert := range payload.Alerts {
		summary := alert.Annotations["summary"]
		if summary == "" {
			summary = "No summary"
		}
		description := alert.Annotations["description"]
		if description == "" {
			description = "No description"
		}

		severity := alert.Labels["severity"]
		if severity == "" {
			severity = "unknown"
		}

		severityColor := "grey"
		switch strings.ToLower(severity) {
		case "critical":
			severityColor = "red"
		case "warning":
			severityColor = "orange"
		}

		content := fmt.Sprintf("**Summary:** %s\n**Description:** %s\n**Severity:** <font color='%s'>%s</font>\n**Starts At:** %s",
			summary, description, severityColor, severity, alert.StartsAt.Format("2006-01-02 15:04:05"))

		elements = append(elements, feishu.MarkdownElement{
			Tag:     "markdown",
			Content: content,
		})

		if alert.GeneratorURL != "" {
			elements = append(elements, feishu.ButtonElement{
				Tag: "button",
				Text: map[string]interface{}{
					"tag":     "plain_text",
					"content": "View in Prometheus",
				},
				Type: "default",
				Behaviors: []map[string]interface{}{
					{
						"type":        "open_url",
						"default_url": alert.GeneratorURL,
					},
				},
			})
		}

		elements = append(elements, feishu.HrElement{
			Tag: "hr",
		})
	}

	// Remove last hr
	if len(elements) > 0 {
		if hr, ok := elements[len(elements)-1].(feishu.HrElement); ok && hr.Tag == "hr" {
			elements = elements[:len(elements)-1]
		}
	}

	title := fmt.Sprintf("[%s] %s", strings.ToUpper(payload.Status), alertname)

	card := &feishu.Card{
		Schema: "2.0",
		Config: &feishu.Config{
			UpdateMulti: true,
			Summary: map[string]string{
				"content": title,
			},
		},
		Header: &feishu.Header{
			Title: map[string]interface{}{
				"tag":     "plain_text",
				"content": title,
			},
			Template: headerTemplate,
		},
		Body: &feishu.Body{
			Direction: "vertical",
			Elements:  elements,
		},
	}

	return feishu.Message{
		MsgType: "interactive",
		Card:    card,
	}
}
