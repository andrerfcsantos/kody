package config

type FlagConfig[T any] struct {
	Key           string
	FlagName      string
	FlagShortHand string
	Default       T
	Description   string
}
