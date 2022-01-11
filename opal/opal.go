package main

import (
	"os"
	"path/filepath"

	diatom "github.com/rgrannell1/diatom/pkg"
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

	diatomArgs := diatom.DiatomArgs{
		Dir:    args.fpath,
		DBPath: dbpath,
	}
	err = diatom.Diatom(&diatomArgs)
	if err != nil {
		return err
	}

	conn, err := NewOpalDb(dbpath)
	if err != nil {
		return err
	}

	err = conn.CreateTables()
	if err != nil {
		return err
	}

	// list modified files and modify them
	notes, err := vault.ListModifiedMarkdown(conn)
	if err := vault.FixFrontmatter(notes, conn); err != nil {
		return err
	}

	if err := vault.FixTitle(notes, conn); err != nil {
		return err
	}

	err = conn.MarkComplete(notes)
	if err != nil {
		return err
	}

	// generate bookmark files using coppermind and diatom data
	err = SyncBookmarks("./template.txt", &vault, conn)
	if err != nil {
		return err
	}

	return nil
}
