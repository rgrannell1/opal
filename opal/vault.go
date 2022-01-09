package main

import (
	"path/filepath"
)

type ObsidianVault struct {
	dpath string
}

/*
 * List markdown files in the vault
 *
 */
func (vault *ObsidianVault) ListMarkdown() ([]string, error) {
	return filepath.Glob(vault.dpath + "/*.md")
}

/*
 * List markdown files modified since last processing
 */
func (vault *ObsidianVault) ListModifiedMarkdown(conn *OpalDb) ([]*ObsidianNote, error) {
	modified := make([]*ObsidianNote, 0)
	fpaths, err := vault.ListMarkdown()

	if err != nil {
		return modified, err
	}

	for _, fpath := range fpaths {
		note, err := NewObsidianNote(vault.dpath, fpath)

		if err != nil {
			return modified, err
		}

		changed, err := note.Changed(conn)

		if changed == true {
			modified = append(modified, note)
		}
	}

	return modified, nil
}

/*
 * Fix frontmatter for all markdown files in a vault
 *
 */
func (vault *ObsidianVault) FixFrontmatter(conn *OpalDb) error {
	notes, err := vault.ListModifiedMarkdown(conn)
	if err != nil {
		return err
	}

	for _, note := range notes {
		if err != nil {
			return err
		}

		if err := note.FixFrontmatter(conn); err != nil {
			return err
		}
	}

	return nil
}

/*
 * Fix titles for all markdown files
 *
 */
func (vault *ObsidianVault) FixTitle(conn *OpalDb) error {
	md, err := vault.ListModifiedMarkdown(conn)
	if err != nil {
		return err
	}

	for _, note := range md {
		if err != nil {
			return err
		}

		if err := note.FixTitle(conn); err != nil {
			return err
		}
	}

	return nil
}
