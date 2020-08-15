package webhooks

type EmbedFooter struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type EmbedAuthor struct {
	URL     string `json:"url,omitempty"`
	Name    string `json:"name,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	URL         string        `json:"url,omitempty"`
	Type        string        `json:"type,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Timestamp   string        `json:"timestamp,omitempty"`
	Color       int           `json:"color,omitempty"`
	Footer      *EmbedFooter  `json:"footer,omitempty"`
	Author      *EmbedAuthor  `json:"author,omitempty"`
	Fields      []*EmbedField `json:"fields,omitempty"`
}

type Message struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}
