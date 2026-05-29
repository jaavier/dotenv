package dotenv_test

import (
	"strings"
	"testing"

	"github.com/jaavier/dotenv"
)

// sampleEnv is a representative .env payload exercising comments, quotes,
// escapes and inline comments.
const sampleEnv = `# Application configuration
APP_ENV=production
APP_DEBUG=false
APP_PORT=8080            # http port

DB_HOST=localhost
DB_PORT=5432
DB_USER=app
DB_PASSWORD="s3cr3t#pass"
DB_DSN="host=localhost port=5432 dbname=app sslmode=disable"

API_KEY=AKIAIOSFODNN7EXAMPLE
API_URL=https://api.example.com/v1

MESSAGE="Hello,\tWorld!\n"
LITERAL='no $expansion or \n here'
`

func BenchmarkParse(b *testing.B) {
	data := []byte(sampleEnv)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dotenv.ParseBytes(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseReader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dotenv.Parse(strings.NewReader(sampleEnv)); err != nil {
			b.Fatal(err)
		}
	}
}
