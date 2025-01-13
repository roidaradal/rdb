package rdb

type kvc struct {
	field  any
	value  any
	column string
}

type klc struct {
	field  any
	values []any
	column string
}

/*
Input: &struct.Field, value

Constraint: Type of key, value must match (compile-time validation)

Output: pointer to new kvc struct (no column yet)

Example: a = Account{}; KV(&a.Name, "john")
*/
func keyValue[T any](key *T, value T) *kvc {
	return &kvc{
		field: key,
		value: value,
	}
}

/*
Input: &struct.Field, values

Constraint: Type of key, values must match (compile-time validation)

Output: pointer to new klc struct (no column yet)

Example: a = Account{}; KV(&a.ID, []uint{5, 6, 7})
*/
func keyList[T any](key *T, values []T) *klc {
	values2 := make([]any, len(values))
	for i, value := range values {
		values2[i] = value
	}
	return &klc{
		field:  key,
		values: values2,
	}
}

/******************************** BUILD METHODS ********************************/

/*
Input: &struct

Constraint: Same struct used for setting Field=&struct.Field

Output: column (string), value (any)

Note: column could be blank string
*/
func (p *kvc) Build(t any) (string, any) {
	p.column = Column(t, p.field)
	return p.column, p.value
}

/*
Input: &struct

Constraint: Same struct used for setting Field=&struct.Field

Output: column (string), values ([]any)

Note: column could be blank string
*/
func (p *klc) Build(t any) (string, []any) {
	p.column = Column(t, p.field)
	return p.column, p.values
}
