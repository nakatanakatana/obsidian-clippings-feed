package clippingsfeed_test

import (
	"encoding/json"
	"testing"

	"github.com/nakatanakatana/obsidian-clippings-feed"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestParseMeta(t *testing.T) {
	md := clippingsfeed.CreateParser()

	for name, tt := range map[string]struct {
		goldenFilename string
		source         string
	}{
		"parse success": {
			goldenFilename: "success",
			source: `---
title: "item title"
site: "dummy"
source: "https://github.com/nakatanakatana/obsidian-feed"
author:
published: 2025-06-01
created: 2025-06-03T12:54:50+09:00
description:
tags:
  - "go"
  - "obsidian"
---
document`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			metadata, err := clippingsfeed.ParseMeta(md, tt.source)
			assert.NilError(t, err)

			result, _ := json.Marshal(metadata)
			golden.AssertBytes(t, result, tt.goldenFilename)
		})
	}
}
