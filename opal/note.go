package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

/*
 * Obsidian note information
 */
type ObsidianNote struct {
	fpath string
	dpath string
	name  string
	date  int
}

/*
 * Construct an Obsidian note representation
 */
func NewObsidianNote(dpath string, fpath string) *ObsidianNote {
	baseName := strings.TrimPrefix(fpath, dpath+"/")
	parts := strings.SplitN(baseName, " - ", 2)

	if len(parts) != 2 {
		return &ObsidianNote{}
	}

	name := parts[1]
	date, err := strconv.Atoi(parts[0])

	if err != nil {
		return &ObsidianNote{}
	}

	return &ObsidianNote{
		fpath: fpath,
		dpath: dpath,
		name:  name,
		date:  date,
	}
}

/*
 * Merge updated frontmatter into existing frontmatter.
 *
 */
func (note *ObsidianNote) UpdateFrontmatter() {

}

/*
 * Write updated frontmatter into a document
 */
func (note *ObsidianNote) WriteFrontmatter() {

}

/*
 * Update a note's frontmatter with tags, title if the file
 * is unprocessed or the hash has changed since last processing
 */
func (note *ObsidianNote) FixFrontmatter(conn *OpalDb) error {
	changed, err := note.Changed(conn)

	if err != nil {
		return err
	}

	if !changed {
		return nil
	}

	// ensure title is present

	return nil
}

/*
 * Write document title to a markdown file, if the file is unprocessed or
 * the hash has changed since last processing
 */
func (note *ObsidianNote) WriteTitle(conn *OpalDb) error {
	return nil
}

/*
 * Write document title to a markdown file, if the file is unprocessed or
 * the hash has changed since last processing
 */
func (note *ObsidianNote) FixTitle(conn *OpalDb) error {
	changed, err := note.Changed(conn)

	if err != nil {
		return err
	}

	if !changed {
		return nil
	}

	// ensure title is present

	return nil
}

/*
 * Has the file changed since being processed?
 *
 */
func (note *ObsidianNote) Changed(conn *OpalDb) (bool, error) {
	fileHash, processedHash, err := conn.GetHashes(note.fpath)

	if err != nil {
		return false, err
	}

	return fileHash != processedHash, nil
}

func (note *ObsidianNote) Exists() (bool, error) {
	_, err := os.Stat(note.fpath)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}
