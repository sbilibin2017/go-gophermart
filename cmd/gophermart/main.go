package main

import "os"

var exitFunc = os.Exit

func main() {
	exitFunc(0)
}
