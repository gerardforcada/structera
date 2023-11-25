package version

import (
    "fmt"
    "github.com/gerardforcada/structera/interfaces"
)

type Type string

const (
    TypeTesting Type = "testing"
    TypeUser Type = "user"
)

func GetHubFromType(t Type) (interfaces.Hub, error) {
    switch t {
    case TypeTesting:
        return Testing{}, nil
    case TypeUser:
        return User{}, nil
    }
    return nil, fmt.Errorf("unknown type %s", t)
}
