package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

type Repository struct {
	Conn *sql.DB
}

type TrackSelect struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Username string `json:"username"`
}
type LikeSelect struct {
	Username string `json:"username"`
}

func (repo *Repository) GetTracks(track string, artist string) ([]LikeSelect, *string, error) {

	var trackList []LikeSelect
	rows, err := repo.Conn.Query("SELECT like_list.username FROM track_list, like_list "+
		"WHERE like_list.track_id = track_list.id AND track_list.name = $1 AND track_list.artist = $2", track, artist)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := LikeSelect{}
		err := rows.Scan(&tl.Username)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	rows, err = repo.Conn.Query("SELECT count(like_list.username) AS track_count, track_list.name, track_list.artist "+
		"WHERE like_list.track_id = track_list.id AND track_list.name = $1 AND track_list.artist = $2 GROUP BY track_list.name, track_list.artist", track, artist)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()
	var likeCount *string
	for rows.Next() {

		err := rows.Scan(&likeCount)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error Scan values")
		}
	}

	return trackList, likeCount, nil
}

func (repo *Repository) SetLike(name, artist, username string) error {
	ctx := context.Background()
	tx, err := repo.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	// defer commit rollback tnx
	var lastinsertedid int
	rows, err := tx.QueryContext(ctx, "SELECT id FROM track_list WHERE name = $1 AND artist =$2", name, artist)
	if err != nil {
		return errors.Wrap(err, "error select in DB!")
	}
	var trackId int
	if rows != nil {
		var tl int
		for rows.Next() {
			err = rows.Scan(&tl)
			if err != nil {
				return errors.Wrap(err, "error Scan values")
			}
		}
		trackId = tl
	}
	if trackId == 0 {
		err = tx.QueryRowContext(ctx, "INSERT INTO track_list (name, artist) VALUES ($1, $2) returning id", name, artist).Scan(&lastinsertedid)
		if err != nil {
			tx.Rollback()
			fmt.Println("TRACK!", err.Error())
			return err
		}
		trackId = lastinsertedid
	}
	rows.Close()
	if trackId == 0 {
		tx.Rollback()
		return errors.New("Empty track")
	}

	rows, err = tx.QueryContext(ctx, "SELECT id FROM like_list WHERE username = $1 AND track_id =$2", username, trackId)
	if err != nil {
		return errors.Wrap(err, "error select in DB!")
	}
	if rows != nil {
		var tl int
		for rows.Next() {
			err = rows.Scan(&tl)
			if err != nil {
				return errors.Wrap(err, "error Scan values")
			}
		}
		if tl != 0 {
			return nil
		}

	}

	_, err = tx.ExecContext(ctx, "INSERT INTO like_list (username, track_id) VALUES ($1, $2)", username, trackId)
	if err != nil && err == sql.ErrNoRows {
		fmt.Println("LIKE!", err.Error())
		tx.Rollback()
		return nil
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	return nil
}
