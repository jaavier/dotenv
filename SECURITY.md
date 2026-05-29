# Security Policy

## Supported versions

The latest released minor version receives security fixes.

| Version | Supported |
| ------- | --------- |
| 1.x     | ✅        |

## Reporting a vulnerability

Please **do not** open a public issue for security vulnerabilities.

Instead, report privately via GitHub's
[security advisories](https://github.com/jaavier/dotenv/security/advisories/new).
You can expect an initial response within a few days. Once a fix is available,
a patched release will be published and the advisory disclosed.

## Security design notes

This package is built to be safe by default:

- **No override of real environment variables.** `Load` never overwrites
  variables already present in the process environment; overriding is the
  explicit, opt-in `Overload`. This prevents a stale or accidentally-present
  `.env` from clobbering securely-injected configuration (12-factor).
- **No code execution.** Command substitution (`$(...)`) and shell evaluation
  are never performed.
- **Resource limits.** A configurable file-size cap (default 1 MiB) is enforced
  on the bytes actually read, protecting against memory-exhaustion from large
  or special files.
- **Atomic application.** A file is fully parsed before any variable is set, so
  malformed input never leaves a partially-applied environment.
- **Side-effect-free parsing.** `Parse` / `ParseBytes` never mutate global
  state, making it safe to inspect untrusted input.
