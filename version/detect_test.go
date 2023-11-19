package version

import (
	"github.com/aws/smithy-go/ptr"
	"reflect"
	"testing"
)

const (
	Version1 Version = 1
	Version2 Version = 2
)

type V1[T Versioned[Version]] *T
type V2[T Versioned[Version]] *T

// Versions struct
type MockEntityVersions struct {
	V1 V1[MockEntityV1]
	V2 V2[MockEntityV2]
}

// Main struct
type MockEntity struct {
	PointerFields
	MockEntityVersions
}

// Initialize function
func (d *MockEntity) Initialize() {
	d.MockEntityVersions = MockEntityVersions{
		V1: &MockEntityV1{},
		V2: &MockEntityV2{},
	}
}

// Methods for the struct
func (d MockEntity) GetVersionStructs() []Versioned[Version] {
	return []Versioned[Version]{
		MockEntityV1{},
		MockEntityV2{},
	}
}

func (d MockEntity) GetBaseStruct() interface{} {
	return d.PointerFields
}

func (d MockEntity) DetectVersion() Version {
	return DetectBestMatch[Version, MockEntity](d)
}

func (d MockEntity) GetVersions() []Version {
	return []Version{
		Version1,
		Version2,
	}
}

func (d MockEntity) GetMinVersion() Version {
	return Version1
}

func (d MockEntity) GetMaxVersion() Version {
	return Version2
}

// Version-specific struct types and methods
type MockEntityV1 struct {
	InEveryVersion string `json:"in_every_version"`
	OnlyIn1        int    `json:"only_in_1"`
}

func (d MockEntityV1) GetVersion() Version {
	return Version1
}

type MockEntityV2 struct {
	InEveryVersion string `json:"in_every_version"`
	From2ToEnd     uint8  `json:"from_2_to_end"`
}

func (d MockEntityV2) GetVersion() Version {
	return Version2
}

type PointerFields struct {
	InEveryVersion *string `json:"in_every_version"`
	OnlyIn1        *int    `json:"only_in_1"`
	From2ToEnd     *uint8  `json:"from_2_to_end"`
}

func TestDetectBestMatch(t *testing.T) {
	type args[V any, T Entity[V]] struct {
		entity T
	}
	type testCase[V any, T Entity[V]] struct {
		name string
		args args[V, T]
		want V
	}
	tests := []testCase[Version, MockEntity]{
		{
			name: "DetectBestMatch V1",
			args: args[Version, MockEntity]{
				entity: MockEntity{
					PointerFields: PointerFields{
						InEveryVersion: ptr.String("InEveryVersion"),
						OnlyIn1:        ptr.Int(1),
						From2ToEnd:     nil,
					},
					MockEntityVersions: MockEntityVersions{
						V1: &MockEntityV1{},
						V2: &MockEntityV2{},
					},
				},
			},
			want: Version1,
		},
		{
			name: "DetectBestMatch V1",
			args: args[Version, MockEntity]{
				entity: MockEntity{
					PointerFields: PointerFields{
						InEveryVersion: ptr.String("InEveryVersion"),
						OnlyIn1:        nil,
						From2ToEnd:     ptr.Uint8(2),
					},
					MockEntityVersions: MockEntityVersions{
						V1: &MockEntityV1{},
						V2: &MockEntityV2{},
					},
				},
			},
			want: Version2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectBestMatch[Version, MockEntity](tt.args.entity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectBestMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
