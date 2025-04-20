package main

var (
	flagsFunc = flags
	runFunc   = run
)

func main() {
	flagsFunc()
	runFunc()
}
