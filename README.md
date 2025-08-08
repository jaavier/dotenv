# dotenv

A minimalist, secure, and robust Go library for loading environment variables from `.env` files. Built with clean architecture principles, focusing on doing one thing with excellence.

[Español](#español) | [English](#english)

---

## English

### Features

- **Minimalist**: Focus on loading environment variables safely - nothing more, nothing less
- **Secure**: File size limits, path sanitization, and permission checks
- **Robust**: Comprehensive error handling with specific error types
- **Clean**: Professional code structure following clean architecture principles
- **Flexible**: Support for multiple files and configuration options
- **Safe**: Input validation and secure parsing

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
    // Load default .env file
    if err := dotenv.Load(); err != nil {
        log.Printf("Warning: %v", err)
    }
    
    // Use your environment variables
    apiKey := os.Getenv("API_KEY")
    dbHost := os.Getenv("DB_HOST")
}
```

### Advanced Usage

#### Load Multiple Files

```go
// Load multiple files (first file has priority)
err := dotenv.Load(".env.local", ".env")
```

#### Custom Options

```go
opts := &dotenv.Options{
    Override: false,  // Don't override existing env vars
    Required: true,   // File must exist
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
# Comments are supported
SIMPLE_KEY=value

# Quoted values
QUOTED_VALUE="value with spaces"
SINGLE_QUOTED='single quotes also work'

# Empty values
EMPTY_VALUE=

# Special characters (escaped)
MULTILINE="Line 1\nLine 2"
WITH_TAB="Column1\tColumn2"

# Trim spaces
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
    default:
        // Other error
    }
}
```

### Security Features

- **File size limit**: Maximum 1MB to prevent memory exhaustion
- **Path sanitization**: Prevents directory traversal attacks
- **Permission checks**: Validates file permissions before reading
- **Key validation**: Only allows valid environment variable names
- **Safe parsing**: Handles malformed input gracefully

### API Reference

#### Functions

**Loading Functions:**
- `Load(filenames ...string) error` - Load one or more .env files
- `LoadWithOptions(opts *Options, filenames ...string) error` - Load with custom options
- `MustLoad(filenames ...string)` - Load files or panic

**Getting Variables:**
- `Get(key string) string` - Get environment variable value (alias for os.Getenv)
- `GetOrDefault(key, defaultValue string) string` - Get variable or return default if empty
- `GetOrPanic(key string) string` - Get variable or panic if not set/empty

#### Types

```go
type Options struct {
    Override bool  // Override existing environment variables
    Required bool  // File must exist (return error if not found)
}
```

#### Errors

- `ErrFileNotFound` - File does not exist
- `ErrInvalidFormat` - Invalid line format (missing =)
- `ErrEmptyKey` - Empty key name
- `ErrPermissionDenied` - No permission to read file

---

## Español

### Características

- **Minimalista**: Enfocado en cargar variables de entorno de manera segura - nada más, nada menos
- **Seguro**: Límites de tamaño de archivo, sanitización de rutas y verificación de permisos
- **Robusto**: Manejo completo de errores con tipos de error específicos
- **Limpio**: Estructura de código profesional siguiendo principios de arquitectura limpia
- **Flexible**: Soporte para múltiples archivos y opciones de configuración
- **Confiable**: Validación de entrada y análisis seguro

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
    // Cargar archivo .env por defecto
    if err := dotenv.Load(); err != nil {
        log.Printf("Advertencia: %v", err)
    }
    
    // Usar tus variables de entorno
    apiKey := os.Getenv("API_KEY")
    dbHost := os.Getenv("DB_HOST")
}
```

### Uso Avanzado

#### Cargar Múltiples Archivos

```go
// Cargar múltiples archivos (el primer archivo tiene prioridad)
err := dotenv.Load(".env.local", ".env")
```

#### Opciones Personalizadas

```go
opts := &dotenv.Options{
    Override: false,  // No sobrescribir variables existentes
    Required: true,   // El archivo debe existir
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
# Los comentarios son soportados
CLAVE_SIMPLE=valor

# Valores con comillas
VALOR_CON_COMILLAS="valor con espacios"
COMILLAS_SIMPLES='las comillas simples también funcionan'

# Valores vacíos
VALOR_VACIO=

# Caracteres especiales (escapados)
MULTILINEA="Línea 1\nLínea 2"
CON_TAB="Columna1\tColumna2"

# Eliminar espacios
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
    default:
        // Otro error
    }
}
```

### Características de Seguridad

- **Límite de tamaño**: Máximo 1MB para prevenir agotamiento de memoria
- **Sanitización de rutas**: Previene ataques de traversal de directorios
- **Verificación de permisos**: Valida permisos antes de leer
- **Validación de claves**: Solo permite nombres válidos de variables
- **Análisis seguro**: Maneja entradas malformadas de manera elegante

### Referencia API

#### Funciones

**Funciones de Carga:**
- `Load(filenames ...string) error` - Cargar uno o más archivos .env
- `LoadWithOptions(opts *Options, filenames ...string) error` - Cargar con opciones personalizadas
- `MustLoad(filenames ...string)` - Cargar archivos o hacer panic

**Funciones para Obtener Variables:**
- `Get(key string) string` - Obtener valor de variable de entorno (alias de os.Getenv)
- `GetOrDefault(key, defaultValue string) string` - Obtener variable o retornar valor por defecto si está vacía
- `GetOrPanic(key string) string` - Obtener variable o hacer panic si no está definida/vacía

#### Tipos

```go
type Options struct {
    Override bool  // Sobrescribir variables de entorno existentes
    Required bool  // El archivo debe existir (retorna error si no)
}
```

#### Errores

- `ErrFileNotFound` - El archivo no existe
- `ErrInvalidFormat` - Formato de línea inválido (falta =)
- `ErrEmptyKey` - Nombre de clave vacío
- `ErrPermissionDenied` - Sin permisos para leer el archivo

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

[@jaavier](https://github.com/jaavier)