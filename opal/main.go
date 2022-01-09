package main

import (
	"github.com/docopt/docopt-go"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	opts, err := docopt.ParseDoc(Usage())
	if err != nil {
		panic(err)
	}

	fpath, _ := opts.String("<fpath>")

	args := &OpalArgs{
		fpath,
		false,
		true,
	}
	err = Opal(args)

	if err != nil {
		panic(err)
	}
}
