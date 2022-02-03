package opal

import (
	"errors"
	"fmt"
)

/*
 * Assert no duplicate files (files with the same hash) are present. This can happen
 * if Opal is broken
 *
 */
func AssertNoDuplicates(conn *OpalDb) error {
	hashes, err := conn.ListHashes()
	if err != nil {
		return err
	}

	ctr := NewCounter()
	for _, pair := range hashes {
		ctr.Add(pair[1], pair[0])
	}

	dupes := ctr.Duplicates()

	if len(dupes) > 0 {
		for _, dupe := range ctr.Duplicates() {
			fmt.Println(dupe)
		}
	}

	return errors.New("files with duplicate hashes found; this indicates files were accidentally duplicated")
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
