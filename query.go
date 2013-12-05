package mediaparty

import (
	"fmt"
)

type rowScanner interface {
	Scan(...interface{}) error
}

type Table struct {
	Schema string
	Table  string
}

func GetSongQuery(schema, table, artist, title, album, albumArtist string) (string, error) {
	song := Song{
		Artist:      artist,
		Title:       title,
		Album:       album,
		AlbumArtist: albumArtist,
	}

	if !song.Valid() {
		return "", InvalidSong(song)
	}

	songTable := Table{
		Schema: schema,
		Table:  table,
	}

	if schema == "" || table == "" {
		return "", fmt.Errorf("Invalid Table:\nschema=%s\ntable=%s", schema, table)
	}

	columns := struct {
		Columns []string
	}{
		Columns: []string{},
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
