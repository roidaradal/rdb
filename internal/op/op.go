package op

import (
	"slices"
	"strings"
)

const (
	Equal        string = "="
	NotEqual     string = "!="
	Greater      string = ">"
	GreaterEqual string = ">="
	Less         string = "<"
	LessEqual    string = "<="
	In           string = "IN"
	NotIn        string = "NOT IN"
	And          string = "AND"
	Or           string = "OR"
	Prefix       string = "PREFIX"
	Suffix       string = "SUFFIX"
	Substring    string = "SUBSTRING"
)

func RepeatString(repeat int, item, glue string) string {
	return strings.Join(slices.Repeat([]string{item}, repeat), glue)
}
