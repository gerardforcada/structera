package testing

type PointerFields struct {
    InEveryVersion *string `json:"in_every_version"`
    OnlyIn1        *int `json:"only_in_1"`
    From2ToEnd     *uint8 `json:"from_2_to_end"`
    FromStartTo3   *[]byte `json:"from_start_to_3"`
    From1to4       *float32 `json:"from_1_to_4"`
}
