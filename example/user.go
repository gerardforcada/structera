package example

// User Original struct with version tags
type User struct {
	InEveryVersion string  `json:"in_every_version"`
	OnlyIn1        int     `version:"1" json:"only_in_1"`
	From2ToEnd     uint8   `version:"2+" json:"from_2_to_end"`
	FromStartTo3   []byte  `version:"-3" json:"from_start_to_3"`
	From1to4       float32 `version:"1-4" json:"from_1_to_4"`
	OnlyIn5        rune    `version:"5" json:"only_in_5"`

	WorksWithMaps     map[string]int64
	AndMapsInMaps     map[string]map[string]int64
	AndSlices         []int
	AndPointers       *int
	AndDoublePointers **int
	AndGenerics       any
	AndOldGenerics    interface{}
}
