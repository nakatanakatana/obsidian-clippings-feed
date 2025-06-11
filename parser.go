package clippingsfeed

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

type Metadata struct {
	Title       string    `json:"title"`
	Site        string    `json:"site"`
	Source      string    `json:"source"`
	Author      []string  `json:"author"`
	Published   string    `json:"published"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
}

func CreateParser() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)
}

func ParseMeta(md goldmark.Markdown, source string) (*Metadata, error) {
	document := md.Parser().Parse(text.NewReader([]byte(source)))

	jsonString, err := json.Marshal(document.OwnerDocument().Meta())
	if err != nil {
		return nil, fmt.Errorf("marshal Error: %w", err)
	}

	metadata := Metadata{}
	err = json.Unmarshal(jsonString, &metadata)
	if err != nil {
		return nil, fmt.Errorf("unmarshal Error: %w", err)
	}

	return &metadata, nil
}
