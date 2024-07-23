package main

import (
	"os"

	"learninggd/analyser"
	"learninggd/reader"
	"learninggd/students"
	"learninggd/writer"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		os.Exit(1)
	}
	app := args[0]

	switch app {
	case "reader":
		reader.Reader()
	case "writer":
		writer.Writer()
	case "students":
		students.StudentApi()
	case "analyser":
		analyser.Analyser()
	default:
		os.Exit(1)
	}
}
