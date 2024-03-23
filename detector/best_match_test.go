package detector

import (
	"fmt"
	"github.com/aws/smithy-go/ptr"
	"github.com/gerardforcada/structera/conversor"
	"github.com/gerardforcada/structera/interfaces"
	"reflect"
	"testing"
)

const (
	Version1 int = 1
	Version2 int = 2
)

// Versions struct
type MockEntityVersions struct {
	V1 MockEntityV1
	V2 MockEntityV2
}

// Main struct
type MockEntity struct {
	MockEntityAllFields
	MockEntityVersions
}

// Initialize function
func (d *MockEntity) Initialize() {
	d.MockEntityVersions = MockEntityVersions{
		V1: MockEntityV1{},
		V2: MockEntityV2{},
	}
}

// Methods for the struct
func (d MockEntity) GetVersionStructs() []interfaces.Era {
	return []interfaces.Era{
		MockEntityV1{},
		MockEntityV2{},
	}
}

func (d MockEntity) GetBaseStruct() any {
	return d.MockEntityAllFields
}

func (d MockEntity) DetectVersion() int {
	return BestMatchingEra[MockEntity](d)
}

func (d MockEntity) GetVersions() []int {
	return []int{
		Version1,
		Version2,
	}
}

func (d MockEntity) ToEra(target any) error {
	return conversor.ToEra(target, d)
}

func (d MockEntity) GetMinVersion() int {
	return Version1
}

func (d MockEntity) GetMaxVersion() int {
	return Version2
}

func (d MockEntity) GetEraFromVersion(version int) (interfaces.Era, error) {
	switch version {
	case Version1:
		return d.MockEntityVersions.V1, nil
	case Version2:
		return d.MockEntityVersions.V2, nil
	default:
		return nil, fmt.Errorf("unknown version %d", version)
	}
}

// Era-specific struct types and methods
type MockEntityV1 struct {
	InEveryVersion string `json:"in_every_version"`
	OnlyIn1        int    `json:"only_in_1"`
}

func (d MockEntityV1) GetVersion() int {
	return Version1
}

func (d MockEntityV1) GetName() string {
	return "MockEntity1"
}

func (d MockEntityV1) GetHub() interfaces.Hub {
	return MockEntity{}
}

type MockEntityV2 struct {
	InEveryVersion string `json:"in_every_version"`
	From2ToEnd     uint8  `json:"from_2_to_end"`
}

func (d MockEntityV2) GetVersion() int {
	return Version2
}

func (d MockEntityV2) GetName() string {
	return "MockEntity2"
}

func (d MockEntityV2) GetHub() interfaces.Hub {
	return MockEntity{}
}

type MockEntityAllFields struct {
	InEveryVersion *string `json:"in_every_version"`
	OnlyIn1        *int    `json:"only_in_1"`
	From2ToEnd     *uint8  `json:"from_2_to_end"`
}

func TestDetectBestMatch(t *testing.T) {
	type args[T interfaces.Hub] struct {
		entity T
	}
	type testCase[T interfaces.Hub] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[MockEntity]{
		{
			name: "BestMatchingEra V1",
			args: args[MockEntity]{
				entity: MockEntity{
					MockEntityAllFields: MockEntityAllFields{
						InEveryVersion: ptr.String("InEveryVersion"),
						OnlyIn1:        ptr.Int(1),
						From2ToEnd:     nil,
					},
					MockEntityVersions: MockEntityVersions{
						V1: MockEntityV1{},
						V2: MockEntityV2{},
					},
				},
			},
			want: Version1,
		},
		{
			name: "BestMatchingEra V1",
			args: args[MockEntity]{
				entity: MockEntity{
					MockEntityAllFields: MockEntityAllFields{
						InEveryVersion: ptr.String("InEveryVersion"),
						OnlyIn1:        nil,
						From2ToEnd:     ptr.Uint8(2),
					},
					MockEntityVersions: MockEntityVersions{
						V1: MockEntityV1{},
						V2: MockEntityV2{},
					},
				},
			},
			want: Version2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BestMatchingEra[MockEntity](tt.args.entity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BestMatchingEra() = %v, want %v", got, tt.want)
			}
		})
	}
}
