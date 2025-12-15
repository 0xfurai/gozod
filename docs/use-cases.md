# Use Cases

Real-world use case examples for Go Validation library.

## Table of Contents

- [HTTP API Request Validation](#http-api-request-validation)
- [Configuration Validation](#configuration-validation)

## HTTP API Request Validation

There are two approaches to validate HTTP request payloads: using maps or decoding directly into structs.

### Approach 1: Decode Directly into Structs (Recommended)

For type-safe validation, decode JSON payloads directly into Go structs and validate them:

```go
type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"isActive"`
}

var userSchema = gozod.Struct(gozod.Shape{
    "name": gozod.String().
        Min(3).
        Max(50).
        CustomError(gozod.ErrCodeTooSmall, "Name must be at least 3 characters").
        CustomError(gozod.ErrCodeRequired, "Name is required"),
    "email": gozod.String().
        Email().
        CustomError(gozod.ErrCodeInvalidString, "Invalid email format"),
    "age": gozod.Int().
        Min(18).
        Max(120).
        CustomError(gozod.ErrCodeTooSmall, "You must be at least 18 years old"),
    "isActive": gozod.Bool().Nilable(),
})

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse JSON body directly to struct
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate the struct
    errors := userSchema.Validate(user, nil)
    if errors != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors.FormatErrorsJSON())
        return
    }

    // Use validated user struct with type safety
    fmt.Printf("User: %+v\n", user)
    // ... save to database, etc.

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

**Benefits:**
- Type-safe access to validated data (no type assertions needed)
- Better IDE autocomplete and compile-time checking
- More idiomatic Go code
- Reduced runtime errors from type assertions

### Approach 2: Using Maps

For dynamic payloads where structure isn't known at compile time:

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var payload map[string]any
    json.NewDecoder(r.Body).Decode(&payload)

    errors := userSchema.Validate(payload, nil)

    if errors != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors.FormatErrorsJSON())
        return
    }

    // Use validated payload (requires type assertions)
    // ... save to database, etc.
}
```

### Complete Example with Struct Decoding

This example shows a complete HTTP server with struct-based validation:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "gozod"
)

type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"isActive"`
}

// Define validation schema using gozod.Struct
var userSchema = gozod.Struct(gozod.Shape{
    "name": gozod.String().
        Min(3).
        Max(50).
        CustomError(gozod.ErrCodeTooSmall, "Name must be at least 3 characters").
        CustomError(gozod.ErrCodeRequired, "Name is required"),
    "email": gozod.String().
        Email().
        CustomError(gozod.ErrCodeInvalidString, "Invalid email format"),
    "age": gozod.Int().
        Min(18).
        Max(120).
        CustomError(gozod.ErrCodeTooSmall, "You must be at least 18 years old"),
    "isActive": gozod.Bool().Nilable(),
})

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Step 1: Decode JSON body directly to struct
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Step 2: Validate the struct
    errors := userSchema.Validate(user, nil)
    if errors != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors.FormatErrorsJSON())
        return
    }

    // Step 3: Use validated user struct with type safety
    fmt.Printf("User: %+v\n", user)

    // Save to database (type-safe access)
    // db.Create(&user)
    // No need for type assertions like user["name"].(string)

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func main() {
    http.HandleFunc("/users", CreateUserHandler)
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: curl -X POST http://localhost:8080/users \\")
    fmt.Println("  -H 'Content-Type: application/json' \\")
    fmt.Println("  -d '{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"age\":25,\"isActive\":true}'")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Example request:**
```bash
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 25,
    "isActive": true
  }'
```

**Example validation error response:**
```json
{
  "age": ["You must be at least 18 years old"],
  "email": ["Invalid email format"],
  "name": ["Name must be at least 3 characters"]
}
```

### Using Flattened Errors for Forms

For frontend form validation, use flattened errors:

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    errors := userSchema.Validate(user, nil)
    if errors != nil {
        // Flatten errors for easier frontend consumption
        flattened := errors.Flatten()

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(flattened)
        return
    }

    // Continue with processing...
}
```

**Flattened error format:**
```json
{
  "name": "Name must be at least 3 characters",
  "email": "Invalid email format",
  "age": "You must be at least 18 years old"
}
```

### Handling Nested Structs

For complex payloads with nested structures:

```go
type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    ZipCode string `json:"zipCode"`
}

type UserWithAddress struct {
    Name    string  `json:"name"`
    Email   string  `json:"email"`
    Age     int     `json:"age"`
    Address Address `json:"address"`
}

