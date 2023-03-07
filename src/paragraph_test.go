package src

import (
	"testing"
)

func TestSplitParagraphs(t *testing.T) {
	text := "This is a test sentence. This is another test sentence."
	maxChars := 50
	expected := []string{"This is a test sentence.", "This is another test sentence."}
	result := SplitParagraphs(text, maxChars)

	if len(expected) != len(result) {
		t.Errorf("Expected %d paragraphs, but got %d", len(expected), len(result))
	}

	for i, paragraph := range result {
		if paragraph != expected[i] {
			t.Errorf("Paragraph %d does not match. Expected '%s', but got '%s'", i, expected[i], paragraph)
		}
	}
}

func TestSplitParagraphsShortInput(t *testing.T) {
	text := "This is a short test."
	maxChars := 100
	expected := []string{"This is a short test."}
	result := SplitParagraphs(text, maxChars)

	if len(expected) != len(result) {
		t.Errorf("Expected %d paragraphs, but got %d", len(expected), len(result))
	}

	if result[0] != expected[0] {
		t.Errorf("Paragraph does not match. Expected '%s', but got '%s'", expected[0], result[0])
	}
}
