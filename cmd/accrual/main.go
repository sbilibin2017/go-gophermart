package main

import (
	"os"

	"github.com/sbilibin2017/go-gophermart/cmd/accrual/app"
)

var exitFunc = os.Exit

func main() {
	exitFunc(app.Run())
}
