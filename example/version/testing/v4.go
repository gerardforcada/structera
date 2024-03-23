package testing

// V4 Version-specific struct types and methods
type V4 struct {
    InEveryVersion string `json:"in_every_version"`
    From2ToEnd     uint8 `json:"from_2_to_end"`
    From1to4       float32 `json:"from_1_to_4"`
}

func (era V4) GetVersion() int {
    return 4
}

func (era V4) GetName() string {
    return "testing"
}
