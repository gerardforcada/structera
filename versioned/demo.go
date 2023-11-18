package versioned

import (
    "github.com/gerardforcada/structera/version"
    "github.com/gerardforcada/structera/versioned/demo"
)

// Versions struct
type DemoVersions struct {
    V1 demo.V1[DemoV1]
    V2 demo.V2[DemoV2]
    V3 demo.V3[DemoV3]
    V4 demo.V4[DemoV4]
    V5 demo.V5[DemoV5]
}

// Main struct
type Demo struct {
    demo.PointerFields
    DemoVersions
}

// Initialize function
func (d *Demo) Initialize() {
    d.DemoVersions = DemoVersions{
        V1: &DemoV1{},
        V2: &DemoV2{},
        V3: &DemoV3{},
        V4: &DemoV4{},
        V5: &DemoV5{},
    }
}

// Methods for the struct
func (d Demo) GetVersionStructs() []version.Versioned[demo.Version] {
    return []version.Versioned[demo.Version]{
        DemoV1{},
        DemoV2{},
        DemoV3{},
        DemoV4{},
        DemoV5{},
    }
}

func (d Demo) GetBaseStruct() interface{} {
    return d.PointerFields
}

func (d Demo) DetectVersion() demo.Version {
    return version.DetectBestMatch[demo.Version, Demo](d)
}

func (d Demo) GetVersions() []demo.Version {
    return []demo.Version{
        demo.Version1,
        demo.Version2,
        demo.Version3,
        demo.Version4,
        demo.Version5,
    }
}

func (d Demo) GetMinVersion() demo.Version {
    return demo.Version1
}

func (d Demo) GetMaxVersion() demo.Version {
    return demo.Version5
}

// Version-specific struct types and methods
type DemoV1 struct {
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

func (d DemoV1) GetVersion() demo.Version {
    return demo.Version1
}
type DemoV2 struct {
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

func (d DemoV2) GetVersion() demo.Version {
    return demo.Version2
}
type DemoV3 struct {
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

func (d DemoV3) GetVersion() demo.Version {
    return demo.Version3
}
type DemoV4 struct {
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

func (d DemoV4) GetVersion() demo.Version {
    return demo.Version4
}
type DemoV5 struct {
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

func (d DemoV5) GetVersion() demo.Version {
    return demo.Version5
}