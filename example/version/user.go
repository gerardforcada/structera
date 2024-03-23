package version

import (
    "fmt"
    "github.com/gerardforcada/structera/conversor"
    "github.com/gerardforcada/structera/detector"
    "github.com/gerardforcada/structera/interfaces"
    "github.com/gerardforcada/structera/example/version/user"
)

type UserAllFields struct {
    InEveryVersion    *string `json:"in_every_version"`
    OnlyIn1           *int `json:"only_in_1"`
    From2ToEnd        *uint8 `json:"from_2_to_end"`
    FromStartTo3      *[]byte `json:"from_start_to_3"`
    From1to4          *float32 `json:"from_1_to_4"`
    OnlyIn5           *rune `json:"only_in_5"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       **int
    AndDoublePointers ***int
    AndGenerics       *any
    AndOldGenerics    *any
}

// UserVersions struct
type UserVersions struct {
    V1 user.V1
    V2 user.V2
    V3 user.V3
    V4 user.V4
    V5 user.V5
}

// User struct
type User struct {
    UserAllFields
    UserVersions
}

// GetVersionStructs method for the struct
func (hub User) GetVersionStructs() []interfaces.Era {
    return []interfaces.Era{
        user.V1{},
        user.V2{},
        user.V3{},
        user.V4{},
        user.V5{},
    }
}

func (hub User) GetEraFromVersion(version int) (interfaces.Era, error) {
    switch version {
    case user.V1{}.GetVersion():
        return hub.UserVersions.V1, nil
    case user.V2{}.GetVersion():
        return hub.UserVersions.V2, nil
    case user.V3{}.GetVersion():
        return hub.UserVersions.V3, nil
    case user.V4{}.GetVersion():
        return hub.UserVersions.V4, nil
    case user.V5{}.GetVersion():
        return hub.UserVersions.V5, nil
    default:
        return nil, fmt.Errorf("unknown version %d", version)
    }
}

func (hub User) ToEra(target any) error {
    return conversor.ToEra(target, hub)
}

func (hub User) GetBaseStruct() any {
    return hub.UserAllFields
}

func (hub User) DetectVersion() int {
    return detector.BestMatchingEra[User](hub)
}

func (hub User) GetVersions() []int {
    return []int{
        user.V1{}.GetVersion(),
        user.V2{}.GetVersion(),
        user.V3{}.GetVersion(),
        user.V4{}.GetVersion(),
        user.V5{}.GetVersion(),
    }
}

func (hub User) GetMinVersion() int {
    return user.V1{}.GetVersion()
}

func (hub User) GetMaxVersion() int {
    return user.V5{}.GetVersion()
}