var addressSchema = gozod.Struct(gozod.Shape{
    "street":  gozod.String().Min(1),
    "city":    gozod.String().Min(1),
    "zipCode": gozod.String().Regex(`^\d{5}$`, "Must be 5 digits"),
})

var userWithAddressSchema = gozod.Struct(gozod.Shape{
    "name":    gozod.String().Min(3).Max(50),
    "email":   gozod.String().Email(),
    "age":     gozod.Int().Min(18),
    "address": addressSchema,
})

func CreateUserWithAddressHandler(w http.ResponseWriter, r *http.Request) {
    var user UserWithAddress
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate nested structure
    errors := userWithAddressSchema.Validate(user, nil)
    if errors != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors.FormatErrorsJSON())
        return
    }

    // Access validated nested fields with type safety
    fmt.Printf("User %s lives in %s\n", user.Name, user.Address.City)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### Partial Updates (PATCH Requests)

For partial updates, use pointers to distinguish between "not provided" and "zero value":

```go
type UserUpdate struct {
    Name     *string `json:"name,omitempty"`
    Email    *string `json:"email,omitempty"`
    Age      *int    `json:"age,omitempty"`
    IsActive *bool   `json:"isActive,omitempty"`
}

var userUpdateSchema = gozod.Struct(gozod.Shape{
    "name":     gozod.String().Min(3).Max(50).Nilable(),
    "email":    gozod.String().Email().Nilable(),
    "age":      gozod.Int().Min(18).Max(120).Nilable(),
    "isActive": gozod.Bool().Nilable(),
})

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    var update UserUpdate
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate only provided fields
    errors := userUpdateSchema.Validate(update, nil)
    if errors != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors.FormatErrorsJSON())
        return
    }

    // Update only provided fields
    if update.Name != nil {
        fmt.Printf("Updating name to: %s\n", *update.Name)
        // db.UpdateName(*update.Name)
    }
    if update.Email != nil {
        fmt.Printf("Updating email to: %s\n", *update.Email)
        // db.UpdateEmail(*update.Email)
    }
    // ... other fields

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
```

## Configuration Validation

Validate application configuration on startup:

```go
type Config struct {
    Port     int      `json:"port"`
    Host     string   `json:"host"`
    Timeout  float64  `json:"timeout"`
    Features []string `json:"features"`
}

var configSchema = gozod.Map(map[string]gozod.Schema{
    "port":     gozod.Int().Min(1).Max(65535),
    "host":     gozod.String().Min(1),
    "timeout":  gozod.Float().Positive(),
    "features": gozod.Array(gozod.String()).Nilable(),
})

func LoadConfig(configPath string) (*Config, error) {
    // Load config from file
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var configData map[string]any
    if err := json.Unmarshal(data, &configData); err != nil {
        return nil, err
    }

    // Validate config
    errors := configSchema.Validate(configData, nil)
    if errors != nil {
        return nil, fmt.Errorf("invalid configuration: %s", errors.FormatErrors())
    }

    // After validation, configData is safe to use
    config := &Config{
        Port:     configData["port"].(int),
        Host:     configData["host"].(string),
        Timeout:  configData["timeout"].(float64),
        Features: convertToStringSlice(configData["features"]),
    }

    return config, nil
}

func convertToStringSlice(val any) []string {
    if val == nil {
        return nil
    }
    if slice, ok := val.([]any); ok {
        result := make([]string, len(slice))
        for i, v := range slice {
            result[i] = v.(string)
        }
        return result
    }
    return nil
}

func main() {
    config, err := LoadConfig("config.json")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Use validated config
    fmt.Printf("Starting server on %s:%d\n", config.Host, config.Port)
}
```

### Environment Variable Validation

```go
func LoadConfigFromEnv() (*Config, error) {
    port, _ := strconv.Atoi(os.Getenv("PORT"))
    host := os.Getenv("HOST")
    timeout, _ := strconv.ParseFloat(os.Getenv("TIMEOUT"), 64)

    configData := map[string]any{
        "port":    port,
        "host":    host,
        "timeout": timeout,
    }

    errors := configSchema.Validate(configData, nil)
    if errors != nil {
        return nil, fmt.Errorf("invalid configuration: %s", errors.FormatErrors())
    }

    config := &Config{
        Port:    port,
        Host:    host,
        Timeout: timeout,
    }
    return config, nil
}
```

## See Also

- [Examples](examples.md) - More validation examples
- [API Reference](api-reference.md) - Complete API documentation
- [Error Handling](error-handling.md) - Error handling guide
