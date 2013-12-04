package mediaparty

import (
	"fmt"
)

type SongQuery struct {
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
}

func (self *SongQuery) GetWhere() string {
	return fmt.Sprintf("WHERE TITLE = %s AND Album = %s AND (Artist = %s OR AlbumArtist = %s)",
		StringToPostgresString(self.Title),
		StringToPostgresString(self.Album),
		StringToPostgresString(self.Artist),
		StringToPostgresString(self.AlbumArtist),
	)
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
