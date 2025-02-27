package file

import "encoding/binary"

type Page struct {
	buffer []byte
}

func NewPage(blockSize int) *Page {
	return &Page{buffer: make([]byte, blockSize)}
}

func NewPageFromBytes(bytes []byte) *Page {
	return &Page{buffer: bytes}
}

func (p *Page) GetInt(offset int) int {
	value := binary.BigEndian.Uint32(p.buffer[offset : offset+4])
	return int(value)
}

func (p *Page) SetInt(offset int, n int) {
	if offset+4 > len(p.buffer) {
		panic("offset is out of bounds")
	}
	binary.BigEndian.PutUint32(p.buffer[offset:offset+4], uint32(n))
}

func (p *Page) GetBytes(offset int) []byte {
	byteLength := p.GetInt(offset)
	return p.buffer[offset+4 : offset+4+byteLength]
}

func (p *Page) SetBytes(offset int, bytes []byte) {
	if offset+4+len(bytes) > len(p.buffer) {
		panic("offset is out of bounds")
	}
	byteLength := len(bytes)
	p.SetInt(offset, byteLength)
	copy(p.buffer[offset+4:offset+4+byteLength], bytes)
}

func (p *Page) GetString(offset int) string {
	bytes := p.GetBytes(offset)
	return string(bytes)
}

func (p *Page) SetString(offset int, s string) {
	bytes := []byte(s)
	p.SetBytes(offset, bytes)
}

func (p *Page) MaxLength(strLen int) int {
	return 4 + strLen*4
}

func (p *Page) Buffer() []byte {
	return p.buffer
}
