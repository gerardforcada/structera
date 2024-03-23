package version

import (
	"encoding/json"
	"fmt"
	"github.com/gerardforcada/structera/conversor"
	"github.com/gerardforcada/structera/detector"
	"github.com/gerardforcada/structera/example/version/testing"
	"github.com/gerardforcada/structera/interfaces"
)

type TestingAllFields struct {
	InEveryVersion *string  `json:"in_every_version"`
	OnlyIn1        *int     `json:"only_in_1"`
	From2ToEnd     *uint8   `json:"from_2_to_end"`
	FromStartTo3   *[]byte  `json:"from_start_to_3"`
	From1to4       *float32 `json:"from_1_to_4"`
}

// TestingVersions struct
type TestingVersions struct {
	V1 testing.V1
	V2 testing.V2
	V3 testing.V3
	V4 testing.V4
}

// Testing struct
type Testing struct {
	TestingAllFields
	TestingVersions
}

// GetVersionStructs method for the struct
func (hub Testing) GetVersionStructs() []interfaces.Era {
	return []interfaces.Era{
		testing.V1{},
		testing.V2{},
		testing.V3{},
		testing.V4{},
	}
}

func (hub Testing) GetEraFromVersion(version int) (interfaces.Era, error) {
	switch version {
	case testing.V1{}.GetVersion():
		return hub.TestingVersions.V1, nil
	case testing.V2{}.GetVersion():
		return hub.TestingVersions.V2, nil
	case testing.V3{}.GetVersion():
		return hub.TestingVersions.V3, nil
	case testing.V4{}.GetVersion():
		return hub.TestingVersions.V4, nil
	default:
		return nil, fmt.Errorf("unknown version %d", version)
	}
}

func (hub Testing) ToEra(target any) error {
	return conversor.ToEra(target, hub)
}

func (hub Testing) GetBaseStruct() any {
	return hub.TestingAllFields
}

func (hub Testing) DetectVersion() int {
	return detector.BestMatchingEra[Testing](hub)
}

func (hub Testing) GetVersions() []int {
	return []int{
		testing.V1{}.GetVersion(),
		testing.V2{}.GetVersion(),
		testing.V3{}.GetVersion(),
		testing.V4{}.GetVersion(),
	}
}

func (hub Testing) GetMinVersion() int {
	return testing.V1{}.GetVersion()
}

func (hub Testing) GetMaxVersion() int {
	return testing.V4{}.GetVersion()
}

func (hub *Testing) FillEra(era interfaces.Era, version int) error {
	eraJSON, err := json.Marshal(era)
	if err != nil {
		return fmt.Errorf("error marshalling era: %w", err)
	}

	switch version {
	case testing.V1{}.GetVersion():
		err = json.Unmarshal(eraJSON, &hub.TestingVersions.V1)
	case testing.V2{}.GetVersion():
		err = json.Unmarshal(eraJSON, &hub.TestingVersions.V2)
	case testing.V3{}.GetVersion():
		err = json.Unmarshal(eraJSON, &hub.TestingVersions.V3)
	case testing.V4{}.GetVersion():
		err = json.Unmarshal(eraJSON, &hub.TestingVersions.V4)
	default:
		return fmt.Errorf("unknown version %d", version)
	}

	return err
}
