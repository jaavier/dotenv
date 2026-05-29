// Package dotenv is a lightweight, dependency-free, secure loader for .env
// files in Go — a modern, actively maintained alternative to godotenv.
//
// The package focuses on doing one thing with excellence: safely loading
// environment variables from files. It uses only the standard library, fits in
// a single file, and is secure by default.
//
// # Secure by default
//
// Following the 12-factor methodology, the process environment is always
// authoritative: Load does NOT override variables that already exist in the
// environment. A .env file only fills in the gaps. If you explicitly want the
// file to win, use Overload (or Options.Override).
//
// Basic usage:
//
//	import "github.com/jaavier/dotenv/v2"
//
//	func main() {
//	    // Load default .env file without overriding existing env vars.
//	    if err := dotenv.Load(); err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// Parsing without side effects:
//
//	vars, err := dotenv.Parse(reader) // returns a map, never touches os.Environ
//
// For more advanced usage, see the README documentation.
package dotenv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Sentinel errors returned by the package. Use errors.Is to test for them.
var (
	ErrFileNotFound     = errors.New("dotenv: file not found")
	ErrInvalidFormat    = errors.New("dotenv: invalid format")
	ErrEmptyKey         = errors.New("dotenv: empty key")
	ErrPermissionDenied = errors.New("dotenv: permission denied")
	ErrFileTooLarge     = errors.New("dotenv: file exceeds maximum size")
)

const (
	defaultFile = ".env"

	// DefaultMaxFileSize is the maximum size, in bytes, of a .env file that
	// will be read when Options.MaxFileSize is not set. It guards against
	// resource-exhaustion from pathologically large files.
	DefaultMaxFileSize = 1 << 20 // 1 MiB
)

// utf8BOM is stripped from the start of a file if present.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// Options configures how environment files are loaded.
//
// The zero value is the recommended, secure configuration: it does not
// override existing environment variables, treats missing files as
// non-fatal, and uses DefaultMaxFileSize.
type Options struct {
	// Override, when true, lets values from the file overwrite variables that
	// already exist in the process environment. The default (false) is the
	// secure choice: the real environment always wins.
	Override bool

	// Required, when true, causes loading to fail if a file does not exist.
	Required bool

	// MaxFileSize is the maximum size, in bytes, of a file to read. Values
	// <= 0 fall back to DefaultMaxFileSize.
	MaxFileSize int64
}

// Load reads one or more .env files (default ".env") and sets the contained
// variables in the process environment WITHOUT overriding variables that are
// already set. Missing files are ignored. See Overload to override.
func Load(filenames ...string) error {
	return LoadWithOptions(nil, filenames...)
}

// Overload is like Load but lets file values override variables that already
// exist in the environment. Use it only when the file is the source of truth.
func Overload(filenames ...string) error {
	return LoadWithOptions(&Options{Override: true}, filenames...)
}

// LoadWithOptions reads one or more .env files using the supplied Options.
// A nil opts means the secure zero-value configuration. Each file is parsed
// in full before any variable is applied, so a malformed file never leaves a
// partially-applied environment.
func LoadWithOptions(opts *Options, filenames ...string) error {
	if opts == nil {
		opts = &Options{}
	}

	if len(filenames) == 0 {
		filenames = []string{defaultFile}
	}

	for _, filename := range filenames {
		vars, err := loadFile(filename, opts)
		if err != nil {
			if opts.Required || !errors.Is(err, ErrFileNotFound) {
				return fmt.Errorf("dotenv: failed to load %s: %w", filename, err)
			}
			continue
		}
		if err := applyVars(vars, opts.Override); err != nil {
			return fmt.Errorf("dotenv: failed to apply %s: %w", filename, err)
		}
	}

	return nil
}

// MustLoad behaves like Load but panics on error. Use it only at program
// startup for configuration the application cannot run without.
func MustLoad(filenames ...string) {
	if err := Load(filenames...); err != nil {
		panic(err)
	}
}

