# Go Web Starter

## Overview
The **Go Web Starter** template is a well-structured, production-ready Go web architecture designed to quickly kickstart web projects. It follows best practices by organizing code into layers (MVC-inspired) with clear separation of concerns.

---

## Project Structure
```text
go-web-starter/
├── cmd/                          # Application entry points
│   ├── main.go                   # Main function
│   └── run.go                    # Server logic
├── internal/                     # Private code (not importable externally)
│   ├── business/                 # Business logic
│   │   ├── controller/           # Controllers (HTTP request handlers)
│   │   │   └── starter.go        # Controller for the home page
│   │   ├── dto/                  # Data Transfer Objects (DTOs)
│   │   ├── model/                # Data models
│   │   │   └── starter.go        # Data structures for the home page
│   │   ├── repository/           # Data access layer (if needed)
│   │   ├── security/             # Security-related logic (if needed)
│   │   └── service/              # Business services (if needed)
│   ├── config/                   # Application configuration
│   │   ├── bcrypt.go             # Password hashing utility
│   │   ├── config.go             # Configuration struct
│   │   ├── createTemplate.go     # Template cache management
│   │   ├── config.go             # Configuration struct
│   │   └── jwt.go                # JWT token handling (if needed)
│   ├── database/                 # Database connection and management (if needed)
│   └── pkg/                      # Utilities and helpers
│       └── error.go              # HTTP error handling
├── templates/                    # HTML templates
│   ├── starter.page.tmpl         # Home page
│   ├── error.page.tmpl           # Error page
│   └── layouts/                  # Reusable layouts
│       └── base.layout.tmpl      # Base layout (HTML structure)
├── assets/                       # Static resources
│   ├── css/
│   │   └── app.css               # Custom CSS styles
│   ├── js/
│   │   └── app.js                # Custom JavaScript (empty by default)
│   └── img/
│       └── 1.ico                 # Favicon
├── go.mod                        # Project dependencies
├── Dockerfile                    # Docker configuration
├── .gitignore                    # Git ignore file
├── LICENSE                       # Project license
└── README.md                     # Basic project overview
```
---

## Entry Points
### `cmd/main.go`
Main entry point of the application. It accepts command-line arguments and passes them to the `run()` function.
```go
func main() {
    run(os.Args[1:])
}
```

### `cmd/run.go`
Contains the main server logic:
- **Configuration initialization**: template cache creation
- **Port configuration**: defaults to `8090`
- **Assets folder**: defaults to `assets`
- **Static routing**: serves CSS, JS, images via `/statics/`
- **Dynamic routing**: routes `/` to the `Starter` controller
- **Server startup**: listens on `http://127.0.0.1:5000`

---

## Simplified MVC Structure
### Controllers (`internal/business/controller/`)

Controllers handle HTTP requests and orchestrate business logic.

**`starter.go`**:
- Accepts GET requests on `/`
- Validates that the path is exactly `/`
- Prepares data (example: `owner: "Stackriv"`)
- Renders the `starter.page.tmpl` template
- Handles 404 errors

```go
func Starter(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        if r.URL.Path != "/" {
            err := pkg.Error(http.StatusNotFound)
            config.RenderTemplate(w, "error", model.Starter{...})
            return
        }
        names := make(map[string]string)
        names["owner"] = "Stackriv"
        config.RenderTemplate(w, "starter", model.Starter{StringData: names})
    }
}
```

### Models (`internal/business/model/`)

Models define data structures.
**`starter.go`**:
```go
type Starter struct {
    StringData map[string]string     // Text data
    IntData    map[string]int        // Numeric data
    Error      ErrorData             // Error information
}

type ErrorData struct {
    Code    string                   // HTTP code (e.g., "404")
    Message string                   // Error message (e.g., "Not Found")
}
```

---

## Configuration (`internal/config/`)
### `config.go`

Defines the application configuration structure:

```go
type Config struct {
    TemplateCache map[string]*template.Template   // Compiled templates cache
    Port          string                          // Server port
    StaticDir     string                          // Static files folder
}
```

### `createTemplate.go`

Manages template cache for optimal performance:

**`CreateTemplateCache()`**:
- Loads all `*.page.tmpl` files from `templates/`
- For each page, loads `*.layout.tmpl` files from `templates/layouts/`
- Returns a map with compiled templates

**`RenderTemplate()`**:
- Retrieves template from cache
- Executes template with data
- Writes result to HTTP response
- Handles execution errors

```go
func RenderTemplate(w http.ResponseWriter, tmplName string, tmplData interface{}) {
    tmpl, ok := appConfig.TemplateCache[tmplName+".page.tmpl"]
    if !ok {
        // Error 500
    }
    buffer := new(bytes.Buffer)
    err := tmpl.Execute(buffer, tmplData)
    // Error handling...
    buffer.WriteTo(w)
}
```

---
## Templates (`templates/`)

