package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	opts, err := docopt.ParseDoc(Usage())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	fpath, _ := opts.String("<fpath>")

	args := &OpalArgs{
		fpath,
		false,
		true,
	}
	err = Opal(args)

	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
