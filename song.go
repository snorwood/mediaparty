package mediaparty

import (
	"fmt"
)

// Song base song class. Contains metadata and file location of song.
type Song struct {
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
	Filepath    string
	ID          int
}

// Valid returns whether the song has all of the neccessary metadata
func (self Song) Valid() bool {
	if self.Artist == "" || self.Album == "" || self.AlbumArtist == "" || self.Title == "" {
		return false
	}

	return true
}

// ScanFromRow reads the data from the sql row into a song struct.
func ScanSongFromRow(row rowScanner) (*song, error) {
	// Scan requires pointers to variables to write into so define them here.
	song = new(Song)
	artist := &song.Artist
	album := &song.Album
	title := &song.Title
	albumArtist := &song.AlbumArtist
	var id *int = new(int)
	filepath := &song.Filepath

	// Perform the scan
	err := row.Scan(artist, title, album, albumArtist, filepath, id)
	if err != nil {
		return nil, err
	}

	// song.Artist = *artist
	// song.Title = *title
	// song.Album = *album
	// song.AlbumArtist = *albumArtist
	// song.Filepath = *filepath
	// song.ID = *id
	// if err != nil {
	// 	song.ID = -1
	// 	return err
	// }

	return song, nil
}

// InvalidSong is a Song that implements Error()
type InvalidSong Song

// Error displays all of the metadata attributes
func (self InvalidSong) Error() string {
	return fmt.Sprintf("Invalid song:\nArtist=%s\nTitle=%s\nAlbum=%s\nAlbumArtist=%s", self.Artist, self.Title, self.Album, self.AlbumArtist)
}
