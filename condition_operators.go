package rdb

const (
	opEqual        string = "="
	opNotEqual     string = "!="
	opGreater      string = ">"
	opGreaterEqual string = ">="
	opLess         string = "<"
	opLessEqual    string = "<="
	opIn           string = "IN"
	opNotIn        string = "NOT IN"
	opAnd          string = "AND"
	opOr           string = "OR"
)

func Equal[T any](key *T, value T) *kvCondition {
	return &kvCondition{
		pair:     keyValue(key, value),
		operator: opEqual,
	}
}

func NotEqual[T any](key *T, value T) *kvCondition {
	return &kvCondition{
		pair:     keyValue(key, value),
		operator: opNotEqual,
	}
}

func Greater[T any](key *T, value T) *kvCondition {
	return &kvCondition{
		pair:     keyValue(key, value),
		operator: opGreater,
	}
}

func GreaterEqual[T any](key *T, value T) *kvCondition {
	return &kvCondition{
		pair:     keyValue(key, value),
		operator: opGreaterEqual,
	}
}

func Less[T any](key *T, value T) *kvCondition {
	return &kvCondition{
		pair:     keyValue(key, value),
		operator: opLess,
	}
}

func LessEqual[T any](key *T, value T) *kvCondition {
	return &kvCondition{
		pair:     keyValue(key, value),
		operator: opLessEqual,
	}
}

func In[T any](key *T, values []T) *klCondition {
	return &klCondition{
		pair:         keyList(key, values),
		listOperator: opIn,
		soloOperator: opEqual,
	}
}

func NotIn[T any](key *T, values []T) *klCondition {
	return &klCondition{
		pair:         keyList(key, values),
		listOperator: opNotIn,
		soloOperator: opNotEqual,
	}
}

func And(conditions ...Condition) *conditionSet {
	return &conditionSet{
		conditions: conditions,
		operator:   opAnd,
	}
}

func Or(conditions ...Condition) *conditionSet {
	return &conditionSet{
		conditions: conditions,
		operator:   opOr,
	}
}
