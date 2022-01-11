package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Bookmark struct {
	description string
	extended    string
	hash        string
	href        string
	meta        string
	shared      string
	tags        string
	time        string
	toread      string
}

func LoadBookmarkTemplate(fpath string) (*template.Template, error) {
	content, err := os.ReadFile("./template.txt")

	tmpl := template.New("bookmark")

	if err != nil {
		return tmpl, err
	}

	return tmpl.Parse(string(content))
}

/*
 * Write a bookmark into a file
 *
 */
func (book *Bookmark) Write(vault *ObsidianVault, template *template.Template) error {
	date := time.Now().Format("20060101") + fmt.Sprintf("%04d", rand.Intn(10000))

	view := struct {
		Date        string
		Description string
		Extended    string
		Hash        string
		Href        string
		Meta        string
		Shared      string
		Tags        string
		Time        string
		Toread      string
	}{
		Date:        date,
		Description: book.description,
		Extended:    book.extended,
		Hash:        book.hash,
		Href:        book.href,
		Meta:        book.meta,
		Shared:      book.shared,
		Tags:        book.tags,
		Time:        book.time,
		Toread:      book.toread,
	}

	buf := new(bytes.Buffer)
	if err := template.Execute(buf, view); err != nil {
		return err
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9- |]+")
	if err != nil {
		return err
	}

	fragment := reg.ReplaceAllString(book.description, "")
	limit := len(fragment) - 1

	if limit > 80 {
		limit = 80
	}

	fname := date + " - " + fragment[0:limit] + ".md"
	fpath := filepath.Join(vault.dpath, "pinboard-bookmarks", fname)

	err = os.MkdirAll(filepath.Join(vault.dpath, "pinboard-bookmarks"), 0700)
	if err != nil {
		return errors.Wrap(err, "failed creating pinboard-bookmarks")
	}

	err = os.WriteFile(fpath, buf.Bytes(), 0700)
	if err != nil {
		return errors.Wrap(err, "failed writing bookmark to file")
	}

	return nil
}

// parse yaml
type Frontmatter struct {
	BookmarkHash string `yaml:"bookmark_hash"`
}

/*
 * Enumerate all bookmark hashes present in Obsidian frontmatter, so we
 * can efficiently find bookmark pages that still need to be created
 *
 */
func GetPresentBookmarkHashes(conn *OpalDb) ([]string, error) {
	rows, err := conn.Db.Query(`select file_id, content from metadata`)
	hashes := []string{}

	if err != nil {
		return hashes, err
	}

	for rows.Next() {
		var content string
		var fileId string

		err := rows.Scan(&fileId, &content)

		if err != nil {
			return hashes, err
		}

		fm := Frontmatter{}
		err = yaml.Unmarshal([]byte(content), &fm)
		if err != nil {
			continue
		}

		if len(fm.BookmarkHash) > 0 {
			hashes = append(hashes, fm.BookmarkHash)
		}
	}

	return hashes, nil
}

/*
 * Sync pinboard bookmarks into Obsidian, ensuring each hash exists as a
 * metadata header
 */
func SyncBookmarks(templatePath string, vault *ObsidianVault, conn *OpalDb) error {
	tmpl, err := LoadBookmarkTemplate(templatePath)
	if err != nil {
		return err
	}

	hashes, err := GetPresentBookmarkHashes(conn)
	if err != nil {
		return err
	}

	bookmarks, err := conn.ListAbsentBookmarks(hashes)
	if err != nil {
		return err

	}
	fmt.Println(len(bookmarks))

	// write bookmarks to Obsidian
	for _, bookmark := range bookmarks {
		err = bookmark.Write(vault, tmpl)

		if err != nil {
			return err
		}
	}

	return nil
}
