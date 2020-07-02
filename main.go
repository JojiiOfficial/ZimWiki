package main

import (
	"fmt"

	zim "github.com/tim-st/go-zim"
)

func main() {
	file, err := zim.Open("archlinux_wiki.zim")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(file.License())
}
