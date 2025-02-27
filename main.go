package main

import (
	"fmt"
	"simpledb/file"
)

func main() {
	block := file.NewBlockId("users.tbl", 0)

	fmt.Println(block)
}
