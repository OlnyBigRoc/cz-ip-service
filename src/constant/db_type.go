package constant

type DbType int

const (
	Invalid DbType = -1
	IPV4    DbType = iota
	IPV6
)
