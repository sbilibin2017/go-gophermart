package main

func main() {
	config := flags()
	err := run(config)
	exit(err)
}
