package cache

type unsignedLong = uint64
type cString struct {
	size    unsignedLong
	content []byte
}

type cacheHeader struct {
	deleted   uint64
	schema    uint64
	showing   uint64
	className cString
}
