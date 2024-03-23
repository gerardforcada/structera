package main

import (
	"flag"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"strings"
)

type Module string

const (
	ModuleFolder  Module = "version"
	ModuleVersion Module = "0.4.0-beta"
	ModulePackage Module = "github.com/gerardforcada/structera"
)

func main() {
	if err := cli(flag.CommandLine); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}

func cli(flagset *flag.FlagSet) error {
	var (
		fileName    string
		structName  string
		outputDir   string
		showVersion bool
		showHelp    bool
		force       bool
	)

	// Define both long and short flag versions
	flagset.StringVar(&fileName, "file", "", "Path to the Go file containing the struct")
	flagset.StringVar(&fileName, "f", "", "Path to the Go file containing the struct (shorthand)")

	flagset.StringVar(&structName, "struct", "", "Name of the struct")
	flagset.StringVar(&structName, "s", "", "Name of the struct (shorthand)")

	flagset.StringVar(&outputDir, "output", "", "Output directory (optional)")
	flagset.StringVar(&outputDir, "o", "", "Output directory (optional) (shorthand)")

	flagset.BoolVar(&showHelp, "help", false, "Print the help page and exit")
	flagset.BoolVar(&showHelp, "h", false, "Print the help page and exit (shorthand)")

	flagset.BoolVar(&showVersion, "version", false, "Print the version of Structera and exit")
	flagset.BoolVar(&showVersion, "v", false, "Print the version of Structera and exit (shorthand)")

	flagset.BoolVar(&force, "force", false, "Replace existing versioned struct files")
	flagset.BoolVar(&force, "F", false, "Replace existing versioned struct files (shorthand)")

	err := flagset.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if showVersion {
		fmt.Printf("Structera version %s\n", ModuleVersion)
		return nil
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
		fmt.Println("  --force,   -F  Replace existing versioned struct files")
		fmt.Println("  --struct,  -s  Name of the struct to version")
		fmt.Println("  --output,  -o  (Optional) Output directory for the versioned struct files")
		fmt.Println("  --help,    -h  Prints this page and exit")
		fmt.Println("  --version, -v  Print the version of Structera and exit")
		fmt.Println("\nExample:")
		fmt.Println("  structera -f ./models/user.go -s User")
		fmt.Println("  structera -f ./models/user.go -s User -o ./models/versioned")
		fmt.Println("  structera --file ./models/user.go --struct User")
		fmt.Println("  structera --file ./models/user.go --struct User --output ./models/versioned")
		fmt.Println()

		if showHelp {
			return nil
		}
		return fmt.Errorf("missing required flags")
	}

	if outputDir == "" {
		outputDir = filepath.Dir(fileName)
	}

	generator := Generator{
		Format:   &Format{},
		Resolver: &Resolver{},
		Filename: fileName,
		StructName: StructName{
			Original: structName,
			Lower:    strings.ToLower(structName),
			Snake:    strcase.SnakeCase(structName),
		},
		OutputDir: outputDir,
		Package:   string(ModuleFolder),
		Replace:   force,
	}

	if err := generator.VersionedStructs(); err != nil {
		return err
	}

	fmt.Println("Versioned structs generated successfully.")
	return nil
}
