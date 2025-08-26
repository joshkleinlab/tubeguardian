package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

// NewYouTubeService initializes a YouTube API client with automated OAuth
func NewYouTubeService(ctx context.Context, credentialsFile string) (*youtube.Service, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeForceSslScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret: %w", err)
	}

	client := getClient(ctx, config)
	return youtube.NewService(ctx, option.WithHTTPClient(client))
}

// getClient retrieves token from file or web and returns an HTTP client
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tokFile := tokenFilePath()
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(ctx, config) // AUTOMATED: opens browser and captures code
		saveToken(tokFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb starts a local HTTP server and opens the browser for OAuth
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	codeCh := make(chan string)

	// Start local HTTP server
	server := &http.Server{Addr: "localhost:8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintln(w, "✅ Authorization successful! You can close this tab.")
			codeCh <- code
		} else {
			fmt.Fprintln(w, "❌ Authorization failed.")
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Open browser automatically
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("🌐 Your browser will open for authorization. If not, visit this URL:\n%s\n", authURL)

	// Wait for code
	code := <-codeCh
	server.Shutdown(ctx)

	// Exchange code for token
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// tokenFilePath returns the path to save token.json
func tokenFilePath() string {
	dir := "configs"
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "token.json")
}

// tokenFromFile reads a token from a file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return &tok, err
}

// saveToken saves a token to a file
func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	log.Printf("💾 Token saved to %s", path)
}
