package version

import (
    "fmt"
    "github.com/gerardforcada/structera/detector"
    "github.com/gerardforcada/structera/interfaces"
    
)

type TestingAllFields struct {
    InEveryVersion *string `json:"in_every_version"`
    OnlyIn1        *int `json:"only_in_1"`
    From2ToEnd     *uint8 `json:"from_2_to_end"`
    FromStartTo3   *[]byte `json:"from_start_to_3"`
    From1to4       *float32 `json:"from_1_to_4"`
}

const (
    TestingVersion1 int = 1
    TestingVersion2 int = 2
    TestingVersion3 int = 3
    TestingVersion4 int = 4
)

// TestingVersions struct
type TestingVersions struct {
    V1 TestingV1
    V2 TestingV2
    V3 TestingV3
    V4 TestingV4
}

// Testing struct
type Testing struct {
    TestingAllFields
    TestingVersions
}

// GetVersionStructs method for the struct
func (hub Testing) GetVersionStructs() []interfaces.Era {
    return []interfaces.Era{
        TestingV1{},
        TestingV2{},
        TestingV3{},
        TestingV4{},
    }
}

func (hub Testing) GetEraFromVersion(version int) (interfaces.Era, error) {
    switch version {
    case TestingVersion1:
        return hub.TestingVersions.V1, nil
    case TestingVersion2:
        return hub.TestingVersions.V2, nil
    case TestingVersion3:
        return hub.TestingVersions.V3, nil
    case TestingVersion4:
        return hub.TestingVersions.V4, nil
    default:
        return nil, fmt.Errorf("unknown version %d", version)
    }
}

func (hub Testing) GetBaseStruct() interface{} {
    return hub.TestingAllFields
}

func (hub Testing) DetectVersion() int {
    return detector.BestMatchingEra[Testing](hub)
}

func (hub Testing) GetVersions() []int {
    return []int{
        TestingVersion1,
        TestingVersion2,
        TestingVersion3,
        TestingVersion4,
    }
}

func (hub Testing) GetMinVersion() int {
    return TestingVersion1
}

func (hub Testing) GetMaxVersion() int {
    return TestingVersion4
}
// TestingV1 Version-specific struct types and methods
type TestingV1 struct {
    InEveryVersion string `json:"in_every_version"`
    OnlyIn1        int `json:"only_in_1"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (era TestingV1) GetVersion() int {
    return TestingVersion1
}

func (era TestingV1) GetName() string {
    return "testing"
}

func (era TestingV1) GetHub() interfaces.Hub {
    return Testing{
        TestingAllFields: TestingAllFields{
            InEveryVersion: &era.InEveryVersion,
            OnlyIn1: &era.OnlyIn1,
            FromStartTo3: &era.FromStartTo3,
            From1to4: &era.From1to4,
        },
    }
}
// TestingV2 Version-specific struct types and methods
type TestingV2 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (era TestingV2) GetVersion() int {
    return TestingVersion2
}

func (era TestingV2) GetName() string {
    return "testing"
}

func (era TestingV2) GetHub() interfaces.Hub {
    return Testing{
        TestingAllFields: TestingAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            FromStartTo3: &era.FromStartTo3,
            From1to4: &era.From1to4,
        },
    }
}
// TestingV3 Version-specific struct types and methods
type TestingV3 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (era TestingV3) GetVersion() int {
    return TestingVersion3
}

func (era TestingV3) GetName() string {
    return "testing"
}

func (era TestingV3) GetHub() interfaces.Hub {
    return Testing{
        TestingAllFields: TestingAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            FromStartTo3: &era.FromStartTo3,
            From1to4: &era.From1to4,
        },
    }
}
// TestingV4 Version-specific struct types and methods
type TestingV4 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (era TestingV4) GetVersion() int {
    return TestingVersion4
}

func (era TestingV4) GetName() string {
    return "testing"
}

func (era TestingV4) GetHub() interfaces.Hub {
    return Testing{
        TestingAllFields: TestingAllFields{
            InEveryVersion: &era.InEveryVersion,
            From2ToEnd: &era.From2ToEnd,
            From1to4: &era.From1to4,
        },
    }
}