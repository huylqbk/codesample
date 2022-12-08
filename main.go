package main

import "github.com/huylqbk/codesample/logger"

func main() {
	log := logger.NewLogger().LogFile("./log").SetCaller().SetLevel(5)
	log.Info("main", "key", "value")
}
