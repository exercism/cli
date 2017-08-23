package main

import (
	"fmt"
	"os"
)

func main() {
	path := "C:\\Users\\robertph\\code\\gopath\\src\\github.com\\exercism\\cli\\fixtures\\detect-path-type\\symlinked-dir-windows"
	info, err := os.Lstat(path)
	fmt.Println(info)
	fmt.Println(err)
	nonExistent := "C:\\Hello"
	info, err = os.Lstat(nonExistent)
	fmt.Println(info)
	fmt.Println(err)
}
