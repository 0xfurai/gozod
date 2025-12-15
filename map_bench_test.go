package gozod

import (
	"testing"
)

type BenchPerson struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type BenchAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
}

type BenchUser struct {
	Name    string       `json:"name"`
	Email   string       `json:"email"`
	Age     int          `json:"age"`
	Address BenchAddress `json:"address"`
}

var personSchema = Map(map[string]Schema{
	"name":  String().Min(2).Max(50),
	"email": String().Email(),
	"age":   Int().Min(0).Max(150),
})

var userSchema = Map(map[string]Schema{
	"name":  String().Min(2).Max(50),
	"email": String().Email(),
	"age":   Int().Min(0).Max(150),
	"address": Map(map[string]Schema{
		"street":   String().Min(5),
		"city":     String().Min(2),
		"zip_code": String().Min(5),
	}),
})

func BenchmarkMapSchema_ValidateStruct_Simple(b *testing.B) {
	person := BenchPerson{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = personSchema.Validate(person, nil)
	}
}

func BenchmarkMapSchema_ValidateMap_Simple(b *testing.B) {
	data := map[string]any{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = personSchema.Validate(data, nil)
	}
}

func BenchmarkMapSchema_ValidateStruct_Nested(b *testing.B) {
	user := BenchUser{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: BenchAddress{
			Street:  "123 Main Street",
			City:    "New York",
			ZipCode: "10001",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userSchema.Validate(user, nil)
	}
}

func BenchmarkMapSchema_ValidateMap_Nested(b *testing.B) {
	data := map[string]any{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
		"address": map[string]any{
			"street":   "123 Main Street",
			"city":     "New York",
			"zip_code": "10001",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userSchema.Validate(data, nil)
	}
}
