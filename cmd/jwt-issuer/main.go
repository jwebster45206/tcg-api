package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jwebster45206/tcg-api/internal/auth"
	"github.com/jwebster45206/tcg-api/internal/config"
)

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}

func loadConfig(configPath string) (*config.Config, error) {
	expandedPath, err := expandPath(configPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(expandedPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	// Parse config path from first argument, default to "config.json"
	configPath := "config.json"
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		configPath = os.Args[1]
		// Remove the config path from args so flag parsing works correctly
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	sub := flag.String("sub", "admin", "subject (user id)")
	scopes := flag.String("scopes", "admin", "comma-separated scopes")
	roles := flag.String("roles", "admin", "comma-separated roles")
	ttl := flag.Duration("ttl", time.Hour, "token ttl (e.g., 1h)")
	flag.Parse()

	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	issuer := auth.HS256IssuerVerifier{
		Issuer:     cfg.Auth.Issuer,
		Audience:   cfg.Auth.Audience,
		CurrentKID: cfg.Auth.CurrentKID,
		Keys:       cfg.Auth.Keys,
	}

	tok, claims, err := issuer.Issue(*sub, splitList(*scopes), *ttl, splitList(*roles))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to issue token: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", tok)
	_ = claims
}

// ...existing code...
func splitList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
