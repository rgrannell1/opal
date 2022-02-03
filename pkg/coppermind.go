package opal

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

/*
 * A Pinboard bookmark
 */
type PinboardBookmark struct {
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

func TitleCase(arg string) (interface{}, error) {
	return strings.Title(strings.ToLower(arg)), nil
}

func LoadBookmarkTemplate(fpath string) (*template.Template, error) {
	content, err := os.ReadFile(fpath)
	tmpl := template.New("bookmark").Funcs(template.FuncMap(map[string]interface{}{
		"TitleCase": TitleCase,
	}))

	if err != nil {
		return tmpl, err
	}

	return tmpl.Parse(string(content))
}

/*
 * Create a description from a bookmark. Modify sites with unpleasant
 * descriptions.
 */
func CreateDescription(book *PinboardBookmark) string {
	url, err := url.Parse(book.href)
	if err != nil {
		return book.description
	}

	if url.Host == "twitter.com" {
		parts := strings.Split(url.Path, "/")

		if len(parts) >= 2 {
			forbidden := regexp.MustCompile(":-")
			return "Tweet from " + parts[1] + " on " + forbidden.ReplaceAllString(book.time, "_")
		}
	}

	if url.Host == "en.wikipedia.org" {
		return strings.ReplaceAll(book.description, " - Wikipedia", "")
	}

	return book.description
}

/*
 * Get a filename for a bookmark
 */
func (book *PinboardBookmark) FileName() (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9- |]+")
	if err != nil {
		return "", err
	}

	fragment := reg.ReplaceAllString(CreateDescription(book), " ")
	limit := len(fragment)

	if limit > 128 {
		limit = 128
	}

	date := time.Now().Format("20060101") + fmt.Sprintf("%04d", rand.Intn(10000))
	fname := date + " - " + strings.Title(strings.ToLower(fragment[0:limit])) + ".md"

	return fname, nil
}

/*
 * Write a bookmark into a file
 *
 */
func (book *PinboardBookmark) Write(vault *ObsidianVault, template *template.Template) error {
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
		Description: CreateDescription(book),
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

	fname, err := book.FileName()
	if err != nil {
		return err
	}

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
	GithubRepo   string `yaml:"github_repo"`
}

/*
 * Enumerate all bookmark hashes present in Obsidian frontmatter, so we
 * can efficiently find bookmark pages that still need to be created
 *
 */
func GetFrontmatter(conn *OpalDb) ([]*Frontmatter, error) {
	rows, err := conn.Db.Query(`
	select file_id, content
	from metadata
	where schema = "!frontmatter"`)

	frontmatter := []*Frontmatter{}
	if err != nil {
		return frontmatter, err
	}

	for rows.Next() {
		var content string
		var fileId string

		// scan the content and fileId
		err := rows.Scan(&fileId, &content)
		if err != nil {
			return frontmatter, err
		}

		fm := Frontmatter{}
		err = yaml.Unmarshal([]byte(content), &fm)
		if err != nil {
			continue
		}

		frontmatter = append(frontmatter, &fm)
	}

	return frontmatter, nil
}

/*
 * Enumerate all bookmark hashes present in Obsidian frontmatter, so we
 * can efficiently find bookmark pages that still need to be created
 *
 */
func GetPresentBookmarkHashes(conn *OpalDb) (*Set, error) {
	set := NewSet([]string{})
	frontmatter, err := GetFrontmatter(conn)

	if err != nil {
		return set, err
	}

	for _, fm := range frontmatter {
		if len(fm.BookmarkHash) > 0 {
			set.Add(fm.BookmarkHash)
		}
	}

	return set, nil
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

	// write bookmarks to Obsidian
	for _, bookmark := range bookmarks {
		err = bookmark.Write(vault, tmpl)

		if err != nil {
			return err
		}
	}

	return nil
}

func LoadGithubStarTemplate(fpath string) (*template.Template, error) {
	content, err := os.ReadFile(fpath)
	tmpl := template.New("github").Funcs(template.FuncMap(map[string]interface{}{
		"TitleCase": TitleCase,
	}))

	if err != nil {
		return tmpl, err
	}

	return tmpl.Parse(string(content))
}

func GetPresentGithubStars(conn *OpalDb) (*Set, error) {
	set := NewSet([]string{})
	frontmatter, err := GetFrontmatter(conn)

	if err != nil {
		return set, err
	}

	for _, fm := range frontmatter {
		if len(fm.GithubRepo) > 0 {
			set.Add(fm.GithubRepo)
		}
	}

	return set, nil
}

type StarredRepository struct {
	Name        string
	Description string
	Login       string
	Url         string
	Language    string
	Topics      string
}

func (repo *StarredRepository) Write(vault *ObsidianVault, template *template.Template) error {
	date := time.Now().Format("20060102") + fmt.Sprintf("%04d", rand.Intn(10000))
	view := struct {
		Name        string
		Description string
		Login       string
		Url         string
		Language    string
		Topics      string
	}{
		Name:        repo.Name,
		Description: repo.Description,
		Login:       repo.Login,
		Url:         repo.Url,
		Language:    repo.Language,
		Topics:      repo.Topics,
	}

	buf := new(bytes.Buffer)
	if err := template.Execute(buf, view); err != nil {
		return err
	}

	reg, err := regexp.Compile("/+")
	if err != nil {
		return err
	}

	name := reg.ReplaceAllString(repo.Name, " ")
	fname := date + " - " + strings.Title(strings.ToLower(name)) + ".md"
	fpath := filepath.Join(vault.dpath, "github-stars", fname)

	err = os.MkdirAll(filepath.Join(vault.dpath, "github-stars"), 0700)
	if err != nil {
		return errors.Wrap(err, "failed creating github-stars")
	}

	err = os.WriteFile(fpath, buf.Bytes(), 0700)
	if err != nil {
		return errors.Wrap(err, "failed writing bookmark to file")
	}

	return nil
}

/*
 * Sync github stars into Obsidian
 *
 */
func SyncGithubStars(templatePath string, vault *ObsidianVault, conn *OpalDb) error {
	tmpl, err := LoadGithubStarTemplate(templatePath)
	if err != nil {
		return err
	}

	repoNames, err := GetPresentGithubStars(conn)
	if err != nil {
		return err
	}

	repos, err := conn.ListAbsentGithubStars(repoNames)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		err = repo.Write(vault, tmpl)

		if err != nil {
			return err
		}
	}

	// -- TODO
	return nil
}
