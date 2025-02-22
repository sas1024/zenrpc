package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"time"

	"github.com/vmkteam/zenrpc/v2/parser"
)

const (
	openIssueURL = "https://github.com/vmkteam/zenrpc/issues/new"
	githubURL    = "https://github.com/vmkteam/zenrpc"
)

func main() {
	start := time.Now()
	fmt.Printf("Generator version: %s\n", parser.GeneratorVersion)

	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[len(os.Args)-1]
	} else {
		filename = os.Getenv("GOFILE")
	}

	if filename == "" {
		fmt.Fprintln(os.Stderr, "File path is empty")
		os.Exit(1)
	}

	fmt.Printf("Entrypoint: %s\n", filename)

	// create package info
	pi, err := parser.NewPackageInfo(filename)
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	outputFileName := pi.OutputFilename()

	// remove output file if it already exists
	if _, err := os.Stat(outputFileName); err == nil {
		if err := os.Remove(outputFileName); err != nil {
			printError(err)
			os.Exit(1)
		}
	}

	if err := pi.Parse(filename); err != nil {
		printError(err)
		os.Exit(1)
	}

	if len(pi.Services) == 0 {
		fmt.Fprintln(os.Stderr, "Services not found")
		os.Exit(1)
	}

	if err := generateFile(outputFileName, pi); err != nil {
		printError(err)
		os.Exit(1)
	}

	fmt.Printf("Generated: %s\n", outputFileName)
	fmt.Printf("Duration: %dms\n", int64(time.Since(start)/time.Millisecond))
	fmt.Println()
	fmt.Print(pi)
	fmt.Println()
}

func printError(err error) {
	// print error to stderr
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)

	// print contact information to stdout
	fmt.Println("\nYou may help us and create issue:")
	fmt.Printf("\t%s\n", openIssueURL)
	fmt.Println("For more information, see:")
	fmt.Printf("\t%s\n\n", githubURL)
}

func generateFile(outputFileName string, pi *parser.PackageInfo) error {
	file, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	output := new(bytes.Buffer)
	if err := serviceTemplate.Execute(output, pi); err != nil {
		return err
	}

	source, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	if _, err = file.Write(source); err != nil {
		return err
	}

	return nil
}
