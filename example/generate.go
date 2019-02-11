//+build ignore

package main

import (
	"log"
	"os"

	"github.com/Evertras/go-ts-type-generator/example"
	"github.com/Evertras/go-ts-type-generator/typegen"
)

func main() {
	// Overwrite our old file if it exists, we don't want it
	outfile, err := os.Create("types.ts")

	if err != nil {
		log.Fatal(err)
	}

	// Be good citizens and close our stuff
	defer outfile.Close()

	// Write a helpful header at the top of our file
	outfile.WriteString("/* THIS FILE IS GENERATED, DO NOT EDIT */\n\n")

	generator := typegen.New()

	// Write all our types to our file
	if err = generator.GenerateTypes(
		outfile,
		example.SomeData{},
		example.Outer{}, // note that Inner will be included due to recursion
	); err != nil {
		log.Fatal(err)
	}

	// End the file with a newline because we're good people
	outfile.WriteString("\n")
}
