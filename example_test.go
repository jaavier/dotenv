package dotenv_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaavier/dotenv"
)

// Load reads a .env file and populates the process environment without
// overriding variables that are already set.
func ExampleLoad() {
	dir, _ := os.MkdirTemp("", "dotenv")
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("GREETING=hello\nPORT=8080"), 0o600)

	if err := dotenv.Load(path); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(os.Getenv("GREETING"))
	fmt.Println(os.Getenv("PORT"))
	// Output:
	// hello
	// 8080
}

// Parse reads key/value pairs into a map without ever touching the process
// environment, which makes it ideal for tests and validation.
func ExampleParse() {
	r := strings.NewReader(`
# comments and blank lines are ignored
HOST=localhost
PORT=5432           # inline comments are stripped
LITERAL='no \n expansion here'
ESCAPED="line1\nline2"
`)
	vars, err := dotenv.Parse(r)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%s:%s\n", vars["HOST"], vars["PORT"])
	fmt.Println(vars["LITERAL"])
	fmt.Printf("%q\n", vars["ESCAPED"])
	// Output:
	// localhost:5432
	// no \n expansion here
	// "line1\nline2"
}

// ParseBytes is a convenience wrapper for in-memory data.
func ExampleParseBytes() {
	vars, _ := dotenv.ParseBytes([]byte("API_KEY=secret123"))
	fmt.Println(vars["API_KEY"])
	// Output: secret123
}

// Overload lets values from the file override variables that already exist in
// the environment. The default Load never overrides; Overload is the explicit
// opt-in.
func ExampleOverload() {
	os.Setenv("REGION", "us-east-1")
	defer os.Unsetenv("REGION")

	dir, _ := os.MkdirTemp("", "dotenv")
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("REGION=eu-west-1"), 0o600)

	_ = dotenv.Overload(path)
	fmt.Println(os.Getenv("REGION"))
	// Output: eu-west-1
}

// GetOrDefault returns a fallback when a variable is unset or empty.
func ExampleGetOrDefault() {
	os.Unsetenv("MISSING_VAR")
	fmt.Println(dotenv.GetOrDefault("MISSING_VAR", "fallback"))
	// Output: fallback
}
