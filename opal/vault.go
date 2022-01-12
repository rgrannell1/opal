package opal

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
		note := NewObsidianNote(vault.dpath, fpath)

		exists, err := note.Exists()
		if err != nil {
			return modified, err
		}

		if !exists {
			continue
		}

		changed, err := note.Changed(conn)
		if err != nil {
			return modified, err
		}

		if changed {
			modified = append(modified, note)
		}
	}

	return modified, nil
}

/*
 * Fix frontmatter for all markdown files in a vault
 *
 */
func (vault *ObsidianVault) FixFrontmatter(notes []*ObsidianNote, conn *OpalDb) error {
	for _, note := range notes {
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
func (vault *ObsidianVault) FixTitle(notes []*ObsidianNote, conn *OpalDb) error {
	for _, note := range notes {
		if err := note.FixTitle(conn); err != nil {
			return err
		}
	}

	return nil
}
