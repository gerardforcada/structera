package interfaces

import (
	"fmt"
	"reflect"
	"testing"
)

type mockHub struct {
	minVersion  int
	maxVersion  int
	detectedVer int
	versions    []int
	baseStruct  any
	eras        map[int]Era
}

func (m *mockHub) GetMinVersion() int {
	return m.minVersion
}

func (m *mockHub) GetMaxVersion() int {
	return m.maxVersion
}

func (m *mockHub) DetectVersion() int {
	return m.detectedVer
}

func (m *mockHub) GetEraFromVersion(version int) (Era, error) {
	era, exists := m.eras[version]
	if !exists {
		return nil, fmt.Errorf("version %d not found", version)
	}
	return era, nil
}

func (m *mockHub) GetVersions() []int {
	return m.versions
}

func (m *mockHub) GetVersionStructs() []Era {
	var versionStructs []Era
	for _, v := range m.versions {
		versionStructs = append(versionStructs, m.eras[v])
	}
	return versionStructs
}

func (m *mockHub) GetBaseStruct() any {
	return m.baseStruct
}

func (m *mockHub) ToEra(target any) error {
	return nil
}

func TestMockHub(t *testing.T) {
	mockEra1 := mockEra{name: "test1", version: 1}
	mockEra2 := mockEra{name: "test2", version: 2}
	mockEra3 := mockEra{name: "test3", version: 3}
	mockHub1 := mockHub{
		minVersion:  1,
		maxVersion:  3,
		detectedVer: 2,
		versions:    []int{1, 2, 3},
		baseStruct:  struct{}{},
		eras: map[int]Era{
			1: mockEra1,
			2: mockEra2,
			3: mockEra3,
		},
	}

	t.Run("TestGetMinVersion", func(t *testing.T) {
		if got := mockHub1.GetMinVersion(); got != mockHub1.minVersion {
			t.Errorf("GetMinVersion() = %v, want %v", got, mockHub1.minVersion)
		}
	})

	t.Run("TestGetMaxVersion", func(t *testing.T) {
		if got := mockHub1.GetMaxVersion(); got != mockHub1.maxVersion {
			t.Errorf("GetMaxVersion() = %v, want %v", got, mockHub1.maxVersion)
		}
	})

	t.Run("TestDetectVersion", func(t *testing.T) {
		if got := mockHub1.DetectVersion(); got != mockHub1.detectedVer {
			t.Errorf("DetectVersion() = %v, want %v", got, mockHub1.detectedVer)
		}
	})

	t.Run("TestGetEraFromVersion", func(t *testing.T) {
		for index, v := range mockHub1.versions {
			era, err := mockHub1.GetEraFromVersion(v)
			if err != nil {
				t.Errorf("GetEraFromVersion(%d) returned an error: %v", v, err)
			}
			if index == 0 && !reflect.DeepEqual(era, mockEra1) {
				t.Errorf("GetEraFromVersion(%d) = %v, want %v", v, era, mockEra1)
			}
			if index == 1 && !reflect.DeepEqual(era, mockEra2) {
				t.Errorf("GetEraFromVersion(%d) = %v, want %v", v, era, mockEra2)
			}
			if index == 2 && !reflect.DeepEqual(era, mockEra3) {
				t.Errorf("GetEraFromVersion(%d) = %v, want %v", v, era, mockEra3)
			}
		}
	})

	t.Run("TestGetVersions", func(t *testing.T) {
		got := mockHub1.GetVersions()
		if !reflect.DeepEqual(got, mockHub1.versions) {
			t.Errorf("GetVersions() = %v, want %v", got, mockHub1.versions)
		}
	})

	t.Run("TestGetVersionStructs", func(t *testing.T) {
		got := mockHub1.GetVersionStructs()
		want := []Era{mockEra1, mockEra2, mockEra3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetVersionStructs() = %v, want %v", got, want)
		}
	})

	t.Run("TestGetBaseStruct", func(t *testing.T) {
		if got := mockHub1.GetBaseStruct(); !reflect.DeepEqual(got, mockHub1.baseStruct) {
			t.Errorf("GetBaseStruct() = %v, want %v", got, mockHub1.baseStruct)
		}
	})

	t.Run("TestToEra", func(t *testing.T) {
		target := struct{}{}
		err := mockHub1.ToEra(&target)
		if err != nil {
			t.Errorf("ToEra() returned an error: %v", err)
		}
	})
}
