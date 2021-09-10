package source

type SourceFieldType string

const (
	FIELDTYPE_TIMESTAMP SourceFieldType = "timestamp"
	FIELDTYPE_STRING    SourceFieldType = "string"
	FIELDTYPE_INT32     SourceFieldType = "int"
	FIELDTYPE_FLOAT32   SourceFieldType = "float"
)

type SourceField struct {
	Name  string          `json:"name"`
	Type  SourceFieldType `json:"type"`
	Range []interface{}   `json:"range,omitempty"`
	Limit []interface{}   `json:"limit,omitempty"`
}

type GeneratorSource struct {
	Name           string        `json:"name"`
	TimestampField string        `json:"timestamp_field"`
	BatchSize      int           `json:"batch_size"`
	Internval      []int         `json:"interval"`
	Fields         []SourceField `json:"fields"`
}
