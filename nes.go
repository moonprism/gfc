package main

import "fmt"

type RAM []byte

func main() {
	err, c := NewCartridgeByPath("/home/aqua/Code/d.nes")
	fmt.Printf("%v, %v\n", err, c)
}
