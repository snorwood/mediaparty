package mediaparty

import (
	"text/template"
)

type writableString string

func (self *writableString) Write(p []byte) (int, error) {
	*self += writableString(p)
	return len(p), nil
}

func (self *writableString) Reset() {
	*self = ""
}

type templateDefinition string

var SelectColumnsTemplate templateDefinition = `
	Select {{parseColumns .Columns}}
	`

var FromTableTemplate templateDefinition = `
	FROM {{.Schema}}.{{.Table}}
	`

var SongWhereTemplate templateDefinition = `
	Where
		LOWER("Artist") = LOWER('{{.Artist}}') AND
		LOWER("Title")  = LOWER('{{.Title}}') AND
		LOWER("Album")  = LOWER('{{.Album}}')
	`
var funcMap template.FuncMap = template.FuncMap{
	"parseColumns": parseColumns,
}

func ExecuteTemplate(templateString templateDefinition, object interface{}) (string, error) {
	temp := template.Must(template.New("temp").Funcs(funcMap).Parse(string(templateString)))
	resultReader := new(writableString)
	err := temp.Execute(resultReader, object)
	if err != nil {
		return "", nil
	}

	return string(*resultReader), nil
}

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
