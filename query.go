package mediaparty

import (
	"fmt"
)

// rowScanner encompasses all row objects with scan functionality
type rowScanner interface {
	Scan(...interface{}) error
}

// Table contains the schema.table data
type Table struct {
	Schema string
	Table  string
}

// GetSongQuery builds the query for a single song
func GetSongQuery(schema, table string, song Song) (string, error) {
	// Make sure the song info is valid
	if !song.Valid() {
		return "", InvalidSong(song)
	}

	// Create/validate the table object
	songTable := Table{
		Schema: schema,
		Table:  table,
	}
	if schema == "" || table == "" {
		return "", fmt.Errorf("Invalid Table:\nschema=%s\ntable=%s", schema, table)
	}

	// Create the columns list
	columns := struct {
		Columns []string
	}{
		Columns: []string{},
	}

	// Build up the query
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

// VariableSongQuery creates a query out of only the given song data.
func VariableSongQuery(schema, table string, song Song) (string, error) {
	// Build and verify the table
	songTable := Table{
		Schema: schema,
		Table:  table,
	}
	if schema == "" || table == "" {
		return "", fmt.Errorf("Invalid Table:\nschema=%s\ntable=%s", schema, table)
	}

	// Build the column list
	columns := struct {
		Columns []string
	}{
		Columns: []string{},
	}

	// Build the query up
	querySegment := ""
	var err error = nil

	query := ""

	querySegment, err = ExecuteTemplate(SelectColumnsTemplate, columns)
	query += querySegment
	querySegment, err = ExecuteTemplate(FromTableTemplate, songTable)
	query += querySegment

	// Add filter data Artist -> Album -> Title
	query += "WHERE "

	if song.Artist == "" {
		return query, err
	}

	query += formatColumn("Artist") + " = " + formatValue(song.Artist)

	if song.Album == "" {
		return query, err
	}

	query += " AND " + formatColumn("Album") + " = " + formatValue(song.Album)

	if song.Title == "" {
		return query, err
	}

	query += " AND " + formatColumn("Title") + " = " + formatValue(song.Title)

	return query, err
}

// StringToPostgresString converts a normal string into a postgres compatible one.
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
