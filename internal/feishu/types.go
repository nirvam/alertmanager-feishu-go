package feishu

// Message is the top-level structure for a Feishu webhook message.
type Message struct {
	MsgType   string                 `json:"msg_type"`
	Content   map[string]interface{} `json:"content,omitempty"`
	Card      *Card                  `json:"card,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
	Sign      string                 `json:"sign,omitempty"`
}

// Card represents a Feishu message card (Schema 2.0).
type Card struct {
	Schema string `json:"schema"`
	Config *Config `json:"config,omitempty"`
	Header *Header `json:"header,omitempty"`
	Body   *Body   `json:"body,omitempty"`
}

// Config contains card configuration.
type Config struct {
	UpdateMulti bool             `json:"update_multi,omitempty"`
	Summary     map[string]string `json:"summary,omitempty"`
}

// Header represents the card header.
type Header struct {
	Title    map[string]interface{} `json:"title,omitempty"`
	Template string                 `json:"template,omitempty"`
}

// Body represents the card body.
type Body struct {
	Direction string        `json:"direction,omitempty"`
	Elements  []interface{} `json:"elements,omitempty"`
}

// MarkdownElement represents a markdown content element.
type MarkdownElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// ButtonElement represents a button element.
type ButtonElement struct {
	Tag       string                 `json:"tag"`
	Text      map[string]interface{} `json:"text"`
	Type      string                 `json:"type"`
	Behaviors []map[string]interface{} `json:"behaviors,omitempty"`
}

// HrElement represents a horizontal rule element.
type HrElement struct {
	Tag string `json:"tag"`
}
