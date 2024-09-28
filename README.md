Вот улучшенная версия README.md, соответствующая стандартам оформления и удобочитаемости:

---

# SLIB

> **For non-commercial use only**

This library is designed to streamline and accelerate development tasks by providing convenient utility functions.

---

## Getting Started

To install the library, run:

```bash
go get github.com/sabbatD/slib
```

### Overview

SLIB includes several utility packages to help with common tasks:

- **Math Package**: Provides a set of mathematical functions.
- **Stringutils Package**: Includes alphabetical letters, digits, and pointer utilities.

---

## Handle Package

### `ChangeConfig` Function

The `ChangeConfig` function allows you to customize the configuration of the library. It accepts the following parameters:

- **`rps`** (uint16) - Requests per second.
- **`duration`** (uint16) - Duration for request repetition.
- **`detain`** (bool) - Controls the request generation method.

**Example**:

```go
ChangeConfig(100, 5, true)
```

If you do not explicitly configure these settings, the library will use the default values:

- `rps`: 100
- `duration`: 5 seconds
- `detain`: false

---

## HTTP Request Functions

### `Get(url)` and `Post(url, body)`

These functions are used to make single HTTP requests:

- **`Get(url string)`**: Sends a GET request to the specified `url`.
- **`Post(url string, body interface{})`**: Sends a POST request to the specified `url` with the provided `body`.

**Example**:

```go
Get("https://example.com")
Post("https://example.com", body)
```

---

## Attack Function

The `Attack` function is designed to perform multiple HTTP requests at a given rate and duration. It supports both GET and POST methods.

### Syntax:

```go
Attack(method, url string, body ...[]byte) string
```

- **`method`**: HTTP method (`"GET"` or `"POST"`).
- **`url`**: The target URL.
- **`body`**: (Optional) Request body for POST requests.

The function returns a string that details how many requests were made and how long the execution took.

**Example**:

```go
fmt.Println(Attack("GET", "https://localhost"))
Attack("POST", "https://localhost", body)
```

**Output**:

```
1200 requests in 3.0003329 seconds
```

---

### Author

SLIB is developed and maintained by **S4BB4T**.

---