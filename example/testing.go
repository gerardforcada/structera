package example

// Testing Original struct with version tags
type Testing struct {
	InEveryVersion string  `json:"in_every_version"`
	OnlyIn1        int     `version:"1" json:"only_in_1"`
	From2ToEnd     uint8   `version:"2+" json:"from_2_to_end"`
	FromStartTo3   []byte  `version:"-3" json:"from_start_to_3"`
	From1to4       float32 `version:"1-4" json:"from_1_to_4"`
}
