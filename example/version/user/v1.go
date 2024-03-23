package user

// V1 Version-specific struct types and methods
type V1 struct {
    InEveryVersion    string `json:"in_every_version"`
    OnlyIn1           int `json:"only_in_1"`
    FromStartTo3      []byte `json:"from_start_to_3"`
    From1to4          float32 `json:"from_1_to_4"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       *int
    AndDoublePointers **int
    AndGenerics       any
    AndOldGenerics    any
}

func (era V1) GetVersion() int {
    return 1
}

func (era V1) GetName() string {
    return "user"
}
