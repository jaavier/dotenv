// Package dotenv provides a minimalist, secure, and robust solution for loading
// environment variables from .env files.
//
// The package focuses on doing one thing with excellence: safely loading environment
// variables from files. It provides comprehensive error handling, security features
// like file size limits and path sanitization, and flexible configuration options.
//
// Basic usage:
//
//	import "github.com/jaavier/dotenv"
//
//	func main() {
//	    // Load default .env file
//	    if err := dotenv.Load(); err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// For more advanced usage, see the README documentation.
package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrFileNotFound    = errors.New("dotenv: file not found")
	ErrInvalidFormat   = errors.New("dotenv: invalid format")
	ErrEmptyKey        = errors.New("dotenv: empty key")
	ErrPermissionDenied = errors.New("dotenv: permission denied")
)

const (
	defaultFile = ".env"
	maxFileSize = 1 << 20
)

type Options struct {
	Override bool
	Required bool
}

func Load(filenames ...string) error {
	return LoadWithOptions(nil, filenames...)
}

func LoadWithOptions(opts *Options, filenames ...string) error {
	if opts == nil {
		opts = &Options{
			Override: true,
			Required: false,
		}
	}

	if len(filenames) == 0 {
		filenames = []string{defaultFile}
	}

	for _, filename := range filenames {
		if err := loadFile(filename, opts); err != nil {
			if opts.Required || !errors.Is(err, ErrFileNotFound) {
				return fmt.Errorf("failed to load %s: %w", filename, err)
			}
		}
	}

	return nil
}

func loadFile(filename string, opts *Options) error {
	filename = filepath.Clean(filename)

	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return err
	}

	if info.Size() > maxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxFileSize)
	}

	if info.IsDir() {
		return fmt.Errorf("expected file, got directory: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return err
	}
	defer file.Close()

	return parse(file, opts)
}

func parse(file *os.File, opts *Options) error {
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}

		if len(key) > 0 {
			if opts.Override || os.Getenv(key) == "" {
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("failed to set environment variable %s: %w", key, err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func parseLine(line string) (string, string, error) {
	if len(line) == 0 {
		return "", "", nil
	}

	equalIndex := strings.Index(line, "=")
	if equalIndex == -1 {
		return "", "", ErrInvalidFormat
	}

	key := strings.TrimSpace(line[:equalIndex])
	if len(key) == 0 {
		return "", "", ErrEmptyKey
	}

	if !isValidKey(key) {
		return "", "", fmt.Errorf("invalid key format: %s", key)
	}

	value := ""
	if equalIndex < len(line)-1 {
		value = line[equalIndex+1:]
		value = processValue(value)
	}

	return key, value, nil
}

func processValue(value string) string {
	value = strings.TrimSpace(value)

	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	value = expandEscapes(value)

	return value
}

func expandEscapes(value string) string {
	replacer := strings.NewReplacer(
		`\n`, "\n",
		`\r`, "\r",
		`\t`, "\t",
		`\\`, "\\",
	)
	return replacer.Replace(value)
}

func isValidKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	for i, ch := range key {
		if !((ch >= 'A' && ch <= 'Z') || 
			(ch >= 'a' && ch <= 'z') || 
			(ch >= '0' && ch <= '9' && i > 0) || 
			ch == '_') {
			return false
		}
	}
	
	return true
}

func MustLoad(filenames ...string) {
	if err := Load(filenames...); err != nil {
		panic(err)
	}
}

// Get is a convenience wrapper around os.Getenv for consistent API.
// It returns the value of the environment variable named by the key.
func Get(key string) string {
	return os.Getenv(key)
}

// GetOrDefault returns the value of the environment variable named by the key,
// or returns defaultValue if the variable is not present or is empty.
func GetOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetOrPanic returns the value of the environment variable named by the key.
// It panics if the variable is not present or is empty.
// This is useful for required configuration that the application cannot run without.
func GetOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("dotenv: required environment variable %s is not set or is empty", key))
	}
	return value
}