package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Module string

const (
	ModuleVersion Module = "0.2.3-alpha"
	ModulePackage Module = "github.com/gerardforcada/structera"
)

func main() {
	cli()
}

func cli() {
	var (
		fileName    string
		structName  string
		outputDir   string
		showVersion bool
		showHelp    bool
	)

	// Define both long and short flag versions
	flag.StringVar(&fileName, "file", "", "Path to the Go file containing the struct")
	flag.StringVar(&fileName, "f", "", "Path to the Go file containing the struct (shorthand)")

	flag.StringVar(&structName, "struct", "", "Name of the struct")
	flag.StringVar(&structName, "s", "", "Name of the struct (shorthand)")

	flag.StringVar(&outputDir, "output", "", "Output directory (optional)")
	flag.StringVar(&outputDir, "o", "", "Output directory (optional) (shorthand)")

	flag.BoolVar(&showHelp, "help", false, "Print the help page and exit")
	flag.BoolVar(&showHelp, "h", false, "Print the help page and exit (shorthand)")

	flag.BoolVar(&showVersion, "version", false, "Print the version of Structera and exit")
	flag.BoolVar(&showVersion, "v", false, "Print the version of Structera and exit (shorthand)")

	flag.Parse()

	if showVersion {
		fmt.Printf("Structera version %s\n", ModuleVersion)
		os.Exit(0)
	}

	// Check if the required flags are set
	if fileName == "" || structName == "" || showHelp {
		fmt.Printf("Structera version %s\n\n", ModuleVersion)
		fmt.Println("Structera is a command-line tool for versioning Go structs.")
		fmt.Printf("For more information, updates, or contributions, visit https://%s\n\n", ModulePackage)
		fmt.Println("Usage:")
		fmt.Println("  structera -f <path-to-struct-file> -s <StructName> [-o <output-directory>]")
		fmt.Println("  structera --file <path-to-struct-file> --struct <StructName> [--output <output-directory>]")
		fmt.Println("\nOptions:")
		fmt.Println("  --file,    -f  Path to the Go file containing the struct")
		fmt.Println("  --struct,  -s  Name of the struct to version")
		fmt.Println("  --output,  -o  (Optional) Output directory for the versioned struct files")
		fmt.Println("  --help,    -h  Prints this page and exit")
		fmt.Println("  --version, -v  Print the version of Structera and exit")
		fmt.Println("\nExample:")
		fmt.Println("  structera -f ./models/user.go -s User")
		fmt.Println("  structera -f ./models/user.go -s User -o ./models/versioned")
		fmt.Println("  structera --file ./models/user.go --struct User")
		fmt.Println("  structera --file ./models/user.go --struct User --output ./models/versioned\n")

		if showHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if outputDir == "" {
		outputDir = filepath.Dir(fileName)
	}

	generator := Generator{
		Version:  &Version{},
		Resolver: &Resolver{},
		Filename: fileName,
		StructName: StructName{
			Original: structName,
			Lower:    strings.ToLower(structName),
		},
		OutputDir: outputDir,
	}

	if err := generator.VersionedStructs(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Versioned structs generated successfully.")
}
