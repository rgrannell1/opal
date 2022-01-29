package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
	opal "github.com/rgrannell1/opal/pkg"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
