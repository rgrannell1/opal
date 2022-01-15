package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	_ "github.com/mattn/go-sqlite3"
	opal "github.com/rgrannell1/opal/pkg"
)

func main() {
	opts, err := docopt.ParseDoc(opal.Usage())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	fpath, _ := opts.String("<fpath>")

	err = opal.Opal(&opal.OpalArgs{
		Fpath: fpath,
		Audit: false,
		Fix:   true,
	})

	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
