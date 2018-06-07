package main

import (
	"gachamachine/gacha"
	"gachamachine/machine"
	"gachamachine/storage"
	"os"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		os.Exit(1)
	}
	server := machine.NewServer()
	gacha.SetStorage(&storage.LocalStroage{})
	server.AddHandlers(gacha.GetHandlers())
	server.Run(port)
}
