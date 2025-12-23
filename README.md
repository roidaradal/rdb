# RDB
A Go library for building type-safe SQL queries.

The `ze` package contains Schema and Request types.

`go get github.com/roidaradal/rdb/...`

## Initialization 

### Initialize
Required call to initialize rdb library

`err := rdb.Initialize()`

### Type Definition
By default, field name is also the column name. 
Fields can be mapped to a custom column name using the `col:"CustomName"` struct tag.
Fields can also be skipped when processing columns using the `col:"-"` struct tag.

```
type Foo struct {
    Field1      string                      // column name: Field1
    Field2      string  `col:"CustomName"`  // column name is CustomName, instead of Field2
    PrivField_  int     `col:"-"`           // skipped
}
```

### AddType 
Registers a new type to RDB. Expects a struct pointer parameter.

```
structRef := &Type{}
err := rdb.AddType(structRef)
```


## DB Connection

### _type:_ SQLConnParams 
Parameters for SQL connection: Host, Port, Username, Password, Database name

### NewSQLConnection
Creates a new MySQL DB connection pool, using the SQLConnParams

```
var dbc *sql.DB 
p := &rdb.SQLConnParams{...}
dbc, err := rdb.NewSQLConnection(p)
```

## Columns and Rows 

### AllColumns
Get all column names of given item's type

```
columnNames := rdb.AllColumns(Type{})
columnNames := rdb.AllColumns(&Type{})
```

### Column
Get column name of given field pointer,
Field must be from the struct used in AddType

`columnName := rdb.Column(&item.Field)`

### Columns
Get column names of given field pointers, 
Fields must be from the struct used in AddType.
If any columns are not found, returns an empty list.

`columnNames := rdb.Columns(&item.Field1, &item.Field2, &item.Field3)`

### Field 
Get field name of given field pointer, 
Field must be from the struct used in AddType

`fieldName := rdb.Field(typeName, &item.Field)`

### Fields 
Get field names of given field pointers, 
Fields must be from the struct used in AddType.
If any fields are not found, returns an empty list.

`fieldNames := rdb.Fields(typeName, &item.Field1, &item.Field2)`

### _type:_ RowReader[T]
Function that reads row values into struct

### NewReader[T]
Creates a RowReader[T] with the given columns 

`reader := rdb.NewReader[T](column1, column2, ...)`

### FullReader[T]
Creates a RowReader[T] using all columns of type T

`reader := rdb.FullReader(&T{})`

### ToRow[T]
Converts given struct to map[string]any for row insertion 

`row := rdb.ToRow(&item)`

## Conditions

### _interface:_ Condition
Unifies the different Condition types into one interface

### NoCondition 
Match all condition 

`condition := rdb.NoCondition()`

### Equal
`condition := rdb.Equal(&item.Field, value)`

### NotEqual
`condition := rdb.NotEqual(&item.Field, value)`

### Prefix
`condition := rdb.Prefix(&item.Field, prefix)`

### Suffix
`condition := rdb.Suffix(&item.Field, suffix)`

### Substring
`condition := rdb.Substring(&item.Field, substring)`

### Greater
`condition := rdb.Greater(&item.Field, value)`

### GreaterEqual
`condition := rdb.GreaterEqual(&item.Field, value)`

### Less
`condition := rdb.Less(&item.Field, value)`

### LessEqual
`condition := rdb.LessEqual(&item.Field, value)`

### In
```
values := []T{...}
condition := rdb.In(&item.Field, values)
```

### NotIn
```
values := []T{...}
condition := rdb.NotIn(&item.Field, values)
```

### And
`condition := rdb.And(condition1, condition2, ...)`

### Or
`condition := rdb.Or(condition1, condition2, ...)`

## Queries 

### _interface:_ Query
Unifies the different Query types into one interface

### QueryString 
Builds the query object and outputs the query string 

`queryString := rdb.QueryString(q)`

### NewCountQuery 
Creates a new CountQuery, can also be used for ExistsQuery

```
q := rdb.NewCountQuery(table)
q.Where(condition)
count, err := q.Count(*sql.DB)   // int
exists, err := q.Exists(*sql.DB) // boolean
```

