package interfaces

type Era interface {
	GetName() string
	GetVersion() int
	GetHub() Hub
}
