package oof

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	f, err := os.OpenFile("/tmp/plugin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic("it broke")
	}

	logger = log.New(f, "", log.LstdFlags)
}
