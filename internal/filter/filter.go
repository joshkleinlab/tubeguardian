package filter

import (
	"bufio"
	"os"
	"strings"

	"github.com/joshkleinlab/tubeguardian/internal/ahocorasick"
)

// Matcher holds the Aho-Corasick automaton for banned words
type Matcher struct {
	ac    *ahocorasick.Matcher
	words []string
}

// LoadKeywords loads keywords from a file into a Matcher
func LoadKeywords(path string) (*Matcher, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if w != "" {
			words = append(words, w)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	m := ahocorasick.NewStringMatcher(words)
	return &Matcher{
		ac:    m,
		words: words,
	}, nil
}

// Match finds all banned keywords inside the given text
func (m *Matcher) Match(text string) []string {
	matches := m.ac.Match([]byte(strings.ToLower(text)))

	var results []string
	for _, idx := range matches {
		results = append(results, m.words[idx])
	}
	return results
}
