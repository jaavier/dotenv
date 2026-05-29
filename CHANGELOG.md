# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - Unreleased

This is a major release. The module path is now
`github.com/jaavier/dotenv/v2` — update your imports accordingly:

```go
import "github.com/jaavier/dotenv/v2"
```

### Added

- `Parse(io.Reader)` and `ParseBytes([]byte)` — pure parsing into a
  `map[string]string` with **no** side effects on the process environment.
- `Overload(...)` — explicit loader that lets file values override existing
  environment variables.
- `Options.MaxFileSize` and the `DefaultMaxFileSize` constant to configure the
  size limit.
- `ErrFileTooLarge` sentinel error (testable via `errors.Is`).
- Support for multi-line quoted values (e.g. PEM keys), the optional leading
  `export` token, UTF-8 BOM, and CRLF line endings.
- Testable examples (rendered on pkg.go.dev) and benchmarks.

### Changed

- **BREAKING (module path):** the module is now `github.com/jaavier/dotenv/v2`.
- **BREAKING (security default):** `Load` and `LoadWithOptions(nil, ...)` no
  longer override variables that already exist in the process environment.
  Following the 12-factor methodology, the real environment is authoritative
  and a `.env` file only fills in what is missing. Use `Overload` (or
  `Options{Override: true}`) for the previous override behavior.
- Files are now applied **atomically**: the whole file is parsed before any
  variable is set, so a malformed file never leaves a half-applied environment.
- Correct POSIX-ish quoting: double quotes expand escapes (`\n \r \t \\ \"`),
  single quotes are literal, and unquoted values support inline `# comments`.
- The size limit is enforced on the bytes actually read (safe for pipes and
  special files), and the 64 KB line cap from `bufio.Scanner` is gone.

### Fixed

- Single-quoted and unquoted values no longer have escape sequences expanded.
- Trailing garbage after a closing quote is now rejected as `ErrInvalidFormat`.

## [1.1.0] - Earlier release

- Clean-architecture refactor and utility functions (`Get`, `GetOrDefault`,
  `GetOrPanic`, `MustLoad`).

## [1.0.0] - Initial release

- Initial `.env` loader with `Load` / `LoadWithOptions`.

[2.0.0]: https://github.com/jaavier/dotenv/compare/v1.1.0...v2.0.0
[1.1.0]: https://github.com/jaavier/dotenv/releases/tag/v1.1.0
[1.0.0]: https://github.com/jaavier/dotenv/releases/tag/v1.0.0
