package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Track describes a individual song in the library or playlist.
type Track struct {
	Title    string
	Album    string
	Artist   string
	Duration int
	Removed  bool
	Index    int
}

func (t Track) String() string {
	var str strings.Builder
	str.WriteString(t.Title)
	str.WriteString(t.Album)
	str.WriteString(t.Artist)
	return str.String()
}

// Tracks is a slice of tracks; used to define custom sort
type Tracks []Track

func (t Tracks) Len() int {
	return len(t)
}

func (t Tracks) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t Tracks) Less(a, b int) bool {
	aValue := t[a].String()
	bValue := t[b].String()
	return aValue < bValue
}

// Entry point
func main() {
	sourceArg := flag.String("source", ".", "location of the Google Takeout export")
	destArg := flag.String("dest", ".", "where to store the flattened files")
	flag.Parse()

	flattenLibrary(*sourceArg, *destArg)
	flattenPlaylists(*sourceArg, *destArg)
	fmt.Println("Done.")
}

// Flatten the library collection
func flattenLibrary(source string, dest string) {
	source = path.Join(source)
	flattenPlaylist("main library", source, path.Join(dest, "Library"))
}

// Find and flatten all of the playlists
func flattenPlaylists(source string, dest string) {
	playlistsFolder := path.Join(source, "Playlists")
	files, err := ioutil.ReadDir(playlistsFolder)
	checkError("Cannot read from "+playlistsFolder, err)

	for _, file := range files {
		playlistSource := path.Join(playlistsFolder, file.Name())
		playlistDest := path.Join(dest, "Playlists", file.Name())
		flattenPlaylist(file.Name(), playlistSource, playlistDest)
	}
}

// Flatten an individual playlist
func flattenPlaylist(name string, source string, dest string) {
	fmt.Println("Flattening", name)

	if name != "Thumbs Up" {
		source = path.Join(source, "Tracks")
	}

	tracks := readTrackCollection(source)
	sort.Sort(Tracks(tracks))
	writeTrackCollection(dest+".csv", tracks)
}

// Load a track collection, e.g. the library or a playlist
func readTrackCollection(collectionPath string) []Track {
	var tracks []Track

	files, err := ioutil.ReadDir(collectionPath)
	checkError("Cannot read from "+collectionPath, err)

	for _, file := range files {
		trackPath := path.Join(collectionPath, file.Name())
		track := readSingleTrack(trackPath)
		tracks = append(tracks, track)
	}

	return tracks
}

// Load one of GPM's tracks CSV files, each of which contains the information for a single song.
func readSingleTrack(trackPath string) Track {
	data, err := os.Open(trackPath)
	checkError("Cannot read from "+trackPath, err)

	defer data.Close()

	reader := csv.NewReader(bufio.NewReader(data))
	reader.Read()              // header line
	line, err := reader.Read() // actual data
	checkError("Cannot read from "+trackPath, err)

	duration, _ := strconv.ParseInt(line[3], 10, 64)

	var index int64
	if len(line) > 7 {
		index, _ = strconv.ParseInt(line[7], 10, 64)
	}

	return Track{
		Title:    line[0],
		Album:    line[1],
		Artist:   line[2],
		Duration: int(duration),
		Removed:  ("Yes" == line[6]),
		Index:    int(index),
	}
}

// Write out the collection CSV, which now contains all of the track info
func writeTrackCollection(collectionPath string, tracks []Track) {
	os.MkdirAll(path.Dir(collectionPath), os.ModePerm)

	file, err := os.Create(collectionPath)
	checkError("Cannot create file at"+collectionPath, err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Title", "Album", "Artist", "Duration", "Index"})
	checkError("Cannot write to file", err)

	for _, track := range tracks {
		if !track.Removed {
			line := []string{
				track.Title,
				track.Album,
				track.Artist,
				strconv.Itoa(track.Duration),
				strconv.Itoa(track.Index),
			}

			err := writer.Write(line)
			checkError("Cannot write to file", err)
		}
	}
}

// Test and exit on error
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
