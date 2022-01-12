package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	_ "github.com/mattn/go-sqlite3"
	opal "github.com/rgrannell1/opal"
)

func main() {
	opts, err := docopt.ParseDoc(Usage())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	fpath, _ := opts.String("<fpath>")

	err = opal.Opal(&opal.OpalArgs{
		fpath,
		false,
		true,
	})

	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
