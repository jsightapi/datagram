package jse

type ParseError struct {
	Message string
	Line    uint32
	Index   uint32
}
