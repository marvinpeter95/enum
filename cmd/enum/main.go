package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/marvinpeter95/enum/cmd/enum/generator"
)

func main() {
	// Define command-line flags for the generator options.
	var (
		flagTypes = flag.String("types", "",
			"Comma-separated list of enum type names")
		flagCaseInsensitive = flag.Bool("case-insensitive", false,
			"Perform case-insensitive matching when paring enum values")
		flagNoTextMarshaling = flag.Bool("no-text-marshaling", false,
			"Do not generate MarshalText and UnmarshalText methods for the enums")
		flagNoParser = flag.Bool("no-parser", false,
			"Do not generate Parse[Enum] functions for the enums")
	)

	flag.Parse()

	var filename string

	// Determine the input filename from command-line arguments or environment variable.
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	} else if value, ok := os.LookupEnv("GOFILE"); ok && value != "" {
		filename = value
	} else {
		fmt.Fprintln(os.Stderr, "Error: no input file specified")
		os.Exit(1)
	}

	// Ensure that types are specified for code generation.
	if flagTypes == nil || *flagTypes == "" {
		fmt.Fprintln(os.Stderr, "Error: no enum types specified")
		os.Exit(1)
	}

	types := strings.Split(*flagTypes, ",")
	for i := range types {
		types[i] = strings.TrimSpace(types[i])
	}

	// Process the input file and generate code for the specified enum types with the given options.
	code, err := generator.Process(filename, types, generator.Options{
		CaseInsensitive:    *flagCaseInsensitive,
		SkipTextMarshaling: *flagNoTextMarshaling,
		SkipParser:         *flagNoParser,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	// Write the generated code to an output file with a suffix "_enum.go".
	outputFilename := strings.TrimSuffix(filename, ".go") + "_enum.go"
	if err := os.WriteFile(outputFilename, code, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to write generated code to file: ", err)
		os.Exit(1)
	}
}
