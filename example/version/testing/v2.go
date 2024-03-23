package testing

// V2 Version-specific struct types and methods
type V2 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    FromStartTo3   []byte `json:"from_start_to_3"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (era V2) GetVersion() int {
    return 2
}

func (era V2) GetName() string {
    return "testing"
}
