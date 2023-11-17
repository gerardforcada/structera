package version

type Version int

type Versioned[T any] interface {
	GetVersion() T
}

type Entity[T any] interface {
	GetMinVersion() T
	GetMaxVersion() T
	DetectVersion() T
	GetVersions() []T
	GetVersionStructs() []Versioned[T]
	GetBaseStruct() interface{}
}
