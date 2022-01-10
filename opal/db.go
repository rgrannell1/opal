package main

import (
	"database/sql"
)

type OpalDb struct {
	Db *sql.DB
}

/*
 * Create opal metadata table
 *
 */
func (conn *OpalDb) CreateTables() error {
	tx, err := conn.Db.Begin()
	defer tx.Rollback()

	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS opal_metadata (
		id                TEXT NOT NULL,
		processed_hash    TEXT NOT NULL,

		PRIMARY KEY(id)
	)`)

	if err != nil {
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

/*
 * Construct a database wrapper
 *
 */
func NewOpalDb(fpath string) (*OpalDb, error) {
	db, err := sql.Open("sqlite3", fpath)
	if err != nil {
		return &OpalDb{}, err
	}

	return &OpalDb{db}, nil
}

/*
 *
 */
func (conn *OpalDb) ReadTitle() (string, error) {
	return "", nil
}

/*
 *
 */
func (conn *OpalDb) WriteProcessedHash() (string, error) {
	return "", nil
}

/*
 * Fetch the file hash, and opal hash
 *
 */
func (conn *OpalDb) GetHashes(fpath string) (string, string, error) {
	var hash string
	row := conn.Db.QueryRow(`SELECT hash FROM file WHERE file.id = ?`, fpath)

	err := row.Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil
		}
	}

	var processedHash string
	opalRow := conn.Db.QueryRow(`SELECT processed_hash FROM opal_metadata WHERE id = ?`, fpath)

	if opalRow == nil {
		return hash, "", nil
	}

	err = opalRow.Scan(&processedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil
		}
	}

	return hash, processedHash, nil
}

func (conn *OpalDb) MarkComplete(notes []*ObsidianNote) error {
	return nil
}

func (conn *OpalDb) GetFrontmatter() error {
	return nil
}

func (conn *OpalDb) ListBookmarks() ([]*Bookmark, error) {
	bookmarks := make([]*Bookmark, 0)

	rows, err := conn.Db.Query(`SELECT description, extended, hash, href, meta, shared, tags, time, toread FROM pinboard_bookmark`)
	if err != nil {
		return bookmarks, err
	}

	for rows.Next() {
		bookmark := Bookmark{}
		err := rows.Scan(
			&bookmark.description,
			&bookmark.extended,
			&bookmark.hash,
			&bookmark.href,
			&bookmark.meta,
			&bookmark.shared,
			&bookmark.tags,
			&bookmark.time,
			&bookmark.toread)

		if err != nil {
			return bookmarks, err
		}

		bookmarks = append(bookmarks, &bookmark)
	}

	err = rows.Close()
	if err != nil {
		return bookmarks, err
	}

	return bookmarks, err
}
