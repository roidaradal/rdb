# RDB
A Go library for type-safe SQL database queries

`go get github.com/roidaradal/rdb/...`

## Initialization 

### Initialize
Required call to initialize rdb library

`err := rdb.Initialize()`

### Type Definition
By default, column name = field name; fields can be skipped or mapped to a custom column name

```
type Foo struct {
    Field1      string                      // column name: Field1
    Field2      string  `col:"CustomName"`  // column name is CustomName, instead of Field2
    PrivField_  int     `col:"-"`           // skipped
}
```

### AddType
Adds new type, pass in a struct pointer

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
Get all column names of given item

```
columns := rdb.AllColumns(Type{})
columns := rdb.AllColumns(&Type{})
```

### Column
Get column name of given field pointer,
Field must be from the object used in AddType

`column := rdb.Column(&item.Field)`

### Columns
Get column names of given field pointers, 
Fields must be from the object used in AddType 

`columns := rdb.Columns(&item.Field1, &item.Field2, &item.Field3)`

### _type:_ RowReader[T]
Function that reads row values into object

### Reader[T]
Creates a RowReader[T] with the given columns 

`reader := rdb.Reader[T](column1, column2, ...)`

### FullReader[T]
Creates a RowReader[T] using all columns of type T

`reader := rdb.FullReader(&T{})`

### ToRow[T]
Converts given object to map[string]any for row insertion 

`row := rdb.ToRow(&item)`

## Conditions

### _interface:_ Condition
Unifies the different Condition types

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
`condition := rdb.And(condition1, condition2, condition3)`

### Or
`condition := rdb.Or(condition1, condition2, condition3)`

### NoCondition 
Match all condition 

`condition := rdb.NoCondition()`

## Queries 

### _interface:_ Query 
Unifies the different Query types

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
Creates a new TopRowQuery

```
q := rdb.NewTopRowQuery(table, reader)
q.Where(condition)
q.OrderAsc(rdb.Column(&item.Field)) // or 
q.OrderDesc(rdb.Column(&item.Field))
topItem, err := q.QueryRow(*sql.DB)
```

### NewTopValueQuery
Creates a new TopValueQuery

```
var topValue V
q := rdb.NewTopValueQuery(table, &item.Field)
q.Where(condition)
q.OrderAsc(rdb.Column(&item.Field)) // or 
q.OrderDesc(rdb.Column(&item.Field))
topValue, err := q.QueryValue(*sql.DB)
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

### _type:_ QueryResultChecker 
Function that checks SQL result if a condition has been satisfied

### RowsAffected
Gets the number of rows affected from SQL result (default: 0)

`affected := rdb.RowsAffected(*sql.Result)`

### LastInsertID 
Gets the last insert ID (uint) from SQL result (default: 0)

`id, ok := rdb.LastInsertID(*sql.Result)`

### AssertNothing
QueryResultChecker that does nothing

`checker := rdb.AssertNothing`

### AssertRowsAffected
Creates a QueryResultChecker that asserts the number of rows affected

`checker := rdb.AssertRowsAffected(1)`

### Exec
Executes an SQL query

`result, err := rdb.Exec(q, *sql.DB)`

### ExecTx
Executes an SQL query as part of a transaction, applies Rollback on any errors

`result, err := rdb.ExecTx(q, *sql.Tx, checker)`

### Rollback 
Rolls back the SQL transaction

`err := rdb.Rollback(*sql.Tx, err)`
