package demo

import "github.com/gerardforcada/structera/version"

type Version version.Version

const (
	Version1 Version = 1
	Version2 Version = 2
	Version3 Version = 3
	Version4 Version = 4
	Version5 Version = 5
)

type Versioned[T any] version.Versioned[T]

type V1[T Versioned[Version]] *T
type V2[T Versioned[Version]] *T
type V3[T Versioned[Version]] *T
type V4[T Versioned[Version]] *T
type V5[T Versioned[Version]] *T
