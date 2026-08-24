package service

import (
	"strings"
	"testing"
)

func TestNormalizeSlugSupportsChineseAndASCII(t *testing.T) {
	tests := map[string]string{
		"hahah":            "hahah",
		"美线新闻":             "美线新闻",
		"  美线 News 2026! ": "美线-news-2026",
	}

	for input, want := range tests {
		if got := normalizeSlug(input); got != want {
			t.Errorf("normalizeSlug(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGeneratedArticleSlugIsUnique(t *testing.T) {
	first := generatedArticleSlug("美线新闻", 1001)
	second := generatedArticleSlug("美线新闻", 1002)

	if first == second {
		t.Fatalf("generated slugs must be unique: %q", first)
	}
	if first != "美线新闻-1001" {
		t.Fatalf("unexpected generated slug: %q", first)
	}
}

func TestGeneratedArticleSlugFitsDatabaseColumn(t *testing.T) {
	slug := generatedArticleSlug(strings.Repeat("新", 220), 1234567890123456789)
	if length := len([]rune(slug)); length > 220 {
		t.Fatalf("slug length is %d, want at most 220", length)
	}
}
