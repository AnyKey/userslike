package repository

import (
	"database/sql"
	"github.com/pkg/errors"
)

type Repository struct {
	Conn *sql.DB
}

type TrackSelect struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}
func (repo *Repository) GetTracks(track string, artist string) ([]TrackSelect, error) {

	var trackList []TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and track.name = $1 AND artist.name = $2", track, artist)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}

	return trackList, nil
}
