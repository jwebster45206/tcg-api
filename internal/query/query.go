package query

type FieldType string

const (
	FieldTypeString   FieldType = "string"
	FieldTypeUUID     FieldType = "uuid"
	FieldTypeInt      FieldType = "int"
	FieldTypeFloat    FieldType = "float"
	FieldTypeBool     FieldType = "bool"
	FieldTypeDateTime FieldType = "datetime"
)

type FilterOperator string

const (
	OpEqual        FilterOperator = "="  // Auto converts to "IN" for arrays, "IS NULL" for nil
	OpNotEqual     FilterOperator = "!=" // Auto converts to "NOT IN" for arrays, "IS NOT NULL" for nil
	OpGreaterThan  FilterOperator = ">"
	OpGreaterEqual FilterOperator = ">="
	OpLessThan     FilterOperator = "<"
	OpLessEqual    FilterOperator = "<="
	OpLike         FilterOperator = "LIKE"
	OpNotLike      FilterOperator = "NOT LIKE"
)

// Filter represents a single SQL filter condition
type Filter struct {
	Column   string
	Operator FilterOperator
	Value    interface{} // nil, single value, or slice
}

// SortOption represents a field to sort by and its direction
type SortOption struct {
	Field string
	Desc  bool
}

type QueryConfig struct {
	AllowedFilters map[string]string    // map[apiField]dbColumn
	AllowedSorts   map[string]string    // map[apiField]dbColumn
	FieldTypes     map[string]FieldType // map[apiField]fieldType
}

func (c QueryConfig) IsFilterAllowed(field string) bool {
	_, exists := c.AllowedFilters[field]
	return exists
}

func (c QueryConfig) IsSortAllowed(field string) bool {
	_, exists := c.AllowedSorts[field]
	return exists
}

func (c QueryConfig) GetFilterDBColumn(field string) (string, bool) {
	dbCol, exists := c.AllowedFilters[field]
	return dbCol, exists
}

func (c QueryConfig) GetSortDBColumn(field string) (string, bool) {
	dbCol, exists := c.AllowedSorts[field]
	return dbCol, exists
}

func (c QueryConfig) GetFieldType(field string) FieldType {
	fieldType, exists := c.FieldTypes[field]
	if !exists {
		return FieldTypeString // Default to string if not specified
	}
	return fieldType
}