### NewDeleteQuery 
Creates a new DeleteQuery 

```
q := rdb.NewDeleteQuery(table)
q.Where(condition)
```

### NewDistinctValuesQuery 
Creates a new DistinctValuesQuery 

```
q := rdb.NewDistinctValuesQuery(table, &item.Field)
q.Where(condition) // optional
uniqueValues, err := q.Query(*sql.DB)
```

### NewInsertRowQuery 
Creates a new InsertRowQuery

```
q := rdb.NewInsertRowQuery(table)
q.Row(rdb.ToRow(&item))
```

### NewInsertRowsQuery
Creates a new InsertRowsQuery 

```
rows := []map[string]any{...}
q := rdb.NewInsertRowsQuery(table)
q.Rows(rows)
```

### NewLookupQuery 
Creates a new LookupQuery 

```
var lookup map[K]V
q := rdb.NewLookupQuery(table, &item.KeyField, &item.ValueField)
q.Where(condition) // optional
lookup, err := q.Lookup(*sql.DB)
```

### NewSelectRowQuery 
Creates a new SelectRowQuery with selected columns (set later)

```
columns := rdb.Columns(&item.Field1, &item.Field2, ...)
reader := rdb.Reader(columns...)
q := rdb.NewSelectRowQuery(table, reader)
q.Columns(columns)
q.Where(condition)
item, err := q.QueryRow(*sql.DB)
```

### NewFullSelectRowQuery 
Creates a new SelectRowQuery that uses all columns

```
q := rdb.NewFullSelectRowQuery(table, reader)
q.Where(condition)
item, err := q.QueryRow(*sql.DB)
```

### NewSelectRowsQuery 
Creates a new SelectRowsQuery with selected columns (set later)

```
columns := rdb.Columns(&item.Field1, &item.Field2, ...)
reader := rdb.Reader(columns...)
q := rdb.NewSelectRowsQuery(table, reader)
q.Columns(columns)
q.Where(condition)                   // optional
q.Limit(limit)                       // optional
q.Page(number, batchSize)            // optional
q.OrderAsc(rdb.Column(&item.Field))  // optional
q.OrderDesc(rdb.Column(&item.Field)) // optional
items, err := q.Query(*sql.DB)
```

### NewFullSelectRowsQuery
Creates a new SelectRowsQuery that uses all columns

```
q := rdb.NewFullSelectRowsQuery(table, reader)
q.Where(condition)                   // optional
q.Limit(limit)                       // optional
q.Page(number, batchSize)            // optional
q.OrderAsc(rdb.Column(&item.Field))  // optional
q.OrderDesc(rdb.Column(&item.Field)) // optional
items, err := q.Query(*sql.DB)
```

### NewTopRowQuery 
Creates a new TopRowQuery.
For top 1 row, use QueryRow().
For top N rows, set N using Limit() and use QueryRows().

```
q := rdb.NewTopRowQuery(table, reader)
q.Where(condition)
q.OrderAsc(rdb.Column(&item.Field)) // or 
q.OrderDesc(rdb.Column(&item.Field))
topItem, err := q.QueryRow(*sql.DB)
```

```
q := rdb.NewTopRowQuery(table, reader)
q.Where(condition)
q.Limit(5)
q.OrderAsc(rdb.Column(&item.Field)) // or 
q.OrderDesc(rdb.Column(&item.Field))
topItems, err := q.QueryRows(*sql.DB)

```

### NewTopValueQuery
Creates a new TopValueQuery.
For top 1 value, use QueryValue().
For top N values, set N using Limit() and use QueryValues().

```
var topValue V
q := rdb.NewTopValueQuery(table, &item.Field)
q.Where(condition)
q.OrderAsc(rdb.Column(&item.Field)) // or 
q.OrderDesc(rdb.Column(&item.Field))
topValue, err := q.QueryValue(*sql.DB)
```

```
q := rdb.NewTopValueQuery(table, &item.Field)
q.Where(condition)
q.Limit(5)
q.OrderAsc(rdb.Column(&item.Field)) // or 
q.OrderDesc(rdb.Column(&item.Field))
topValues, err := q.QueryValues(*sql.DB)
```

