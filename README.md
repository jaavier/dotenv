# dotenv — lightweight & secure `.env` loader for Go

> A tiny, **zero-dependency**, security-first Go library to load environment
> variables from `.env` files. A modern, actively maintained alternative to
> [godotenv](https://github.com/joho/godotenv).

[![Go Reference](https://pkg.go.dev/badge/github.com/jaavier/dotenv.svg)](https://pkg.go.dev/github.com/jaavier/dotenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaavier/dotenv)](https://goreportcard.com/report/github.com/jaavier/dotenv)
[![CI](https://github.com/jaavier/dotenv/actions/workflows/ci.yml/badge.svg)](https://github.com/jaavier/dotenv/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/jaavier/dotenv/branch/main/graph/badge.svg)](https://codecov.io/gh/jaavier/dotenv)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jaavier/dotenv)](go.mod)
[![Release](https://img.shields.io/github/v/release/jaavier/dotenv?sort=semver)](https://github.com/jaavier/dotenv/releases)
[![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Keywords: golang dotenv, load `.env` file in Go, environment variables, 12-factor
config, godotenv alternative, secure env loader.

[Español](#español) | [English](#english)

---

## English

### Features

- **Minimalist**: Zero dependencies (standard library only), single file. Loads environment variables safely - nothing more, nothing less
- **Secure by default**: Following 12-factor, the real process environment always wins — a `.env` file never overrides existing variables unless you explicitly ask it to (`Overload`)
- **Robust**: Comprehensive error handling with specific error types; files are applied atomically (a malformed file never leaves a half-applied environment)
- **Correct parsing**: POSIX-ish quoting — double quotes expand escapes, single quotes are literal, unquoted values support inline `# comments`, optional `export` prefix, multi-line quoted values
- **Side-effect free option**: `Parse`/`ParseBytes` return a `map[string]string` without ever touching the global environment — ideal for testing
- **Hardened**: Configurable file-size limit to prevent memory exhaustion; no `bufio.Scanner` 64KB line cap

### Why dotenv? (vs godotenv)

[godotenv](https://github.com/joho/godotenv) is great and battle-tested, but it
has been declared **feature-complete** and no longer accepts new functionality.
`dotenv` is a small, modern, actively-maintained alternative with safer
defaults.

| | `jaavier/dotenv` | `joho/godotenv` |
| --- | --- | --- |
| Dependencies | **0** (stdlib only) | 0 |
| Overrides existing env by default | **No** (safe; `Overload` to opt in) | No (`Load`) / Yes (`Overload`) |
| Atomic apply (all-or-nothing per file) | **Yes** | No |
| Side-effect-free `Parse` / `ParseBytes` | **Yes** | `Read` returns a map |
| Inline `# comments` (unquoted) | **Yes** | Yes |
| Single-quote literal vs double-quote escapes | **Yes** | Yes |
| Multi-line quoted values (PEM keys) | **Yes** | Yes |
| Configurable max file size (DoS guard) | **Yes** | No |
| Long lines > 64 KB | **Yes** | Limited by Scanner |
| Variable expansion `${VAR}` | No (by design — no injection surface) | Yes |
| Actively maintained | **Yes** | Feature-complete |

> `dotenv` deliberately omits `${VAR}` interpolation to keep the attack surface
> minimal. If you need interpolation, expand values yourself after `Parse`.

### Performance

Parsing a representative `.env` (15 keys, comments, quotes, escapes) on a
commodity CPU:

```
BenchmarkParse-4         217534    ~4.6 µs/op    77 MB/s    3416 B/op    24 allocs/op
```

Run it yourself: `make bench` (or `go test -bench=. -benchmem ./...`).

### Installation

```bash
go get github.com/jaavier/dotenv
```

### Quick Start

Create a `.env` file in your project root:

```env
API_KEY=your_secret_key
DB_HOST=localhost
DB_PORT=5432
```

Load it in your Go application:

```go
package main

import (
    "log"
    "os"
    
    "github.com/jaavier/dotenv"
)

func main() {
    // Load default .env file. Existing environment variables are NOT
    // overridden; the file only fills in what is missing.
    if err := dotenv.Load(); err != nil {
        log.Printf("Warning: %v", err)
    }
    
    // Use your environment variables
    apiKey := os.Getenv("API_KEY")
    dbHost := os.Getenv("DB_HOST")
}
```

> **Secure default:** `Load` never overrides variables that already exist in
> the process environment. This matches `godotenv`/12-factor: the runtime is the
> source of truth and a stale or accidental `.env` cannot clobber injected
> secrets. Use `Overload` when you explicitly want the file to win.

### Advanced Usage

#### Load Multiple Files

```go
// Load multiple files (first file has priority)
err := dotenv.Load(".env.local", ".env")
```

#### Override Existing Variables

```go
// Overload lets file values overwrite existing env vars (opt-in).
err := dotenv.Overload(".env")
```

#### Parse Without Side Effects

```go
// Parse returns a map and never mutates the global environment.
vars, err := dotenv.Parse(reader)        // any io.Reader
vars, err = dotenv.ParseBytes(data)      // or raw bytes
fmt.Println(vars["API_KEY"])
```

#### Custom Options

```go
opts := &dotenv.Options{
    Override:    false,        // default: do not override existing env vars
    Required:    true,         // file must exist
    MaxFileSize: 64 * 1024,    // optional cap in bytes (0 => DefaultMaxFileSize)
}

err := dotenv.LoadWithOptions(opts, ".env.production")
```

#### Panic on Error

```go
// Use MustLoad to panic if loading fails
dotenv.MustLoad(".env.required")
```

#### Getting Environment Variables

```go
// Simple get (same as os.Getenv)
apiKey := dotenv.Get("API_KEY")

// Get with default value
port := dotenv.GetOrDefault("PORT", "8080")

// Get required variable (panics if not set)
dbHost := dotenv.GetOrPanic("DB_HOST")
```

### Supported Formats

```env
# Full-line comments are supported
SIMPLE_KEY=value

# Inline comments (unquoted values only; '#' must follow whitespace)
PORT=8080            # the http port
PASSWORD=p#ss        # '#' kept: not preceded by whitespace

# Optional leading 'export' (so the file can also be `source`d by a shell)
export API_KEY=secret

# Double quotes: escape sequences (\n \r \t \\ \") are expanded
MULTILINE="Line 1\nLine 2"
WITH_TAB="Column1\tColumn2"

# Single quotes: fully literal, no escapes, no comment stripping
LITERAL='value with \n kept as-is and a # too'

# Multi-line quoted values (e.g. PEM keys)
PRIVATE_KEY="-----BEGIN KEY-----
line2
-----END KEY-----"

# Empty values, and surrounding whitespace is trimmed on unquoted values
EMPTY_VALUE=
TRIMMED_VALUE=  value with spaces trimmed  
```

### Error Handling

The library provides specific error types for better error handling:

```go
if err := dotenv.Load(".env"); err != nil {
    switch {
    case errors.Is(err, dotenv.ErrFileNotFound):
        // File doesn't exist
    case errors.Is(err, dotenv.ErrPermissionDenied):
        // No permission to read file
    case errors.Is(err, dotenv.ErrInvalidFormat):
        // Invalid line format
    case errors.Is(err, dotenv.ErrEmptyKey):
        // Empty key found
    case errors.Is(err, dotenv.ErrFileTooLarge):
        // File exceeds the configured size limit
    default:
        // Other error
    }
}
```

### Security Features

- **No-clobber default**: `.env` never overrides existing process env vars (12-factor); overriding is opt-in via `Overload`
- **Atomic apply**: each file is fully parsed before any variable is set, so a parse error never leaves a partially-applied environment
- **Side-effect-free parsing**: `Parse`/`ParseBytes` never mutate the global environment
- **Resource limits**: configurable file-size cap (default 1 MiB), enforced on the bytes actually read (safe for pipes/special files), with no 64KB line cap
- **No code execution**: command substitution (`$(...)`) and shell evaluation are never performed
- **Key validation**: only valid environment variable names (`[A-Za-z_][A-Za-z0-9_]*`) are accepted

### API Reference

#### Functions

**Loading Functions:**
- `Load(filenames ...string) error` - Load one or more .env files without overriding existing env vars
- `Overload(filenames ...string) error` - Load files, letting them override existing env vars
- `LoadWithOptions(opts *Options, filenames ...string) error` - Load with custom options
- `MustLoad(filenames ...string)` - Load files or panic

**Parsing (no side effects):**
- `Parse(r io.Reader) (map[string]string, error)` - Parse into a map without touching the environment
- `ParseBytes(data []byte) (map[string]string, error)` - Convenience wrapper for in-memory data

**Getting Variables:**
- `Get(key string) string` - Get environment variable value (alias for os.Getenv)
- `GetOrDefault(key, defaultValue string) string` - Get variable or return default if empty
- `GetOrPanic(key string) string` - Get variable or panic if not set/empty

#### Types

```go
type Options struct {
    Override    bool  // override existing environment variables (default false)
    Required    bool  // file must exist (return error if not found)
    MaxFileSize int64 // max bytes to read; <= 0 uses DefaultMaxFileSize
}

const DefaultMaxFileSize = 1 << 20 // 1 MiB
```

#### Errors

- `ErrFileNotFound` - File does not exist
- `ErrInvalidFormat` - Invalid line format (missing `=`, unterminated quote, or garbage after a quote)
- `ErrEmptyKey` - Empty key name
- `ErrPermissionDenied` - No permission to read file
- `ErrFileTooLarge` - File exceeds the maximum size

### FAQ

**How do I load a `.env` file in Go?**
`go get github.com/jaavier/dotenv`, then call `dotenv.Load()` at startup and read
values with `os.Getenv` (or `dotenv.Get` / `GetOrDefault`).

**Does it override my existing environment variables?**
No. `Load` only fills in variables that are not already set — the real
environment always wins. Use `dotenv.Overload(...)` if you want the file to win.

**Is it a drop-in replacement for godotenv?**
The API differs, but migration is trivial: `godotenv.Load` → `dotenv.Load`,
`godotenv.Overload` → `dotenv.Overload`, `godotenv.Read` → `dotenv.Parse`.
The main intentional difference is that `${VAR}` interpolation is not performed.

**Does it support variable expansion like `${OTHER}`?**
No, by design — this avoids an injection surface. Expand values yourself after
calling `Parse` if you need it.

**Can I parse a string or stream without touching the environment?**
Yes: `dotenv.Parse(io.Reader)` and `dotenv.ParseBytes([]byte)` return a map and
never mutate global state.

**Which Go versions are supported?**
Go 1.17 and newer (tested on Linux, macOS and Windows in CI).

---

If this package is useful to you, please consider giving it a ⭐ on
[GitHub](https://github.com/jaavier/dotenv) — it genuinely helps others discover it.

---

## Español

### Características

- **Minimalista**: Cero dependencias (solo librería estándar), un único archivo. Carga variables de entorno de manera segura - nada más, nada menos
- **Seguro por defecto**: Siguiendo 12-factor, el entorno real del proceso siempre gana — un archivo `.env` nunca sobrescribe variables existentes a menos que lo pidas explícitamente (`Overload`)
- **Robusto**: Manejo completo de errores con tipos específicos; los archivos se aplican de forma atómica (un archivo malformado nunca deja el entorno aplicado a medias)
- **Análisis correcto**: Comillas estilo POSIX — las comillas dobles expanden escapes, las simples son literales, los valores sin comillas soportan comentarios `# en línea`, prefijo `export` opcional y valores multilínea entre comillas
- **Opción sin efectos secundarios**: `Parse`/`ParseBytes` devuelven un `map[string]string` sin tocar nunca el entorno global — ideal para tests
- **Endurecido**: Límite de tamaño de archivo configurable para prevenir agotamiento de memoria; sin el límite de 64KB por línea de `bufio.Scanner`

### Instalación

```bash
go get github.com/jaavier/dotenv
```

### Inicio Rápido

Crea un archivo `.env` en la raíz de tu proyecto:

```env
API_KEY=tu_clave_secreta
DB_HOST=localhost
DB_PORT=5432
```

Cárgalo en tu aplicación Go:

```go
package main

import (
    "log"
    "os"
    
    "github.com/jaavier/dotenv"
)

func main() {
    // Cargar archivo .env por defecto. Las variables de entorno existentes
    // NO se sobrescriben; el archivo solo rellena lo que falta.
    if err := dotenv.Load(); err != nil {
        log.Printf("Advertencia: %v", err)
    }
    
    // Usar tus variables de entorno
    apiKey := os.Getenv("API_KEY")
    dbHost := os.Getenv("DB_HOST")
}
```

> **Default seguro:** `Load` nunca sobrescribe variables que ya existen en el
> entorno del proceso. Esto coincide con `godotenv`/12-factor: el runtime es la
> fuente de verdad y un `.env` accidental o desactualizado no puede pisar
> secretos inyectados. Usa `Overload` cuando quieras que el archivo gane.

### Uso Avanzado

#### Cargar Múltiples Archivos

```go
// Cargar múltiples archivos (el primer archivo tiene prioridad)
err := dotenv.Load(".env.local", ".env")
```

#### Sobrescribir Variables Existentes

```go
// Overload permite que los valores del archivo pisen las variables existentes.
err := dotenv.Overload(".env")
```

#### Analizar Sin Efectos Secundarios

```go
// Parse devuelve un map y nunca muta el entorno global.
vars, err := dotenv.Parse(reader)        // cualquier io.Reader
vars, err = dotenv.ParseBytes(data)      // o bytes en crudo
fmt.Println(vars["API_KEY"])
```

#### Opciones Personalizadas

```go
opts := &dotenv.Options{
    Override:    false,        // por defecto: no sobrescribir variables existentes
    Required:    true,         // el archivo debe existir
    MaxFileSize: 64 * 1024,    // límite opcional en bytes (0 => DefaultMaxFileSize)
}

err := dotenv.LoadWithOptions(opts, ".env.production")
```

#### Panic en Error

```go
// Usar MustLoad para hacer panic si la carga falla
dotenv.MustLoad(".env.required")
```

#### Obteniendo Variables de Entorno

```go
// Obtener simple (igual que os.Getenv)
apiKey := dotenv.Get("API_KEY")

// Obtener con valor por defecto
port := dotenv.GetOrDefault("PORT", "8080")

// Obtener variable requerida (hace panic si no está definida)
dbHost := dotenv.GetOrPanic("DB_HOST")
```

### Formatos Soportados

```env
# Comentarios de línea completa soportados
CLAVE_SIMPLE=valor

# Comentarios en línea (solo valores sin comillas; '#' debe seguir a un espacio)
PUERTO=8080            # el puerto http
PASSWORD=p#ss          # '#' conservado: no va precedido de espacio

# Prefijo 'export' opcional (para que el archivo también pueda `source`arse)
export API_KEY=secreto

# Comillas dobles: se expanden los escapes (\n \r \t \\ \")
MULTILINEA="Línea 1\nLínea 2"
CON_TAB="Columna1\tColumna2"

# Comillas simples: totalmente literal, sin escapes ni comentarios
LITERAL='valor con \n tal cual y un # también'

# Valores multilínea entre comillas (p. ej. claves PEM)
PRIVATE_KEY="-----BEGIN KEY-----
line2
-----END KEY-----"

# Valores vacíos; los espacios alrededor se recortan en valores sin comillas
VALOR_VACIO=
VALOR_LIMPIO=  valor con espacios eliminados  
```

### Manejo de Errores

La librería proporciona tipos de error específicos para un mejor manejo:

```go
if err := dotenv.Load(".env"); err != nil {
    switch {
    case errors.Is(err, dotenv.ErrFileNotFound):
        // El archivo no existe
    case errors.Is(err, dotenv.ErrPermissionDenied):
        // Sin permisos para leer el archivo
    case errors.Is(err, dotenv.ErrInvalidFormat):
        // Formato de línea inválido
    case errors.Is(err, dotenv.ErrEmptyKey):
        // Se encontró una clave vacía
    case errors.Is(err, dotenv.ErrFileTooLarge):
        // El archivo excede el límite de tamaño configurado
    default:
        // Otro error
    }
}
```

### Características de Seguridad

- **Default sin sobrescritura**: el `.env` nunca pisa variables existentes del proceso (12-factor); sobrescribir es explícito con `Overload`
- **Aplicación atómica**: cada archivo se analiza por completo antes de fijar ninguna variable, así un error de formato nunca deja el entorno a medias
- **Análisis sin efectos secundarios**: `Parse`/`ParseBytes` nunca mutan el entorno global
- **Límites de recursos**: tope de tamaño de archivo configurable (1 MiB por defecto), aplicado sobre los bytes realmente leídos (seguro para pipes/archivos especiales), sin límite de 64KB por línea
- **Sin ejecución de código**: nunca se realiza sustitución de comandos (`$(...)`) ni evaluación de shell
- **Validación de claves**: solo se aceptan nombres válidos de variables (`[A-Za-z_][A-Za-z0-9_]*`)

### Referencia API

#### Funciones

**Funciones de Carga:**
- `Load(filenames ...string) error` - Cargar uno o más archivos .env sin sobrescribir variables existentes
- `Overload(filenames ...string) error` - Cargar archivos dejando que sobrescriban las variables existentes
- `LoadWithOptions(opts *Options, filenames ...string) error` - Cargar con opciones personalizadas
- `MustLoad(filenames ...string)` - Cargar archivos o hacer panic

**Análisis (sin efectos secundarios):**
- `Parse(r io.Reader) (map[string]string, error)` - Analizar a un map sin tocar el entorno
- `ParseBytes(data []byte) (map[string]string, error)` - Atajo para datos en memoria

**Funciones para Obtener Variables:**
- `Get(key string) string` - Obtener valor de variable de entorno (alias de os.Getenv)
- `GetOrDefault(key, defaultValue string) string` - Obtener variable o retornar valor por defecto si está vacía
- `GetOrPanic(key string) string` - Obtener variable o hacer panic si no está definida/vacía

#### Tipos

```go
type Options struct {
    Override    bool  // sobrescribir variables existentes (por defecto false)
    Required    bool  // el archivo debe existir (retorna error si no)
    MaxFileSize int64 // máximo de bytes a leer; <= 0 usa DefaultMaxFileSize
}

const DefaultMaxFileSize = 1 << 20 // 1 MiB
```

#### Errores

- `ErrFileNotFound` - El archivo no existe
- `ErrInvalidFormat` - Formato inválido (falta `=`, comilla sin cerrar o basura tras una comilla)
- `ErrEmptyKey` - Nombre de clave vacío
- `ErrPermissionDenied` - Sin permisos para leer el archivo
- `ErrFileTooLarge` - El archivo excede el tamaño máximo

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

[@jaavier](https://github.com/jaavier)