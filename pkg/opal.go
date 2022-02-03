package opal

import (
	"os"
	"path/filepath"

	copper "github.com/rgrannell1/coppermind/pkg"
	diatom "github.com/rgrannell1/diatom/pkg"
)

type OpalArgs struct {
	Fpath string
	Audit bool
	Fix   bool
}

/*
 * Main application; audit or fix Obsidian notes
 */
func Opal(args *OpalArgs) error {
	vault := ObsidianVault{args.Fpath}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbpath := filepath.Join(home, ".diatom.sqlite")
	err = diatom.Diatom(&diatom.DiatomArgs{
		Dir:    args.Fpath,
		DBPath: dbpath,
	})

	if err != nil {
		return err
	}

	err = copper.Coppermind()
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

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	root := filepath.Dir(ex)

	// generate bookmark files using coppermind and diatom data
	err = SyncBookmarks(filepath.Join(root, "./pinboard-template.txt"), &vault, conn)
	if err != nil {
		return err
	}

	//err = SyncGithubStars(filepath.Join(root, "./github-template.txt"), &vault, conn)
	//if err != nil {
	//	return err
	//}

	return nil
}