### NewUpdateQuery, Update
Creates a new UpdateQuery and adds field updates

```
q := rdb.NewUpdateQuery(table)
q.Where(Condition)
rdb.Update(q, &item.Field1, value1)
rdb.Update(q, &item.Field2, value2) // or 
q.Update(fieldName, value)
q.Updates(map[fieldName]value)      // values = any type
```

### NewValueQuery 
Creates a new ValueQuery

```
q := rdb.NewValueQuery(table, &item.Field)
q.Where(condition)
value, err := q.QueryValue(*sql.DB)
```

## Execution and Results

### _type:_ ResultChecker 
Function that checks SQL result if a condition has been satisfied

### AssertNothing 
ResultChecker that does nothing 

`checker := rdb.AssertNothing`

### AssertRowsAffected
Creates a ResultChecker that asserts number of rows affected

`checker := rdb.AssertRowsAffected(1)`
s
### RowsAffected
Gets the number of rows affected from SQL result (default: 0)

`affected := rdb.RowsAffected(*sql.Result)`

### LastInsertID 
Gets the last insert ID (uint) from SQL result (default: 0)

`id, ok := rdb.LastInsertID(*sql.Result)`

### Exec
Executes an SQL query

`result, err := rdb.Exec(q, *sql.DB)`

### ExecTx
Executes an SQL query as part of a transaction, applies Rollback on any errors

`result, err := rdb.ExecTx(q, *sql.Tx, checker)`

### Rollback 
Rolls back the SQL transaction

`err := rdb.Rollback(*sql.Tx, err)`

## Ze 

### Initialize 

`err := ze.Initialize(*rdb.SQLConnParams)`

### Errors and Status Codes 

* _error_: ze.ErrInactiveItem
* _error_: ze.ErrInvalidField
* _error_: ze.ErrMissingField
* _error_: ze.ErrMissingParams
* _error_: ze.ErrNotFoundItem
* _error_: ze.ErrMissingSchema
* _status_: ze.OK200 (OK)
* _status_: ze.OK201 (Created)
* _status_: ze.Err400 (Missing client parameters)
* _status_: ze.Err401 (Unauthenticated)
* _status_: ze.Err403 (Unauthorized)
* _status_: ze.Err404 (Not Found)
* _status_: ze.Err429 (Rate limited)
* _status_: ze.Err500 (Server-side Error)

### Types 

* ID
* Date 
* DateTime 
* UniqueItem    : ID 
* CodedItem     : Code 
* CreatedItem   : CreatedAt 
* ActiveItem    : IsActive
* Identity      : ID, Code   
* Item          : ID, Code, IsActive, CreatedAt  


### Items, ItemsRef 
Items schema and get Item reference object

```
ze.Items // *Schema[Item]
itemsRef := ze.ItemsRef()
```

### _type_: Task 
Contains Action and Target of Task 

### _type_: Request 
Application request that holds DB connection, transaction, checker, 
transaction queries, request start time, and logs 

```
var rq *Request 
rq, err := NewRequest(name string, args ...any)
rq.AddLog(message)
rq.AddFmtLog(format, args ...any)
rq.AddDurationLog(time.Time)
rq.AddErrorLog(error)
rq.AddTxStep(rdb.Query)
err := rq.StartTransaction(numSteps int)
err := rq.CommitTransaction()
output := rq.Output()
```

```
srq := rq.SubRequest()
rq.MergeLogs(srq)
```

### _type_: Schema[T]
Schema object for given type 

* ValidatorFn = func(any) bool
* TransformFn = func(any) any

```
schema, err := NewSchema[T](&T{}, table)
schema, err := NewSharedSchema[T](&T{})
AddRequiredField(schema, &item.Field)
AddEditableField(schema, &item.Field)
AddTransformer(schema, &item.Field, transformKey)
AddTransformFn(schema, &item.Field, TransformFn)

var validator ValidatorFn 
validator = NewStringValidator(func(string) bool)
AddValidator(schema, &item.Field, validator)
```

Available transformer keys:
* upper
* lower
* upperdot 
* lowerdot