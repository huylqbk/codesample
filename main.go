package main

import "github.com/huylqbk/codesample/logger"

func main() {
	log := logger.NewLogger("").LogFile().SetCaller().SetLevel(5)
	log.Info("main", "abc", "test", "111")
}