// Parse reads key/value pairs from r and returns them as a map. It performs
// no I/O beyond reading r and never mutates the process environment, which
// makes it ideal for testing and for inspecting values before applying them.
// Reads are capped at DefaultMaxFileSize.
func Parse(r io.Reader) (map[string]string, error) {
	return parseReader(r, DefaultMaxFileSize)
}

// ParseBytes is a convenience wrapper around Parse for in-memory data.
func ParseBytes(data []byte) (map[string]string, error) {
	return parseReader(bytes.NewReader(data), DefaultMaxFileSize)
}

// applyVars writes vars to the process environment. When override is false a
// variable is only set if it is not already present (an existing empty value
// still counts as present and is preserved).
func applyVars(vars map[string]string, override bool) error {
	for k, v := range vars {
		if !override {
			if _, ok := os.LookupEnv(k); ok {
				continue
			}
		}
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", k, err)
		}
	}
	return nil
}

// loadFile validates and reads a single file into a map of key/value pairs.
func loadFile(filename string, opts *Options) (map[string]string, error) {
	filename = filepath.Clean(filename)

	info, err := os.Stat(filename)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return nil, ErrFileNotFound
		case os.IsPermission(err):
			return nil, ErrPermissionDenied
		default:
			return nil, err
		}
	}

	if info.IsDir() {
		return nil, fmt.Errorf("expected file, got directory: %s", filename)
	}

	maxSize := opts.MaxFileSize
	if maxSize <= 0 {
		maxSize = DefaultMaxFileSize
	}
	if info.Size() > maxSize {
		return nil, ErrFileTooLarge
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsPermission(err) {
			return nil, ErrPermissionDenied
		}
		return nil, err
	}
	defer func() { _ = file.Close() }()

	return parseReader(file, maxSize)
}

// parseReader reads at most maxSize bytes from r and parses them. The limit is
// enforced on the actual bytes read (not just a reported size), so it is safe
// for pipes, sockets and special files where Stat lies about the size.
func parseReader(r io.Reader, maxSize int64) (map[string]string, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxFileSize
	}

	data, err := io.ReadAll(io.LimitReader(r, maxSize+1))
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	if int64(len(data)) > maxSize {
		return nil, ErrFileTooLarge
	}

	data = bytes.TrimPrefix(data, utf8BOM)
	lines := strings.Split(string(data), "\n")
	out := make(map[string]string)

	for i := 0; i < len(lines); i++ {
		raw := strings.TrimRight(lines[i], "\r")
		trimmed := strings.TrimSpace(raw)

		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// A quoted value may span multiple physical lines (e.g. a PEM key).
		// Accumulate following lines until the quote is closed. Errors are
		// reported against the line where the assignment begins.
		start := i
		logical := raw
		for needsMoreLines(logical) {
			i++
			if i >= len(lines) {
				return nil, fmt.Errorf("line %d: %w", start+1, ErrInvalidFormat)
			}
			logical += "\n" + strings.TrimRight(lines[i], "\r")
		}

		key, value, err := parseLine(logical)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", start+1, err)
		}
		if len(key) > 0 {
			out[key] = value
		}
	}

	return out, nil
}

// needsMoreLines reports whether logical opens a quoted value that has not yet
// been closed, meaning more physical lines must be appended.
func needsMoreLines(logical string) bool {
	line := strings.TrimSpace(logical)
	if rest, ok := cutExport(line); ok {
		line = rest
	}

	eq := strings.IndexByte(line, '=')
	if eq < 0 {
		return false
	}

	v := strings.TrimLeft(line[eq+1:], " \t")
	if len(v) == 0 {
		return false
	}

	q := v[0]
	if q != '"' && q != '\'' {
		return false
	}
	return indexClosingQuote(v[1:], q) < 0
}

