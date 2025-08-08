package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("loads basic env file", func(t *testing.T) {
		content := `TEST_KEY=test_value
DB_HOST=localhost
DB_PORT=5432`
		
		tmpFile := filepath.Join(t.TempDir(), ".env")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		
		if err := Load(tmpFile); err != nil {
			t.Errorf("Load() error = %v", err)
		}
		
		if got := os.Getenv("TEST_KEY"); got != "test_value" {
			t.Errorf("TEST_KEY = %v, want %v", got, "test_value")
		}
	})
	
	t.Run("handles quoted values", func(t *testing.T) {
		content := `QUOTED="value with spaces"
SINGLE_QUOTED='single quotes'`
		
		tmpFile := filepath.Join(t.TempDir(), ".env")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		
		if err := Load(tmpFile); err != nil {
			t.Errorf("Load() error = %v", err)
		}
		
		if got := os.Getenv("QUOTED"); got != "value with spaces" {
			t.Errorf("QUOTED = %v, want %v", got, "value with spaces")
		}
		
		if got := os.Getenv("SINGLE_QUOTED"); got != "single quotes" {
			t.Errorf("SINGLE_QUOTED = %v, want %v", got, "single quotes")
		}
	})
	
	t.Run("ignores comments and empty lines", func(t *testing.T) {
		content := `# This is a comment
KEY1=value1

# Another comment
KEY2=value2
`
		
		tmpFile := filepath.Join(t.TempDir(), ".env")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		
		if err := Load(tmpFile); err != nil {
			t.Errorf("Load() error = %v", err)
		}
		
		if got := os.Getenv("KEY1"); got != "value1" {
			t.Errorf("KEY1 = %v, want %v", got, "value1")
		}
		
		if got := os.Getenv("KEY2"); got != "value2" {
			t.Errorf("KEY2 = %v, want %v", got, "value2")
		}
	})
	
	t.Run("handles escape sequences", func(t *testing.T) {
		content := `MULTILINE="Line1\nLine2"
WITH_TAB="Col1\tCol2"`
		
		tmpFile := filepath.Join(t.TempDir(), ".env")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		
		if err := Load(tmpFile); err != nil {
			t.Errorf("Load() error = %v", err)
		}
		
		if got := os.Getenv("MULTILINE"); got != "Line1\nLine2" {
			t.Errorf("MULTILINE = %v, want %v", got, "Line1\nLine2")
		}
		
		if got := os.Getenv("WITH_TAB"); got != "Col1\tCol2" {
			t.Errorf("WITH_TAB = %v, want %v", got, "Col1\tCol2")
		}
	})
	
	t.Run("returns error for non-existent file with Required option", func(t *testing.T) {
		opts := &Options{Required: true}
		err := LoadWithOptions(opts, "/non/existent/file.env")
		if err == nil {
			t.Error("LoadWithOptions() should return error for non-existent file when Required=true")
		}
	})
	
	t.Run("validates key format", func(t *testing.T) {
		content := `123INVALID=value
VALID_KEY=value`
		
		tmpFile := filepath.Join(t.TempDir(), ".env")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		
		err := Load(tmpFile)
		if err == nil {
			t.Error("Load() should return error for invalid key format")
		}
	})
}

func TestLoadWithOptions(t *testing.T) {
	t.Run("override option", func(t *testing.T) {
		os.Setenv("EXISTING_KEY", "original")
		defer os.Unsetenv("EXISTING_KEY")
		
		content := `EXISTING_KEY=new_value`
		tmpFile := filepath.Join(t.TempDir(), ".env")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		
		opts := &Options{Override: false}
		if err := LoadWithOptions(opts, tmpFile); err != nil {
			t.Errorf("LoadWithOptions() error = %v", err)
		}
		
		if got := os.Getenv("EXISTING_KEY"); got != "original" {
			t.Errorf("EXISTING_KEY should not be overridden, got %v", got)
		}
		
		opts.Override = true
		if err := LoadWithOptions(opts, tmpFile); err != nil {
			t.Errorf("LoadWithOptions() error = %v", err)
		}
		
		if got := os.Getenv("EXISTING_KEY"); got != "new_value" {
			t.Errorf("EXISTING_KEY should be overridden, got %v", got)
		}
	})
	
	t.Run("required option", func(t *testing.T) {
		opts := &Options{Required: true}
		err := LoadWithOptions(opts, "/non/existent/file.env")
		if err == nil {
			t.Error("LoadWithOptions() should return error when Required=true and file doesn't exist")
		}
	})
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{"basic", "KEY=value", "KEY", "value", false},
		{"empty value", "KEY=", "KEY", "", false},
		{"spaces", "KEY = value ", "KEY", "value", false},
		{"quoted", `KEY="value"`, "KEY", "value", false},
		{"no equals", "KEY", "", "", true},
		{"empty key", "=value", "", "", true},
		{"valid underscore", "MY_KEY=value", "MY_KEY", "value", false},
		{"number in key", "KEY123=value", "KEY123", "value", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, err := parseLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("parseLine() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("parseLine() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestIsValidKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"VALID_KEY", true},
		{"_UNDERSCORE", true},
		{"lowercase", true},
		{"MixedCase", true},
		{"WITH_123", true},
		{"", false},
		{"123START", false},
		{"WITH-DASH", false},
		{"WITH.DOT", false},
		{"WITH SPACE", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := isValidKey(tt.key); got != tt.valid {
				t.Errorf("isValidKey(%v) = %v, want %v", tt.key, got, tt.valid)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_GET_VAR", "test_value")
	defer os.Unsetenv("TEST_GET_VAR")
	
	if got := Get("TEST_GET_VAR"); got != "test_value" {
		t.Errorf("Get() = %v, want %v", got, "test_value")
	}
	
	if got := Get("NON_EXISTENT_VAR"); got != "" {
		t.Errorf("Get() for non-existent var = %v, want empty string", got)
	}
}

func TestGetOrDefault(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_DEFAULT_VAR", "actual_value")
	defer os.Unsetenv("TEST_DEFAULT_VAR")
	
	tests := []struct {
		name         string
		key          string
		defaultValue string
		want         string
	}{
		{"existing var", "TEST_DEFAULT_VAR", "default", "actual_value"},
		{"non-existent var", "NON_EXISTENT_VAR", "default", "default"},
		{"empty var", "EMPTY_VAR", "default", "default"},
	}
	
	// Set empty var for testing
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOrDefault(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOrPanic(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_PANIC_VAR", "safe_value")
	defer os.Unsetenv("TEST_PANIC_VAR")
	
	t.Run("returns value for existing var", func(t *testing.T) {
		if got := GetOrPanic("TEST_PANIC_VAR"); got != "safe_value" {
			t.Errorf("GetOrPanic() = %v, want %v", got, "safe_value")
		}
	})
	
	t.Run("panics for non-existent var", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("GetOrPanic() should have panicked for non-existent var")
			} else {
				expected := "dotenv: required environment variable NON_EXISTENT_PANIC_VAR is not set or is empty"
				if r != expected {
					t.Errorf("GetOrPanic() panic message = %v, want %v", r, expected)
				}
			}
		}()
		GetOrPanic("NON_EXISTENT_PANIC_VAR")
	})
	
	t.Run("panics for empty var", func(t *testing.T) {
		os.Setenv("EMPTY_PANIC_VAR", "")
		defer os.Unsetenv("EMPTY_PANIC_VAR")
		
		defer func() {
			if r := recover(); r == nil {
				t.Error("GetOrPanic() should have panicked for empty var")
			}
		}()
		GetOrPanic("EMPTY_PANIC_VAR")
	})
}