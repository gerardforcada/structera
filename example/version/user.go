package version

import (
    "fmt"
    "github.com/gerardforcada/structera/detector"
    "github.com/gerardforcada/structera/interfaces"
    
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
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *any
}

const (
    UserVersion1 int = 1
    UserVersion2 int = 2
    UserVersion3 int = 3
    UserVersion4 int = 4
    UserVersion5 int = 5
)

// UserVersions struct
type UserVersions struct {
    V1 UserV1
    V2 UserV2
    V3 UserV3
    V4 UserV4
    V5 UserV5
}

// User struct
type User struct {
    UserAllFields
    UserVersions
}

// GetVersionStructs method for the struct
func (hub User) GetVersionStructs() []interfaces.Era {
    return []interfaces.Era{
        UserV1{},
        UserV2{},
        UserV3{},
        UserV4{},
        UserV5{},
    }
}

func (hub User) GetEraFromVersion(version int) (interfaces.Era, error) {
    switch version {
    case UserVersion1:
        return hub.UserVersions.V1, nil
    case UserVersion2:
        return hub.UserVersions.V2, nil
    case UserVersion3:
        return hub.UserVersions.V3, nil
    case UserVersion4:
        return hub.UserVersions.V4, nil
    case UserVersion5:
        return hub.UserVersions.V5, nil
    default:
        return nil, fmt.Errorf("unknown version %d", version)
    }
}

func (hub User) GetBaseStruct() interface{} {
    return hub.UserAllFields
}

func (hub User) DetectVersion() int {
    return detector.BestMatchingEra[User](hub)
}

func (hub User) GetVersions() []int {
    return []int{
        UserVersion1,
        UserVersion2,
        UserVersion3,
        UserVersion4,
        UserVersion5,
    }
}

func (hub User) GetMinVersion() int {
    return UserVersion1
}

func (hub User) GetMaxVersion() int {
    return UserVersion5
}
// UserV1 Version-specific struct types and methods
type UserV1 struct {
    InEveryVersion    string `json:"in_every_version"`
    OnlyIn1           int `json:"only_in_1"`
    FromStartTo3      []byte `json:"from_start_to_3"`
    From1to4          float32 `json:"from_1_to_4"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       int
    AndDoublePointers int
    AndGenerics       any
    AndOldGenerics    any
}

func (era UserV1) GetVersion() int {
    return UserVersion1
}

func (era UserV1) GetName() string {
    return "user"
}

func (era UserV1) GetHub() interfaces.Hub {
    return User{
        UserAllFields: UserAllFields{
            InEveryVersion: &era.InEveryVersion,
            OnlyIn1: &era.OnlyIn1,
            FromStartTo3: &era.FromStartTo3,
            From1to4: &era.From1to4,
            WorksWithMaps: &era.WorksWithMaps,
            AndMapsInMaps: &era.AndMapsInMaps,
            AndSlices: &era.AndSlices,
            AndPointers: &era.AndPointers,
            AndDoublePointers: &era.AndDoublePointers,
            AndGenerics: &era.AndGenerics,
            AndOldGenerics: &era.AndOldGenerics,
        },
    }
}
// UserV2 Version-specific struct types and methods
type UserV2 struct {
    InEveryVersion    string `json:"in_every_version"`
    From2ToEnd        uint8 `json:"from_2_to_end"`
    FromStartTo3      []byte `json:"from_start_to_3"`
    From1to4          float32 `json:"from_1_to_4"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       int
    AndDoublePointers int
    AndGenerics       any
    AndOldGenerics    any
}

func (era UserV2) GetVersion() int {
    return UserVersion2
}

func (era UserV2) GetName() string {
    return "user"
}

func (era UserV2) GetHub() interfaces.Hub {
    return User{
        UserAllFields: UserAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            FromStartTo3: &era.FromStartTo3,
            From1to4: &era.From1to4,
            WorksWithMaps: &era.WorksWithMaps,
            AndMapsInMaps: &era.AndMapsInMaps,
            AndSlices: &era.AndSlices,
            AndPointers: &era.AndPointers,
            AndDoublePointers: &era.AndDoublePointers,
            AndGenerics: &era.AndGenerics,
            AndOldGenerics: &era.AndOldGenerics,
        },
    }
}
// UserV3 Version-specific struct types and methods
type UserV3 struct {
    InEveryVersion    string `json:"in_every_version"`
    From2ToEnd        uint8 `json:"from_2_to_end"`
    FromStartTo3      []byte `json:"from_start_to_3"`
    From1to4          float32 `json:"from_1_to_4"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       int
    AndDoublePointers int
    AndGenerics       any
    AndOldGenerics    any
}

func (era UserV3) GetVersion() int {
    return UserVersion3
}

func (era UserV3) GetName() string {
    return "user"
}

func (era UserV3) GetHub() interfaces.Hub {
    return User{
        UserAllFields: UserAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            FromStartTo3: &era.FromStartTo3,
            From1to4: &era.From1to4,
            WorksWithMaps: &era.WorksWithMaps,
            AndMapsInMaps: &era.AndMapsInMaps,
            AndSlices: &era.AndSlices,
            AndPointers: &era.AndPointers,
            AndDoublePointers: &era.AndDoublePointers,
            AndGenerics: &era.AndGenerics,
            AndOldGenerics: &era.AndOldGenerics,
        },
    }
}
// UserV4 Version-specific struct types and methods
type UserV4 struct {
    InEveryVersion    string `json:"in_every_version"`
    From2ToEnd        uint8 `json:"from_2_to_end"`
    From1to4          float32 `json:"from_1_to_4"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       int
    AndDoublePointers int
    AndGenerics       any
    AndOldGenerics    any
}

func (era UserV4) GetVersion() int {
    return UserVersion4
}

func (era UserV4) GetName() string {
    return "user"
}

func (era UserV4) GetHub() interfaces.Hub {
    return User{
        UserAllFields: UserAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            From1to4: &era.From1to4,
            WorksWithMaps: &era.WorksWithMaps,
            AndMapsInMaps: &era.AndMapsInMaps,
            AndSlices: &era.AndSlices,
            AndPointers: &era.AndPointers,
            AndDoublePointers: &era.AndDoublePointers,
            AndGenerics: &era.AndGenerics,
            AndOldGenerics: &era.AndOldGenerics,
        },
    }
}
// UserV5 Version-specific struct types and methods
type UserV5 struct {
    InEveryVersion    string `json:"in_every_version"`
    From2ToEnd        uint8 `json:"from_2_to_end"`
    OnlyIn5           rune `json:"only_in_5"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       int
    AndDoublePointers int
    AndGenerics       any
    AndOldGenerics    any
}

func (era UserV5) GetVersion() int {
    return UserVersion5
}

func (era UserV5) GetName() string {
    return "user"
}

func (era UserV5) GetHub() interfaces.Hub {
    return User{
        UserAllFields: UserAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            OnlyIn5: &era.OnlyIn5,
            WorksWithMaps: &era.WorksWithMaps,
            AndMapsInMaps: &era.AndMapsInMaps,
            AndSlices: &era.AndSlices,
            AndPointers: &era.AndPointers,
            AndDoublePointers: &era.AndDoublePointers,
            AndGenerics: &era.AndGenerics,
            AndOldGenerics: &era.AndOldGenerics,
        },
    }
}