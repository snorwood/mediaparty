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

var SelectColumnsTemplate string = `
	Select {{parseColumns .Columns}}
	`

func ParseColumns(columns []string) string {
	if len(columns) == 0 {
		return "*"
	}

	parsedString := columns[0]
	for i := 1; i < len(columns); i++ {
		parsedString += ", " + columns[i]
	}
	return parsedString
}

var funcMap template.FuncMap = template.FuncMap{
	"parseColumns": ParseColumns,
}

var FromTableTemplate string = `
	FROM {{.Schema}}.{{.Table}}
	`

var SongWhereTemplate string = `
	Where
		LOWER(Artist) = LOWER({{.Artist}}) AND
		LOWER(Title)  = LOWER({{.Title}}) AND
		LOWER(Album)  = LOWER({{.Album}})
	`

func ExecuteTemplate(templateDefinition string, object interface{}) (string, error) {
	temp := template.Must(template.New("temp").Funcs(funcMap).Parse(templateDefinition))
	resultReader := new(writableString)
	err := temp.Execute(resultReader, object)
	if err != nil {
		return "", nil
	}

	return string(*resultReader), nil
}
