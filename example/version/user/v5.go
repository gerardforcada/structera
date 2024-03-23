package user

// V5 Version-specific struct types and methods
type V5 struct {
    InEveryVersion    string `json:"in_every_version"`
    From2ToEnd        uint8 `json:"from_2_to_end"`
    OnlyIn5           rune `json:"only_in_5"`
    WorksWithMaps     map[string]int64
    AndMapsInMaps     map[string]map[string]int64
    AndSlices         []int
    AndPointers       int
    AndDoublePointers int
    AndGenerics       any
    AndOldGenerics    any
}

func (era V5) GetVersion() int {
    return 5
}

func (era V5) GetName() string {
    return "user"
}
