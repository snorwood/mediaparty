package mediaparty

import (
	"fmt"
	"io/ioutil"
	"text/template"
)

// writableString is a string that can be wrote to using writers
type writableString string

func (self *writableString) Write(p []byte) (int, error) {
	*self += writableString(p)
	return len(p), nil
}

func (self *writableString) Reset() {
	*self = ""
}

// templateDefinition is a string type that specifies it is to be used for templates
type templateDefinition string

// SelectColumnsTemplate defines the columns to select
var SelectColumnsTemplate templateDefinition = `
	Select {{parseColumns .Columns}}
	`

// FromTableTemplate defines the table to select from
var FromTableTemplate templateDefinition = `
	FROM {{.Schema}}.{{.Table}}
	`

// SongWhereTemplate defines the song info
var SongWhereTemplate templateDefinition = `
	Where
		LOWER("Artist") = LOWER('{{.Artist}}') AND
		LOWER("Title")  = LOWER('{{.Title}}') AND
		LOWER("Album")  = LOWER('{{.Album}}')
	`

// MusicPlayer returns the html template for an embedded music player
func MusicPlayer() templateDefinition {
	content, err := ioutil.ReadFile("../musicplayer.html")
	html := string(content)

	if err != nil {
		fmt.Printf("Could not find the musicplayer definition: %s", err)
	}

	return templateDefinition(html)
}

// funcMap maps function to be used in templates
var funcMap template.FuncMap = template.FuncMap{
	"parseColumns": parseColumns,
}

// ExecuteTemplate returns a copy of the template filled in with the object.
func ExecuteTemplate(templateString templateDefinition, object interface{}) (string, error) {
	temp := template.Must(template.New("temp").Funcs(funcMap).Parse(string(templateString)))
	resultReader := new(writableString)
	err := temp.Execute(resultReader, object)
	if err != nil {
		return "", nil
	}

	return string(*resultReader), nil
}

// parseColumns creates a Postgres formatted list of columns. * if no columns passed.
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

// formatColumn provides Postgres column definition formatting to a normal string
func formatColumn(column string) string {
	return "\"" + column + "\""
}

// formatValue provides Postgres string value formatting to a normal string
func formatValue(value string) string {
	return "'" + value + "'"
}
