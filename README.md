# RDB
A Go library for type-safe SQL database queries.

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