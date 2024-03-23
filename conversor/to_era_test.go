package conversor

import (
	"fmt"
	"github.com/gerardforcada/structera/interfaces"
	"testing"
)

type mockHub struct {
	minVersion  int
	maxVersion  int
	detectedVer int
	versions    []int
	baseStruct  any
	eras        map[int]interfaces.Era
}

func (m mockHub) GetMinVersion() int {
	return m.minVersion
}

func (m mockHub) GetMaxVersion() int {
	return m.maxVersion
}

func (m mockHub) DetectVersion() int {
	return m.detectedVer
}

func (m mockHub) GetEraFromVersion(version int) (interfaces.Era, error) {
	era, exists := m.eras[version]
	if !exists {
		return era, fmt.Errorf("version %d not found", version)
	}
	return era, nil
}

func (m mockHub) GetVersions() []int {
	return m.versions
}

func (m mockHub) GetVersionStructs() []interfaces.Era {
	var versionStructs []interfaces.Era
	for _, v := range m.versions {
		versionStructs = append(versionStructs, m.eras[v])
	}
	return versionStructs
}

func (m mockHub) GetBaseStruct() any {
	return m.baseStruct
}

func (m mockHub) ToEra(target any) error {
	return nil
}

type mockEra struct {
	Name    string
	Version int
}

func (m mockEra) GetName() string {
	return m.Name
}

func (m mockEra) GetVersion() int {
	return m.Version
}

func TestToEra(t *testing.T) {
	mockEra1 := mockEra{Name: "test1", Version: 1}
	mockHub1 := mockHub{
		minVersion:  1,
		maxVersion:  1,
		detectedVer: 1,
		versions:    []int{1},
		baseStruct:  struct{}{},
		eras: map[int]interfaces.Era{
			1: mockEra1,
		},
	}

	t.Run("ValidConversion", func(t *testing.T) {
		var era interfaces.Era = mockEra{}
		err := ToEra(&era, mockHub1)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("NonPointerTarget", func(t *testing.T) {
		var era mockEra
		err := ToEra(era, mockHub1)
		if err == nil {
			t.Error("Expected error for non-pointer target, got none")
		}
	})

	t.Run("NilPointerTarget", func(t *testing.T) {
		var era *mockEra = nil
		err := ToEra(era, mockHub1)
		if err == nil {
			t.Error("Expected error for nil pointer target, got none")
		}
	})
}