// parseLine parses a single logical line into a key/value pair. An empty or
// comment-only line yields ("", "", nil). The leading "export " token, if
// present, is ignored.
func parseLine(line string) (string, string, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", "", nil
	}

	if rest, ok := cutExport(line); ok {
		line = rest
	}

	eq := strings.IndexByte(line, '=')
	if eq == -1 {
		return "", "", ErrInvalidFormat
	}

	key := strings.TrimSpace(line[:eq])
	if len(key) == 0 {
		return "", "", ErrEmptyKey
	}
	if !isValidKey(key) {
		return "", "", fmt.Errorf("invalid key format: %s", key)
	}

	value, err := parseValue(line[eq+1:])
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

// parseValue interprets the raw right-hand side of an assignment following
// POSIX-ish conventions:
//   - double-quoted values expand escape sequences (\n \r \t \\ \");
//   - single-quoted values are taken literally (no escapes, no comments);
//   - unquoted values are trimmed and have inline "# comments" removed.
func parseValue(raw string) (string, error) {
	v := strings.TrimLeft(raw, " \t")
	if len(v) == 0 {
		return "", nil
	}

	switch v[0] {
	case '"':
		idx := indexClosingQuote(v[1:], '"')
		if idx < 0 || !validTrailer(v[2+idx:]) {
			return "", ErrInvalidFormat // unterminated quote or trailing garbage
		}
		return expandEscapes(v[1 : 1+idx]), nil
	case '\'':
		idx := indexClosingQuote(v[1:], '\'')
		if idx < 0 || !validTrailer(v[2+idx:]) {
			return "", ErrInvalidFormat // unterminated quote or trailing garbage
		}
		return v[1 : 1+idx], nil
	default:
		return stripInlineComment(v), nil
	}
}

// validTrailer reports whether the text following a closing quote is allowed:
// only optional whitespace and an optional inline "# comment".
func validTrailer(s string) bool {
	s = strings.TrimLeft(s, " \t")
	return len(s) == 0 || s[0] == '#'
}

// indexClosingQuote returns the index of the closing quote in s, or -1 if it
// is not found. For double quotes a backslash escapes the next character;
// single quotes are literal and the first quote closes.
func indexClosingQuote(s string, quote byte) int {
	escaped := false
	for i := 0; i < len(s); i++ {
		if escaped {
			escaped = false
			continue
		}
		c := s[i]
		if quote == '"' && c == '\\' {
			escaped = true
			continue
		}
		if c == quote {
			return i
		}
	}
	return -1
}

// stripInlineComment removes a trailing "# comment" from an unquoted value.
// The '#' is only treated as a comment when it starts the value or is preceded
// by whitespace, so values like "pa#ss" are preserved. Trailing whitespace is
// trimmed.
func stripInlineComment(v string) string {
	for i := 0; i < len(v); i++ {
		if v[i] == '#' && (i == 0 || v[i-1] == ' ' || v[i-1] == '\t') {
			v = v[:i]
			break
		}
	}
	return strings.TrimRight(v, " \t")
}

// cutExport strips an optional leading "export " (or "export\t") token,
// allowing files that can also be sourced by a shell.
func cutExport(line string) (string, bool) {
	const kw = "export"
	if strings.HasPrefix(line, kw) && len(line) > len(kw) &&
		(line[len(kw)] == ' ' || line[len(kw)] == '\t') {
		return strings.TrimLeft(line[len(kw):], " \t"), true
	}
	return line, false
}

// expandEscapes interprets the escape sequences allowed inside double quotes.
func expandEscapes(value string) string {
	if !strings.ContainsRune(value, '\\') {
		return value
	}
	replacer := strings.NewReplacer(
		`\n`, "\n",
		`\r`, "\r",
		`\t`, "\t",
		`\"`, `"`,
		`\\`, `\`,
	)
	return replacer.Replace(value)
}

// isValidKey reports whether key matches [A-Za-z_][A-Za-z0-9_]*.
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
// It panics if the variable is not present or is empty. This is useful for
// required configuration that the application cannot run without.
func GetOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("dotenv: required environment variable %s is not set or is empty", key))
	}
	return value
}
