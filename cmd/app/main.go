package main

import (
	"github.com/kylerequez/go-sample-dashboard/src/servers"
)

func main() {
	if err := servers.Init(); err != nil {
		panic(err)
	}
}
