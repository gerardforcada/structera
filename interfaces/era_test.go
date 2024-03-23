package interfaces

import (
	"testing"
)

type mockEra struct {
	name    string
	version int
}

func (m mockEra) GetName() string {
	return m.name
}

func (m mockEra) GetVersion() int {
	return m.version
}

func TestMockEra(t *testing.T) {
	tests := []struct {
		name        string
		version     int
		wantName    string
		wantVersion int
	}{
		{"Test Era", 1, "Test Era", 1},
		{"Test Era 2", 2, "Test Era 2", 2},
	}

	for _, tt := range tests {
		mock := mockEra{name: tt.name, version: tt.version}
		if got := mock.GetName(); got != tt.wantName {
			t.Errorf("mockEra.GetName() = %v, want %v", got, tt.wantName)
		}
		if got := mock.GetVersion(); got != tt.wantVersion {
			t.Errorf("mockEra.GetVersion() = %v, want %v", got, tt.wantVersion)
		}
	}
}
