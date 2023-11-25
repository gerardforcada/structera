package interfaces

type Hub interface {
	GetMinVersion() int
	GetMaxVersion() int
	DetectVersion() int
	GetEraFromVersion(int) (Era, error)
	GetVersions() []int
	GetVersionStructs() []Era
	GetBaseStruct() interface{}
}
