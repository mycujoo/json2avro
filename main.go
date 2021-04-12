package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	file := flag.String("file", "", "name of the file with JSON payload")
	name := flag.String("name", "", "name of the schema")
	nullable := flag.Bool("nullable", false, "Whether all the fields should be nullable")

	flag.Parse()

	if file == nil || *file == "" {
		log.Fatal(errors.New("--file is required"))
	}

	if name == nil || *name == "" {
		log.Fatal(errors.New("--name is required"))
	}

	if nullable == nil {
		log.Fatal(errors.New("nullable cannot be nil"))
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.ReadFile(filepath.Join(wd, *file))
	if err != nil {
		log.Fatal(err)
	}

	res, err := Parse(*name, f, *nullable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
}
