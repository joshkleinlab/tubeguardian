package ahocorasick

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// loadKeywords loads banned words from a text file.
// Each line = one keyword (case-insensitive).
func LoadKeywords(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to open keywords file: %v", err)
	}
	defer file.Close()

	var keywords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			keywords = append(keywords, strings.ToLower(word))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading keywords file: %v", err)
	}

	log.Printf("Loaded %d banned keywords\n", len(keywords))
	return keywords
}
