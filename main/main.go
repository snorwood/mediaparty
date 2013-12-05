package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/snorwood/mediaparty"
	"log"
	"net/http"
	"strings"
)

func mp3Handler(w http.ResponseWriter, r *http.Request) {
	db, dberr := sql.Open("postgres",
		"user=postgres password=9K2Po2Tg4Es dbname=musicplayer sslmode=disable port=5432")
	defer db.Close()
	if dberr != nil {
		log.Fatal("Error connecting to database: ", dberr)
	}

	s := strings.Split(r.URL.Path, "/")
	fmt.Println(s)
	song := mediaparty.Song{
		Artist:      s[2],
		Title:       s[3],
		Album:       s[4],
		AlbumArtist: s[5],
	}

	query, err := mediaparty.GetSongQuery(Table{"musicplayer", "music"}, []string{}, song)
	if err != nil {
		log.Fatalf("Error getting query: %s", err)
	}

	songRow := db.QueryRow(query)
	filledSong, scanErr := ScanSongFromRow(songRow)
	if scanErr != nil {
		log.Fatalf("Error scanning row: %s", scanErr)
	}
	fmt.Println("Served", song.Filepath)
	http.ServeFile(w, r, song.Filepath)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<body>
			<audio controls>
  				<source src="/music/StudioEIM/The Third Flight/MapleStory/Wizet" type="audio/mpeg">
  				<embed height="50" width="100" src="/music/StudioEIM/The Third Flight/MapleStory/Wize">
			</audio>
		</body>
	`)
}

func main() {
	http.HandleFunc("/music/", mp3Handler)
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8000", nil)
}
