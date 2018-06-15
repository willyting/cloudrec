package iface

// FolderOperator ...
type FolderOperator interface {
	Readdirnames(n int) (names []string, err error)
}
