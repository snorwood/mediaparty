package mediaparty

import (
	"fmt"
)

// rowScanner allows items to be read from either a *sql.Rows or *sql.Row object
type rowScanner interface {
	Scan(...interface{}) error
}

// Table contains the table's indentity within the database
type Table struct {
	Schema    string
	TableName string
}

// GetSongQuery builds a query to retrieve the given song
func GetSongQuery(songTable Table, columnNames []string, song Song) (string, error) {

	if !song.Valid() {
		return "", InvalidSong(song)
	}

	if songTable.Schema == "" || songTable.Tableable == "" {
		return "", fmt.Errorf("Invalid Table:\nschema=%s\ntable=%s", schema, table)
	}

	columns := struct {
		Columns []string
	}{
		Columns: columnNames,
	}

	querySegment := ""
	var err error = nil

	query := ""

	querySegment, err = ExecuteTemplate(SelectColumnsTemplate, columns)
	query += querySegment
	querySegment, err = ExecuteTemplate(FromTableTemplate, songTable)
	query += querySegment
	querySegment, err = ExecuteTemplate(SongWhereTemplate, song)
	query += querySegment

	if err != nil {
		return "", fmt.Errorf("Could not put together query: %s", err)
	}

	return query, nil
}

// StringToPostgresString fixes all of the syntax from a normal
// string that will mess up a postgresql query
func StringToPostgresString(s string) string {
	for i := 0; i < len(s); i++ {
		escapeCharNext := (i < len(s)-1 && (string(s[i+1]) == "\\" || string(s[i+1]) == "'"))
		if string(s[i]) == "\\" {
			if !escapeCharNext {
				s = string(s[:i]) + "\\" + string(s[i:])
			}

			i += 1
		} else {
			if string(s[i]) == "'" {
				s = string(s[:i]) + "\\" + string(s[i:])
				i += 1
			}
		}
	}

	return s
}
