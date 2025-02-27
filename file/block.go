package file

import (
	"fmt"
	"hash/fnv"
)

type Block struct {
	filename string
	blknum   int
}

func NewBlockId(filename string, blknum int) *Block {
	if filename == "" {
		panic("filename can not be empty")
	}
	return &Block{
		filename: filename,
		blknum:   blknum,
	}
}

func (b *Block) FileName() string {
	return b.filename
}

func (b *Block) Number() int {
	return b.blknum
}

func (b *Block) Equals(other *Block) bool {
	if other == nil {
		return false
	}

	return b.filename == other.filename && b.blknum == other.blknum
}

func (b *Block) String() string {
	return fmt.Sprintf("[file %s, block %d]", b.filename, b.blknum)
}

func (b *Block) Hash() int {
	h := fnv.New32()

	s := b.String()

	h.Write([]byte(s))

	return int(h.Sum32())
}
