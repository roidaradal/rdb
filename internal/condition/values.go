package condition

import (
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
)

const (
	falseCondition string = "false"
	trueCondition  string = "true"
)

// Returns 'false' as condition, empty list for values
func falseConditionValues() (string, []any) {
	return falseCondition, []any{}
}

// Returns 'true' as condition, empty list for values
func trueConditionValues() (string, []any) {
	return trueCondition, []any{}
}

// Build condition string and query parameter values list,
// Values list corresponds to the ? in the query,
// Used for solo value conditions
func soloConditionValues(column, operator string, value any) (string, []any) {
	isValueNil := dyn.IsNull(value)
	if operator == Equal && isValueNil {
		return fmt.Sprintf("%s IS NULL", column), []any{}
	} else if operator == NotEqual && isValueNil {
		return fmt.Sprintf("%s IS NOT NULL", column), []any{}
	} else if operator == Prefix {
		prefix := fmt.Sprintf("%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{prefix}
	} else if operator == Suffix {
		suffix := fmt.Sprintf("%%%v", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{suffix}
	} else if operator == Substring {
		substring := fmt.Sprintf("%%%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{substring}
	} else {
		return fmt.Sprintf("%s %s ?", column, operator), []any{value}
	}
}

// Build condition string for list value conditions,
// Adds the repeated placeholder ? to end of the condition
func listCondition(column, operator string, numValues int) string {
	placeholders := str.Repeat(numValues, "?", ", ")
	return fmt.Sprintf("%s %s (%s)", column, operator, placeholders)
}
