package user

type PointerFields struct {
    InEveryVersion    *string `json:"in_every_version"`
    OnlyIn1           *int `json:"only_in_1"`
    From2ToEnd        *uint8 `json:"from_2_to_end"`
    FromStartTo3      *[]byte `json:"from_start_to_3"`
    From1to4          *float32 `json:"from_1_to_4"`
    OnlyIn5           *rune `json:"only_in_5"`
    WorksWithMaps     *map[string]int64
    AndMapsInMaps     *map[string]map[string]int64
    AndSlices         *[]int
    AndPointers       *int
    AndDoublePointers *int
    AndGenerics       *any
    AndOldGenerics    *interface{}
}