Templates use Go's `text/template` syntax.

### `base.layout.tmpl`

Base HTML structure reused by all pages:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/dark.css">
    <link rel="stylesheet" href="/statics/css/app.css">
    <link rel="shortcut icon" href="/statics/img/1.ico" type="image/x-icon">
    <title>Go Web Starter</title>
</head>
<body>
    {{ block "content" . }}{{ end }}
</body>
</html>
```

**Key elements**:
- **Water.css**: Minimalist CSS framework (dark theme)
- **app.css**: Custom styles
- **favicon**: Site icon
- **block "content"**: Area replaced with each page's content

### `starter.page.tmpl`

Home page of the site:

```go
{{ template "base" . }}

{{ define "content" }}
    <h1>Home Page</h1>
    <p>Hello everyone! I am {{ index .StringData "owner" }}</p>
{{ end }}
```

**How it works**:
- Uses the `base` layout
- Defines the `content` block
- Accesses data via `.StringData` and `index`

### `error.page.tmpl`

Generic error page:

```go
{{ template "base" . }}

{{ define "content" }}
    <h1>Error Page</h1>
    <h2>{{ index .Error.Code }}</h2>
    <p>{{ index .Error.Message }}</p>
{{ end }}
```

---

## Utilities (`internal/pkg/`)

### `error.go`

Helper function to map HTTP codes to readable messages:

```go
func Error(code int) map[string]string {
    switch code {
    case http.StatusBadRequest:
        msg = "Bad Request"
    case http.StatusNotFound:
        msg = "Not Found"
    // ... other codes ...
    }
    return map[string]string{
        "code": strconv.Itoa(code),
        "msg":  msg,
    }
}
```

**Supported codes**:
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `405` - Method Not Allowed
- `409` - Conflict
- `500` - Internal Server Error

---

## Static Assets (`assets/`)

### CSS (`assets/css/app.css`)

Minimalist custom styles:

```css
* {
    margin: 0;
    padding: 0;
    user-select: none;
}
```

**Usage**:
- Resets default margins/paddings
- Disables text selection by mouse

### JavaScript (`assets/js/app.js`)

Empty by default, ready for custom code.

### Images (`assets/img/`)

- `1.ico`: Site favicon

---

## Docker (`Dockerfile`)

Minimal Docker configuration:

```dockerfile
FROM golang:lts-alpine
```

**Enhancement points**:
- Uses Alpine for a lightweight image
- Can be extended to copy code, download dependencies, and compile

---

## Dependencies (`go.mod`)

```
module github.com/stackriv/go-web-starter
go 1.25.4
```

**Status**:
- No external dependencies
- Uses only Go stdlib
- Requires Go 1.25.4 or higher

---

## Quick Start

### 1. Preparation
```bash
cd /path/to/go-web-starter
```

### 2. Launch the server
```bash
go run ./cmd/main.go
```

### 3. Access the application
Open browser at: **http://127.0.0.1:5000**

---

## Extension Guide

### Add a new route

1. **Create a controller** in `internal/business/controller/`:
```go
func MyPage(w http.ResponseWriter, r *http.Request) {
    data := model.Starter{StringData: map[string]string{"title": "My page"}}
    config.RenderTemplate(w, "mypage", data)
}
```

2. **Add the route** in `cmd/run.go`:
```go
http.HandleFunc("/mypage", controller.MyPage)
```

3. **Create the template** `templates/mypage.page.tmpl`:
```go
{{ template "base" . }}
{{ define "content" }}
    <h1>{{ index .StringData "title" }}</h1>
{{ end }}
```

### Add a new route with query parameters

```go
func Article(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    data := model.Starter{StringData: map[string]string{"id": id}}
    config.RenderTemplate(w, "article", data)
}
```

### Add CSS styles

Modify `assets/css/app.css`:
```css
body {
    font-family: Arial, sans-serif;
    max-width: 1200px;
    margin: 0 auto;
}
```

---

## Built-in Best Practices

✔️ **Separation of Concerns**: Controllers, Models, Config clearly separated  
✔️ **Template Caching**: Templates compiled once at startup  
✔️ **Error Handling**: HTTP errors mapped to readable messages  
✔️ **No External Dependencies**: Uses only Go stdlib  
✔️ **Separated Static Assets**: CSS, JS, images in dedicated folder  
✔️ **Docker Ready**: Dockerfile provided for containerization  
✔️ **Git Friendly**: Complete `.gitignore` for Go, macOS, Windows, Linux, IDEs

---

## Useful Go Resources

- [Go Standard Library](https://pkg.go.dev/std)
- [Text/Template Documentation](https://pkg.go.dev/text/template)
- [Net/HTTP Documentation](https://pkg.go.dev/net/http)

---

**Version**: 1.0  
**Module**: github.com/stackriv/go-web-starter  
**Go**: 1.25.4+  
**License**: See LICENSE