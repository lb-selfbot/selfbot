package authsync

// SyncType represents the type of synchronization operation.
type SyncType int

const (
	SyncTypeConfigUpdate SyncType = iota
	SyncTypeConsoleAppend
	SyncTypeConsoleClear
	SyncTypeConsoleSet
)
