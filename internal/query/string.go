package query

import (
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/types"
)

func QueryString(q Query) string {
	query, rawValues := q.Build()
	values := make([]any, len(rawValues))
	for i, value := range rawValues {
		typ := fmt.Sprintf("%T", value)
		if strings.HasPrefix(typ, "*") {
			values[i] = types.Deref(value)
		} else {
			values[i] = value
		}
	}
	query = strings.Replace(query, "?", "%v", -1)
	return fmt.Sprintf(query, values...)
}
