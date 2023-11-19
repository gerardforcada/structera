package versioned

import (
    "github.com/gerardforcada/structera/version"
    "github.com/gerardforcada/structera/example/versioned/testing"
)

// Versions struct
type TestingVersions struct {
    V1 testing.V1[TestingV1]
    V2 testing.V2[TestingV2]
    V3 testing.V3[TestingV3]
    V4 testing.V4[TestingV4]
}

// Main struct
type Testing struct {
    testing.PointerFields
    TestingVersions
}

// Initialize function
func (d *Testing) Initialize() {
    d.TestingVersions = TestingVersions{
        V1: &TestingV1{},
        V2: &TestingV2{},
        V3: &TestingV3{},
        V4: &TestingV4{},
    }
}

// Methods for the struct
func (d Testing) GetVersionStructs() []version.Versioned[testing.Version] {
    return []version.Versioned[testing.Version]{
        TestingV1{},
        TestingV2{},
        TestingV3{},
        TestingV4{},
    }
}

func (d Testing) GetBaseStruct() interface{} {
    return d.PointerFields
}

func (d Testing) DetectVersion() testing.Version {
    return version.DetectBestMatch[testing.Version, Testing](d)
}

func (d Testing) GetVersions() []testing.Version {
    return []testing.Version{
        testing.Version1,
        testing.Version2,
        testing.Version3,
        testing.Version4,
    }
}

func (d Testing) GetMinVersion() testing.Version {
    return testing.Version1
}

func (d Testing) GetMaxVersion() testing.Version {
    return testing.Version4
}

// Version-specific struct types and methods
type TestingV1 struct {
    InEveryVersion string `json:"in_every_version"`
    OnlyIn1        int `json:"only_in_1"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (d TestingV1) GetVersion() testing.Version {
    return testing.Version1
}
type TestingV2 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (d TestingV2) GetVersion() testing.Version {
    return testing.Version2
}
type TestingV3 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (d TestingV3) GetVersion() testing.Version {
    return testing.Version3
}
type TestingV4 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (d TestingV4) GetVersion() testing.Version {
    return testing.Version4
}