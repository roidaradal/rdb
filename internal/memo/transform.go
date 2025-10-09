package memo

import (
	"strings"

	"github.com/roidaradal/fn"
)

// Transforms a field value (any)
type TransformFn = func(any) any

var transformers = map[string]TransformFn{
	"upper":    upper,
	"lower":    lower,
	"upperdot": upperdot,
	"lowerdot": lowerdot,
}

// Uppercase transform
func upper(item any) any {
	// assumes item is a string
	text, ok := item.(string)
	if !ok {
		return item
	}
	text = strings.TrimSpace(text)
	return strings.ToUpper(text)
}

// Lowercase transform
func lower(item any) any {
	// assumes item is a string
	text, ok := item.(string)
	if !ok {
		return item
	}
	text = strings.TrimSpace(text)
	return strings.ToLower(text)
}

// If empty string, default to '.'
func guardDot(item any) any {
	// assumes item is a string
	text, ok := item.(string)
	if !ok {
		return item
	}
	text = strings.TrimSpace(text)
	return fn.Ternary(text == "", ".", text)
}

// Uppercase transform + default to '.' if empty string
func upperdot(item any) any {
	return upper(guardDot(item))
}

// Lowercase transform + default to '.' if empty string
func lowerdot(item any) any {
	return lower(guardDot(item))
}
