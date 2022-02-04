package opal

import (
	"errors"
	"fmt"
)

/*
 * Assert no duplicate files (files with the same hash) are present. This can happen
 * if Opal is broken & outputs duplicate files.
 *
 */
func AssertNoDuplicates(conn *OpalDb) error {
	hashes, err := conn.ListHashes()
	if err != nil {
		return err
	}

	failed := false

	ctr := NewCounter()
	for _, pair := range hashes {
		// -- don't complain about empty files
		if pair[1] != "0" {
			ctr.Add(pair[1], pair[0])
		}
	}

	dupes := ctr.Duplicates()

	if len(dupes) > 0 {
		failed = true
		for _, dupe := range ctr.Duplicates() {
			fmt.Println(dupe)
		}
	}

	if failed {
		return errors.New("files with duplicate hashes found; this indicates files were accidentally duplicated")
	} else {
		return nil
	}
}

/*
 * Validate the Obsidian repository.
 *
 */
func (vault *ObsidianVault) Validate(conn *OpalDb) error {
	if err := AssertNoDuplicates(conn); err != nil {
		return err
	}

	return nil
}
