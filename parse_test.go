package dotenv

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse_PureNoSideEffects(t *testing.T) {
	const key = "DOTENV_PARSE_SENTINEL"
	os.Setenv(key, "original")
	defer os.Unsetenv(key)

	vars, err := ParseBytes([]byte(key + "=fromfile\nOTHER=x"))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	if vars[key] != "fromfile" {
		t.Errorf("parsed %s = %q, want %q", key, vars[key], "fromfile")
	}
	if got := os.Getenv(key); got != "original" {
		t.Errorf("Parse must not mutate env: %s = %q, want %q", key, got, "original")
	}
}

func TestParse_Quoting(t *testing.T) {
	tests := []struct {
		name string
		in   string
		key  string
		want string
	}{
		{"single quote literal", `S='a\nb'`, "S", `a\nb`},
		{"double quote escapes", `D="a\nb"`, "D", "a\nb"},
		{"escaped quote", `Q="he said \"hi\"\n"`, "Q", "he said \"hi\"\n"},
		{"hash in double quotes", `H="a#b"`, "H", "a#b"},
		{"hash in single quotes", `H='a#b'`, "H", "a#b"},
		{"inline comment", "K=val # comment", "K", "val"},
		{"hash not a comment", "K=pa#ss", "K", "pa#ss"},
		{"equals in value", "K=a=b=c", "K", "a=b=c"},
		{"empty value", "K=", "K", ""},
		{"whitespace value", "K=    ", "K", ""},
		{"export prefix", "export FOO=bar", "FOO", "bar"},
		{"export with quotes", `export   BAZ="v"`, "BAZ", "v"},
		{"tab escape", `T="c1\tc2"`, "T", "c1\tc2"},
		{"backslash literal", `B="a\\b"`, "B", `a\b`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars, err := ParseBytes([]byte(tt.in))
			if err != nil {
				t.Fatalf("ParseBytes(%q) error = %v", tt.in, err)
			}
			if got := vars[tt.key]; got != tt.want {
				t.Errorf("ParseBytes(%q)[%q] = %q, want %q", tt.in, tt.key, got, tt.want)
			}
		})
	}
}

func TestParse_TrailingGarbageAfterQuoteRejected(t *testing.T) {
	tests := []string{
		`K="a"extra`,
		`K="a"="b"`,
		`K='a'bcd`,
		`K="a" b`,
	}
	for _, in := range tests {
		if _, err := ParseBytes([]byte(in)); !errors.Is(err, ErrInvalidFormat) {
			t.Errorf("ParseBytes(%q) error = %v, want ErrInvalidFormat", in, err)
		}
	}
}

func TestParse_TrailingCommentAfterQuoteAllowed(t *testing.T) {
	vars, err := ParseBytes([]byte(`K="a" # trailing comment`))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}
	if vars["K"] != "a" {
		t.Errorf("K = %q, want %q", vars["K"], "a")
	}
}

func TestParse_DuplicateKeysLastWins(t *testing.T) {
	vars, err := ParseBytes([]byte("K=1\nK=2\nK=3"))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}
	if vars["K"] != "3" {
		t.Errorf("duplicate key = %q, want %q", vars["K"], "3")
	}
}

func TestParse_MultilineQuotedValue(t *testing.T) {
	content := "CERT=\"line1\nline2\nline3\"\nNEXT=ok"
	vars, err := ParseBytes([]byte(content))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}
	if want := "line1\nline2\nline3"; vars["CERT"] != want {
		t.Errorf("CERT = %q, want %q", vars["CERT"], want)
	}
	if vars["NEXT"] != "ok" {
		t.Errorf("NEXT = %q, want %q", vars["NEXT"], "ok")
	}
}

func TestParse_CRLFAndBOM(t *testing.T) {
	content := "\xEF\xBB\xBFA=1\r\nB=2\r\n"
	vars, err := ParseBytes([]byte(content))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}
	if vars["A"] != "1" || vars["B"] != "2" {
		t.Errorf("CRLF/BOM parse = %v, want A=1 B=2", vars)
	}
}

func TestParse_UnterminatedQuote(t *testing.T) {
	for _, in := range []string{`K="oops`, `K='oops`} {
		_, err := ParseBytes([]byte(in))
		if !errors.Is(err, ErrInvalidFormat) {
			t.Errorf("ParseBytes(%q) error = %v, want ErrInvalidFormat", in, err)
		}
	}
}

func TestParse_LongValueNo64KCap(t *testing.T) {
	// bufio.Scanner would choke past 64KiB; the byte-based parser must not.
	long := strings.Repeat("x", 200_000)
	vars, err := ParseBytes([]byte("BIG=" + long))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}
	if len(vars["BIG"]) != len(long) {
		t.Errorf("BIG length = %d, want %d", len(vars["BIG"]), len(long))
	}
}

func TestLoad_DoesNotOverrideByDefault(t *testing.T) {
	const key = "DOTENV_NO_OVERRIDE"
	os.Setenv(key, "real")
	defer os.Unsetenv(key)

	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmpFile, []byte(key+"=fromfile"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Load(tmpFile); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got := os.Getenv(key); got != "real" {
		t.Errorf("default Load must not override: %s = %q, want %q", key, got, "real")
	}
}

func TestOverload_OverridesExisting(t *testing.T) {
	const key = "DOTENV_OVERLOAD"
	os.Setenv(key, "real")
	defer os.Unsetenv(key)

	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmpFile, []byte(key+"=fromfile"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Overload(tmpFile); err != nil {
		t.Fatalf("Overload() error = %v", err)
	}
	if got := os.Getenv(key); got != "fromfile" {
		t.Errorf("Overload must override: %s = %q, want %q", key, got, "fromfile")
	}
}

func TestLoad_AtomicOnError(t *testing.T) {
	const good = "DOTENV_ATOMIC_GOOD"
	os.Unsetenv(good)
	defer os.Unsetenv(good)

	// Second line is malformed; nothing from this file should be applied.
	content := good + "=value\n123BAD=value"
	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmpFile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Load(tmpFile); err == nil {
		t.Fatal("Load() should fail on malformed file")
	}
	if _, ok := os.LookupEnv(good); ok {
		t.Errorf("no variable should be applied when the file fails to parse")
	}
}

func TestLoad_FileTooLarge(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmpFile, []byte("K=value123"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := LoadWithOptions(&Options{MaxFileSize: 4}, tmpFile)
	if !errors.Is(err, ErrFileTooLarge) {
		t.Errorf("LoadWithOptions() error = %v, want ErrFileTooLarge", err)
	}
}

func TestLoad_MissingFileNotRequired(t *testing.T) {
	if err := Load(filepath.Join(t.TempDir(), "does-not-exist.env")); err != nil {
		t.Errorf("missing file with default options should be ignored, got %v", err)
	}
}
