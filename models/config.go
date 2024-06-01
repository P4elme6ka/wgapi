package models

type Config struct {
	ListenAddr          string
	DefaultWgDevice     string
	LogPath             string
	Debug               bool
	StoreFile           string
	PersistentStore     bool
	InactivePeersDelete bool
}
