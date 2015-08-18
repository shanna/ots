package ots

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

// Defaults from libots.
var libotsLanguages = []string{"bg", "ca", "cs", "cy", "da", "de", "el", "en", "eo", "es", "et", "eu", "fi", "fr", "ga", "gl", "he", "hu", "ia", "id", "is", "it", "lv", "mi", "ms", "mt", "nl", "nn", "pl", "pt", "ro", "ru", "sv", "tl", "tr", "uk", "yi"}

var sampleText = `The hawksbill turtle is a critically endangered sea turtle belonging to the family Cheloniidae.
It is the only species in its genus. The species has a worldwide distribution, with Atlantic and
Pacific subspecies.`

func TestLanguages(t *testing.T) {
	languages, err := Languages()
	assert.Equal(t, nil, err)
	assert.Equal(t, libotsLanguages, languages)
}

func TestParse(t *testing.T) {
	article, err := Parse(sampleText, "en")
	assert.Equal(t, nil, err)
	assert.Equal(t, "en", article.Language)
}

func TestKeywords(t *testing.T) {
	article, _ := Parse(sampleText, "en")
	keywords := article.Keywords()
	assert.Equal(t, []string{"species", "turtle", "subspecies", "pacific", "atlantic"}, keywords)
}

func TestArticleSentences(t *testing.T) {
	article, _ := Parse(sampleText, "en")
	summary := article.Sentences(1)
	fmt.Printf("%+v\n", summary)
	// assert.Equal(t, []Summary{Summary{"test", 1.0}}, summary)
}
