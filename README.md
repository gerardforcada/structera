# Structera: Go Versioned Structs

Structera is a command-line tool developed in Go. It facilitates the automatic generation of versioned Go structs based on custom version tags, simplifying the management of different struct versions.

[![Go Reference](https://pkg.go.dev/badge/github.com/gerardforcada/structera.svg)](https://pkg.go.dev/github.com/gerardforcada/structera)
[![Go Tests](https://github.com/gerardforcada/structera/actions/workflows/test.yml/badge.svg)](https://github.com/gerardforcada/structera/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/gerardforcada/structera/badge.svg?branch=main)](https://coveralls.io/github/gerardforcada/structera?branch=main)

## Installation

Ensure Go is installed on your system before installing Structera. Execute the following command to install:

```bash
go install github.com/gerardforcada/structera@latest
```

This command downloads and installs the Structera binary.

## Usage

Structera is used through the command line with these arguments:

- `--file, -f`: Path to the Go file containing the struct.
- `--struct, -s`: Name of the struct for versioning.
- `--output, -o` (optional): Destination directory for the versioned struct files.

For example:

```bash
structera -f ./models/user.go -s User -o ./models/
```

This generates a versioned folder with structs based on the `User` struct in `user.go`, placing them in the `./models/versioned` directory.

```bash
$ tree models/
models/
├── user.go
└── versioned
    ├── user
    │   ├── fields.go
    │   └── version.go
    └── user.go

2 directories, 4 files
```

For more details about the command-line options, run `structera --help`.

## How It Works

Structera processes a specified Go struct and creates different struct versions based on version tags in struct fields. Consider this struct:

```go
type User struct {
    InEveryVersion string
    OnlyIn1        string `version:"1"`
    From2ToEnd     string `version:"2+"`
    FromStartTo3   string `version:"-3"`
    From1to4       string `version:"1-4"`
    OnlyIn5        string `version:"5"`
}
```

Structera produces version-specific structs for each tag, enabling easy management of multiple versions.

```go
type UserV1 struct {
    InEveryVersion string
    OnlyIn1        string
    FromStartTo3   string
    From1to4       string
}

type UserV2 struct {
    InEveryVersion string
    From2ToEnd     string
    FromStartTo3   string
    From1to4       string
}

type UserV3 struct {
    InEveryVersion string
    From2ToEnd     string
    FromStartTo3   string
    From1to4       string
}

type UserV4 struct {
    InEveryVersion string
    From2ToEnd     string
    From1to4       string
}

type UserV5 struct {
    InEveryVersion string
    From2ToEnd     string
    OnlyIn5        string
}
```

## Version Tag

The version tag defines the struct version that includes a particular field. The tag formats are:

- `<no tag>`: The field will be included in all versions of the struct.
- `version:"1"`: The field will only be included in version 1 of the struct.
- `version:"2+"`: The field will be included in version 2 and all subsequent versions of the struct.
- `version:"-3"`: The field will be included in version 3 and all previous versions of the struct.
- `version:"1-4"`: The field will be included in versions 1 to 4 of the struct.

## Usage 

After generation, use these structs directly in your code. Example of using a generated struct:

```go
package main

import (
	"encoding/json"
	"fmt"
	
	models "main/models/versioned" // Import your versioned models package
)

func main() {
	// Create a new struct
	user := &models.UserV1{
		InEveryVersion: "hey",
	}
	fmt.Println(user.InEveryVersion) // Prints "hey"

	// ----------------------------------------------- //

	// Or unmarshal directly into the struct
	jsonString := `{"in_every_version":"hey"}`
	
	var user models.UserV1
	_ = json.Unmarshal([]byte(jsonString), &user)
	fmt.Println(user.InEveryVersion) // Prints "hey"
}
```

However, the generated structs can hold the versioned models. This is how you can use them:

```go
package main

import (
	"encoding/json"
	"fmt"
	
	models "main/models/versioned" // Import your versioned models package
)

func main() {
	// Create a new struct
	var user models.User
	user.V1 = &models.UserV1{
		InEveryVersion: "hey",
	}
	fmt.Println(user.V1.InEveryVersion) // Prints "hey"

	// ----------------------------------------------- //

	// Or unmarshal directly into the struct
	jsonString := `{"in_every_version":"hey"}`

	var user models.User
	user.Initialize() // Initializes all the versions

	err := json.Unmarshal([]byte(jsonString), user.V1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user.V1.InEveryVersion) // Prints "hey"
}
```

To handle unknown struct versions, unmarshall your content into the model and use the `DetectVersion` method:

```go
package main

import (
	"encoding/json"
	"fmt"
	
	models "main/models/versioned" // Import your versioned models package
)

func main() {
	jsonString := `{"in_every_version":"hey"}`

	var user models.User
	err := json.Unmarshal([]byte(jsonString), &user)
	if err != nil {
		panic(err)
	}

	version := user.DetectVersion()
	fmt.Println(version) // Prints "1" (user.Version1)
}
```

## Supporting Extra Tags

Structera can retain additional tags in generated structs, useful for preserving extra information like JSON tags.

```go
type User struct {
    InEveryVersion string `json:"in_every_version"`
    OnlyIn1        string `version:"1" json:"only_in_1"`
    From2ToEnd     string `version:"2+" json:"from_2_to_end"`
    FromStartTo3   string `version:"-3" json:"from_start_to_3"`
    From1to4       string `version:"1-4" json:"from_1_to_4"`
    OnlyIn5        string `version:"5" json:"only_in_5"`
}
```

Resulting struct with retained tags:

```go
type UserV1 struct {
    InEveryVersion string `json:"in_every_version"`
    OnlyIn1        string `json:"only_in_1"`
    FromStartTo3   string `json:"from_start_to_3"`
    From1to4       string `json:"from_1_to_4"`
}
```

## Contributing

Contributions to Structera are welcome! Please feel free to submit pull requests or create issues for bugs and feature requests.

## License

Structera is licensed under the GNU GPLv3 License.