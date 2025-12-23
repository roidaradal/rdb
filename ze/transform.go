package ze

import (
	"strings"

	"github.com/roidaradal/fn/str"
)

// Transforms a field value (any)
type TransformFn = func(any) any

var transformers = map[string]TransformFn{
	"upper":    upper,
	"lower":    lower,
	"upperdot": upperdot,
	"lowerdot": lowerdot,
}

// TransformFn: trimSpace + uppercase
func upper(item any) any {
	// Assumes item is string
	text, ok := item.(string)
	if !ok {
		return item
	}
	return strings.ToUpper(strings.TrimSpace(text))
}

// TransformFn: trimSpace + lowercase
func lower(item any) any {
	// Assumes item is string
	text, ok := item.(string)
	if !ok {
		return item
	}
	return strings.ToLower(strings.TrimSpace(text))
}

// TransformFn: trimSpace + uppercase + guardDot
func upperdot(item any) any {
	return upper(guardDot(item))
}

// TransformFn: trimSpace + lowercase + guardDot
func lowerdot(item any) any {
	return lower(guardDot(item))
}

// Common: guard blank string with dot
func guardDot(item any) any {
	// Assumes item is string
	text, ok := item.(string)
	if !ok {
		return item
	}
	return str.GuardDot(strings.TrimSpace(text))
}
