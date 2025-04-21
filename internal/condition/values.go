package condition

import (
	"fmt"

	"github.com/roidaradal/rdb/internal/op"
	"github.com/roidaradal/rdb/internal/types"
)

const (
	defaultCondition string = "false"
	trueCondition    string = "true"
)

func defaultConditionValues() (string, []any) {
	return defaultCondition, []any{}
}

func matchAllConditionValues() (string, []any) {
	return trueCondition, []any{}
}

func singleConditionValues(column, operator string, value any) (string, []any) {
	isValueNil := types.IsNil(value)
	if operator == op.Equal && isValueNil {
		return fmt.Sprintf("%s IS NULL", column), []any{}
	} else if operator == op.NotEqual && isValueNil {
		return fmt.Sprintf("%s IS NOT NULL", column), []any{}
	} else if operator == op.Prefix {
		prefix := fmt.Sprintf("%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{prefix}
	} else if operator == op.Suffix {
		suffix := fmt.Sprintf("%%%v", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{suffix}
	} else if operator == op.Substring {
		substring := fmt.Sprintf("%%%v%%", value)
		return fmt.Sprintf("%s LIKE ?", column), []any{substring}
	} else {
		return fmt.Sprintf("%s %s ?", column, operator), []any{value}
	}
}

func listCondition(column, operator string, numValues int) string {
	placeholders := op.RepeatString(numValues, "?", ", ")
	return fmt.Sprintf("%s %s (%s)", column, operator, placeholders)
}
