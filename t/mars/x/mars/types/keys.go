package types

const (
	// ModuleName defines the module name
	ModuleName = "mars"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_mars"
)

var (
	ParamsKey = []byte("p_mars")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
