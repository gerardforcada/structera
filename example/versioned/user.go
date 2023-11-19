package versioned

import (
    "github.com/gerardforcada/structera/version"
    "github.com/gerardforcada/structera/example/versioned/user"
)

// Versions struct
type UserVersions struct {
    V1 user.V1[UserV1]
    V2 user.V2[UserV2]
    V3 user.V3[UserV3]
    V4 user.V4[UserV4]
    V5 user.V5[UserV5]
}

// Main struct
type User struct {
    user.PointerFields
    UserVersions
}

// Initialize function
func (d *User) Initialize() {
    d.UserVersions = UserVersions{
        V1: &UserV1{},
        V2: &UserV2{},
        V3: &UserV3{},
        V4: &UserV4{},
        V5: &UserV5{},
    }
}

// Methods for the struct
func (d User) GetVersionStructs() []version.Versioned[user.Version] {
    return []version.Versioned[user.Version]{
        UserV1{},
        UserV2{},
        UserV3{},
        UserV4{},
        UserV5{},
    }
}

func (d User) GetBaseStruct() interface{} {
    return d.PointerFields
}

func (d User) DetectVersion() user.Version {
    return version.DetectBestMatch[user.Version, User](d)
}

func (d User) GetVersions() []user.Version {
    return []user.Version{
        user.Version1,
        user.Version2,
        user.Version3,
        user.Version4,
        user.Version5,
    }
}

func (d User) GetMinVersion() user.Version {
    return user.Version1
}

func (d User) GetMaxVersion() user.Version {
    return user.Version5
}

// Version-specific struct types and methods
type UserV1 struct {
    InEveryVersion    *string `json:"in_every_version"`
    OnlyIn1           *int `json:"only_in_1"`
    FromStartTo3      *[]byte `json:"from_start_to_3"`
    From1to4          *float32 `json:"from_1_to_4"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *interface{}
}

func (d UserV1) GetVersion() user.Version {
    return user.Version1
}
type UserV2 struct {
    InEveryVersion    *string `json:"in_every_version"`
    From2ToEnd        *uint8 `json:"from_2_to_end"`
    FromStartTo3      *[]byte `json:"from_start_to_3"`
    From1to4          *float32 `json:"from_1_to_4"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *interface{}
}

func (d UserV2) GetVersion() user.Version {
    return user.Version2
}
type UserV3 struct {
    InEveryVersion    *string `json:"in_every_version"`
    From2ToEnd        *uint8 `json:"from_2_to_end"`
    FromStartTo3      *[]byte `json:"from_start_to_3"`
    From1to4          *float32 `json:"from_1_to_4"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *interface{}
}

func (d UserV3) GetVersion() user.Version {
    return user.Version3
}
type UserV4 struct {
    InEveryVersion    *string `json:"in_every_version"`
    From2ToEnd        *uint8 `json:"from_2_to_end"`
    From1to4          *float32 `json:"from_1_to_4"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *interface{}
}

func (d UserV4) GetVersion() user.Version {
    return user.Version4
}
type UserV5 struct {
    InEveryVersion    *string `json:"in_every_version"`
    From2ToEnd        *uint8 `json:"from_2_to_end"`
    OnlyIn5           *rune `json:"only_in_5"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *interface{}
}

func (d UserV5) GetVersion() user.Version {
    return user.Version5
}