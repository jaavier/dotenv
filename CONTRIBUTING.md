# Contributing

Thanks for your interest in improving `dotenv`! This package aims to stay
**small, dependency-free, and secure**, so contributions are evaluated with
that philosophy in mind.

## Getting started

```bash
git clone https://github.com/jaavier/dotenv
cd dotenv
make test      # run the test suite
make race      # run with the race detector
make cover     # coverage report
make lint      # golangci-lint
```

If you don't have `make`, the equivalent commands are:

```bash
go test ./...
go test -race ./...
golangci-lint run ./...
gofmt -l .     # should print nothing
```

## Guidelines

- **Keep it minimal.** No third-party dependencies. The package should remain a
  single, focused file.
- **Tests first.** New behavior needs tests; bug fixes should come with a
  regression test.
- **Security by default.** Changes must not weaken the safe defaults (e.g. not
  overriding existing environment variables).
- **Formatting & linting.** `gofmt -l .` must be empty and `golangci-lint run`
  must be clean before you open a PR.
- **Document it.** Update doc comments and the README when behavior changes, and
  add an entry to `CHANGELOG.md`.

## Pull requests

1. Fork and create a feature branch.
2. Make your change with tests and docs.
3. Ensure CI is green.
4. Open a PR using the template and describe the motivation clearly.

By contributing, you agree that your contributions are licensed under the
project's [MIT License](LICENSE).
