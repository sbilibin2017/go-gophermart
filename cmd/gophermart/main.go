package main

import "os"

func main() {
	config := flags()
	err := run(config)
	code := exit(err)
	os.Exit(code)
}
