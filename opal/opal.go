package main

import (
	"os"
	"path/filepath"
)

type OpalArgs struct {
	fpath string
	audit bool
	fix   bool
}

/*
 * Main application; audit or fix Obsidian notes
 */
func Opal(args *OpalArgs) error {
	vault := ObsidianVault{args.fpath}
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	dbpath := filepath.Join(home, ".diatom.sqlite")

	conn, err := NewOpalDb(dbpath)
	if err != nil {
		return err
	}

	err = conn.CreateTables()
	if err != nil {
		return err
	}

	if err := vault.FixFrontmatter(conn); err != nil {
		return err
	}

	if err := vault.FixTitle(conn); err != nil {
		return err
	}

	return nil
}
