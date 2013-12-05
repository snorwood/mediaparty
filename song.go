package mediaparty

import (
	"fmt"
)

type Song struct {
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
	Filepath    string
	ID          int
}

func (self Song) Valid() bool {
	if self.Artist == "" || self.Album == "" || self.AlbumArtist == "" || self.Title == "" {
		return false
	}

	return true
}

func (self *Song) ScanFromRow(row rowScanner) error {
	artist := &self.Artist
	album := &self.Album
	title := &self.Title
	albumArtist := &self.AlbumArtist
	var id *int = new(int)
	filepath := &self.Filepath
	err := row.Scan(artist, title, album, albumArtist, filepath, id)

	if err != nil {
		return err
	}

	self.Artist = *artist
	self.Title = *title
	self.Album = *album
	self.AlbumArtist = *albumArtist
	self.Filepath = *filepath
	self.ID = *id
	if err != nil {
		self.ID = -1
		return err
	}
	return nil
}

type InvalidSong Song

func (self InvalidSong) Error() string {
	return fmt.Sprintf("Invalid song:\nArtist=%s\nTitle=%s\nAlbum=%s\nAlbumArtist=%s", self.Artist, self.Title, self.Album, self.AlbumArtist)
}
