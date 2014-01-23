package mediaparty

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

// mp3Handler responds to request for mp3 files
func Mp3Handler(w http.ResponseWriter, r *http.Request) {
	// Connect to database -Move to function or outside requests?
	db, dberr := sql.Open("postgres",
		"user=postgres password=9K2Po2Tg4Es dbname=musicplayer sslmode=disable port=5432")

	defer db.Close()

	if dberr != nil {
		log.Fatal("Error connecting to database: ", dberr)
	}

	// Parse the path
	path := strings.Split(r.URL.Path, "/")

	if len(path) >= 6 {
		songInfo := Song{
			Artist:      path[2],
			Title:       path[3],
			Album:       path[4],
			AlbumArtist: path[5],
		}

		// Form the query for the song request
		query, err := GetSongQuery("musicplayer", "music", songInfo)
		if err != nil {
			log.Fatalf("Error getting query: %s", err)
		}

		// Retrieve the song from the database.
		songRow := db.QueryRow(query)
		song, scanErr := ScanSongFromRow(songRow)
		if scanErr != nil {
			log.Fatalf("Error scanning row: %s", scanErr)
		}
		fmt.Println("Served", song.Filepath)
		http.ServeFile(w, r, song.Filepath)
	}
}

// mainHandler responds to requests to the root path
func MainHandler(w http.ResponseWriter, r *http.Request) {
	// Check path for valid song info
	path := strings.Split(r.URL.Path, "/")
	if len(path) >= 5 {
		song := Song{
			Artist:      path[1],
			Title:       path[2],
			Album:       path[3],
			AlbumArtist: path[4],
		}

		// Embed song in page
		html, err := ExecuteTemplate(MusicPlayer(), song)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(html)
			fmt.Fprint(w, html)
		}
	} else {
		fmt.Fprint(w, "THIS IS THE HOMEPAGE CONGRATS -_-")
	}
}

// tableHandler will return an object for listing songs
func TableHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tabletest.html")
}

// queryHandler handles querys for a list of filtered songs
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to database
	db, dberr := sql.Open("postgres",
		"user=postgres password=9K2Po2Tg4Es dbname=musicplayer sslmode=disable port=5432")

	defer db.Close()

	if dberr != nil {
		log.Fatal("Error connecting to database: ", dberr)
	}

	// Parse the path for the requested song info
	rawFilters := r.URL.Query()
	filters := make(map[string]string)
	for key, value := range rawFilters {
		if len(value) >= 1 {
			filters[strings.ToLower(key)] = value[0]
		}
	}

	fmt.Println(filters)

	// Request the songs using the given info as a filter
	songInfo := Song{
		Artist:      filters["artist"],
		Title:       filters["title"],
		Album:       filters["album"],
		AlbumArtist: filters["albumartist"],
	}

	query, err := VariableSongQuery("musicplayer", "music", songInfo)
	if err != nil {
		// log.Fatalf("Error getting query: %s", err)
		fmt.Printf("Error getting query: %s\n", err)
		return
	}

	// Recieve the song list
	songRows, dbErr := db.Query(query)
	if dbErr != nil {
		// log.Fatalf("Invalid response from database: %s \n Query: %s", dbErr, query)
		fmt.Printf("Invalid response from database: %s \n Query: %s\n", dbErr, query)
		return
	}
	songs := make([]*Song, 0)
	for songRows.Next() {
		song, scanErr := ScanSongFromRow(songRows)

		// Security measure. Don't want this being sent out.
		song.Filepath = ""

		if scanErr != nil {
			log.Fatalf("Error scanning row: %s", scanErr)
		} else {
			songs = append(songs, song)
		}
	}
	fmt.Println(query, songs)

	// Convert song list to json
	responseBytes, jsonErr := json.Marshal(songs)
	response := string(responseBytes)

	// Send the song list
	if jsonErr != nil {
		fmt.Fprintf(w, jsonErr.Error())
		log.Fatalf("Error converting to json: %s", jsonErr)
	} else {
		fmt.Println("Served", response)
		fmt.Fprintf(w, response)
	}
}
