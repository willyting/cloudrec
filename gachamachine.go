package main

import (
	"gachamachine/gacha"
	"gachamachine/machine"
	"os"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		os.Exit(1)
	}
	server := machine.NewServer()
	server.AddHandlers(gacha.GetHandlers())
	server.Run(port)
}
