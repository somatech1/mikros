package service

// ServerPort is a framework type representing a port listened by a specific
// type of service.
type ServerPort int32

// Int32 returns the current ServerPort in int32 format.
func (s ServerPort) Int32() int32 {
	return int32(s)
}

// Int returns the current ServerPort in int format.
func (s ServerPort) Int() int {
	return int(s)
}
