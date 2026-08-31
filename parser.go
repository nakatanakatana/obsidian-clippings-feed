package clippingsfeed

import (
	"encoding/json"
	"fmt"
	"time"

	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
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

func CreateParser() parser.Parser {
	return parser.New(
		parser.WithExtensions(
			meta.Parser,
		),
	)
}

func ParseMeta(p parser.Parser, source string) (*Metadata, error) {
	document := p.ParseStringSource(source)

	doc, ok := document.(*ast.Document)
	if !ok {
		return nil, fmt.Errorf("failed to cast to ast.Document")
	}

	jsonString, err := json.Marshal(doc.Metadata())
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
