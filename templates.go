package mediaparty

import (
	"text/template"
)

// writableString is a string that implements the Write method.
type writableString string

// Write adds whatever is written to the end of the string
func (self *writableString) Write(p []byte) (int, error) {
	*self += writableString(p)
	return len(p), nil
}

// Reset clears the current string
func (self *writableString) Reset() {
	*self = ""
}

// templateDefinition is a string that can be used as a go template
type templateDefinition string

// SelectColumnsTemplate defines the SELECT part of a query. Takes struct with Columns field.
var SelectColumnsTemplate templateDefinition = `
	Select {{parseColumns .Columns}}
	`

// FromTableTemplate defines the FROM part of a query. Takes struct with Schema and TableName fields.
var FromTableTemplate templateDefinition = `
	FROM {{.Schema}}.{{.TableName}}
	`

// SongWhereTemplate defines the WHERE clause for a song query. Takes a song struct.
var SongWhereTemplate templateDefinition = `
	Where
		LOWER("Artist") = LOWER('{{.Artist}}') AND
		LOWER("Title")  = LOWER('{{.Title}}') AND
		LOWER("Album")  = LOWER('{{.Album}}')
	`

// funcMap stores functions available for use in template strings
var funcMap template.FuncMap = template.FuncMap{
	"parseColumns": parseColumns,
}

// ExecuteTemplate fills in the given template with the given object and returns the final string.
func ExecuteTemplate(templateString templateDefinition, object interface{}) (string, error) {
	temp := template.Must(template.New("temp").Funcs(funcMap).Parse(string(templateString)))
	resultReader := new(writableString)
	err := temp.Execute(resultReader, object)
	if err != nil {
		return "", nil
	}

	return string(*resultReader), nil
}

// ParseColumns transforms an array of strings into a string of column names formatted for sql. Returns * for an empty list.
func parseColumns(columns []string) string {
	if len(columns) == 0 {
		return "*"
	}

	parsedString := columns[0]
	for i := 1; i < len(columns); i++ {
		parsedString += ", " + columns[i]
	}
	return parsedString
}
