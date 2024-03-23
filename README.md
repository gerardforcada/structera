# Structera: Go Versioned Structs

Structera is a command-line tool developed in Go. It facilitates the automatic generation of versioned Go structs based on custom version tags, simplifying the management of different struct versions.

[![Go Reference](https://pkg.go.dev/badge/github.com/gerardforcada/structera.svg)](https://pkg.go.dev/github.com/gerardforcada/structera)
[![Go Tests](https://github.com/gerardforcada/structera/actions/workflows/test.yml/badge.svg)](https://github.com/gerardforcada/structera/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/gerardforcada/structera/badge.svg?branch=main)](https://coveralls.io/github/gerardforcada/structera?branch=main)

## Structera Summary

- [Installation](#installation)
- [Usage](#usage)
  - [Key Concepts](#key-concepts)
  - [Command-Line](#command-line)
- [How It Works](#how-it-works)
- [Version Tag](#version-tag)
- [Supporting Extra Tags](#supporting-extra-tags)
- [Code Usage](#code-usage)
- [Examples](#examples)
  - [Unmarshall into an Era Directly](#unmarshall-into-an-era-directly)
  - [Fill an Specific Era from the Hub](#fill-an-specific-era-from-the-hub)
  - [Use the Hub to Detect an Era Based on the Content](#use-the-hub-to-detect-an-era-based-on-the-content)
  - [Use Generics to Detect Hub Models](#use-generics-to-detect-hub-models)
  - [Fill Hub with Eras to Use Specific Fields](#fill-hub-with-eras-to-use-specific-fields)
- [Advanced Usage](#advanced-usage)
  - [Hub Details](#hub-details)
  - [Era Details](#era-details)
  - [Type Details](#type-details)
- [Contributing](#contributing)
- [License](#license)

----------------------------

## Installation

Ensure Go is installed on your system before installing Structera. Execute the following command to install:

```bash
go install github.com/gerardforcada/structera@latest
```

This command downloads and installs the Structera binary.

## Usage

### Key concepts

- **Hub**: The central place where all the versions of the struct are managed.
- **Era**: A specific version of the struct.

### Command-line

Structera is used through the command line with these arguments:

- `--file, -f`: Path to the Go file containing the struct.
- `--struct, -s`: Name of the struct for versioning.
- `--output, -o` (optional): Destination directory for the versioned struct files.
- `--force, -F` (optional): Overwrite the already existing eras
For example:

```bash
structera -f ./models/user.go -s User -o ./models/
```

This generates a folder with structs based on the `User` struct in `user.go`, placing them in the `./models/version` directory.

```bash
$ tree models/
models/
├── user.go # Original struct
└── version
    ├── user
    │   ├── v1.go # Era
    │   └── v2.go # Era
    ├── types.go
    └── user.go # Hub

2 directories, 5 files
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

## Code Usage 

After generation, use these structs directly in your code. Use the following summary to understand how to use the generated structs:

## Examples
- [Unmarshall into an Era directly](#unmarshall-into-an-era-directly)
- [Fill an specific Era from the Hub](#fill-an-specific-era-from-the-hub)
- [Use the Hub to detect an Era based on the content](#use-the-hub-to-detect-an-era-based-on-the-content)
- [Use Generics to detect Hub models](#use-generics-to-detect-hub-models)
- [Fill Hub with Eras to use specific fields](#fill-hub-with-eras-to-use-specific-fields)

## Unmarshall into an Era directly

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"main/models/version/user" // Import the user versioned model package
)

func main() {
	jsonString := `{"in_every_version":"hey"}`

	var era user.V1
	err := json.Unmarshal([]byte(jsonString), &era)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", era) // Prints {InEveryVersion:hey}
}
```

## Fill an specific Era from the Hub

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"main/models/version" // Import the version package to access the user hub
	"main/models/version/user" // Import the user versioned model package
)

func main() {
	jsonString := `{"in_every_version":"hey"}`

	var hub version.User
	err := json.Unmarshal([]byte(jsonString), &hub)
	if err != nil {
		panic(err)
	}

	var era user.V1
	err = hub.ToEra(&era)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", era) // Prints {InEveryVersion:hey}
}
```

## Use the Hub to detect an Era based on the content

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"main/models/version" // Import the version package to access the user hub
)

func main() {
	jsonString := `{"in_every_version":"hey"}`

	var hub version.User
	err := json.Unmarshal([]byte(jsonString), &hub)
	if err != nil {
		panic(err)
	}

	version := hub.DetectVersion() // Returns the lowest matching version where the content fits
	fmt.Printf("Detected version: %d\n", version) // Prints 1

	era, err := hub.GetEraFromVersion(version) // Returns the specific era based on the detected version
	if err != nil {
		panic(err)
	}

	err = hub.ToEra(&era) // Fill an era object with the generic hub content
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", era) // Prints {InEveryVersion:hey}
}
```

## Use Generics to detect Hub models

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"main/models/version" // Import the version package to access the user hub
)

func handleEra[hubType version.Type](input string) {
	if hubType == version.TypeUser {
		hub, err := version.GetHubFromType(hubType)
		if err != nil {
			return
		}

		err := json.Unmarshal([]byte(input), &hub)
		if err != nil {
			panic(err)
		}

		var era user.V1
		err = hub.ToEra(&era)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", era)
	}
}

func main() {
	jsonString := `{"in_every_version":"hey"}`
	handleEra[version.TypeUser](jsonString) // Prints {InEveryVersion:hey}
	handleEra[version.TypeAdmin](jsonString) // Does nothing
}
```

## Fill Hub with Eras to use specific fields

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"main/models/version" // Import the version package to access the user hub
)

func main() {
	jsonString := `{"in_every_version":"hey"}`

	var hub version.User
	err := json.Unmarshal([]byte(jsonString), &hub)
	if err != nil {
		panic(err)
	}

	version := hub.DetectVersion() // Returns the lowest matching version where the content fits
	fmt.Printf("Detected version: %d\n", version) // Prints 1

	era, err := hub.GetEraFromVersion(version) // Returns the specific era based on the detected version
	if err != nil {
		panic(err)
	}

	err = hub.ToEra(&era) // Fill an era object with the generic hub content
	if err != nil {
		panic(err)
	}

	err = hub.FillEra(era, version) // Fill the specific hub era with an era object content
	if err != nil {
		return
	}
	fmt.Println(hub.V1.InEveryVersion) // Prints "hey"
}
```

----------------------------

# Advanced usage

## Hub details

### Hub attributes
- `V<version>`: The specific era struct (e.g. V1, V2, V3, ...) => `hub.V1`
- `<OriginalField>`: The original field in the generic hub struct => `hub.InEveryVersion`
- `<Model>AllFields`: All the fields in the generic hub struct => `hub.UserAllFields`

### Hub methods
- `DetectVersion() int`: Returns the lowest matching version where the content fits => `hub.DetectVersion()`
- `GetEraFromVersion(version int) (interfaces.Era, error)`: Returns the specific era based on the detected version => `hub.GetEraFromVersion(1)`
- `ToEra(era any) error`: Fill an era object with the generic hub content => `hub.ToEra(&era)`
- `FillEra(era interfaces.Era, version int) error`: Fill the specific hub era with an era object content => `hub.FillEra(era, 1)`
- `GetVersions() []int`: Returns the list of versions available in the hub => `hub.GetVersions()`
- `GetMinVersion() int`: Returns the lowest version available in the hub => `hub.GetMinVersion()`
- `GetMaxVersion() int`: Returns the highest version available in the hub => `hub.GetMaxVersion()`

## Era details

### Era attributes
- `<OriginalField>`: The original field in the specific era struct => `era.InEveryVersion`

### Era methods
- `GetVersion() int`: Returns the version of the era => `era.GetVersion()`
- `GetName() string`: Returns the name of the era model => `era.GetName()`

## Type details

Structera generates a `types.go` file containing the `Type` enum, which is used to identify the hub model type. The `Type` enum is used to handle different hub models in a generic way.

```go
package version

type Type string

const (
    TypeAdmin Type = "admin"
    TypeUser Type = "user"
)
```

### Type methods
- `GetHubFromType(t Type) (interfaces.Hub, error)`: Returns the specific hub model based on the type => `version.GetHubFromType(version.TypeUser)`

----------------------------

## Contributing

Contributions to Structera are welcome! Please feel free to submit pull requests or create issues for bugs and feature requests.

## License

Structera is licensed under the GNU GPLv3 License.