package main

import (
	"time"

	"github.com/micro/go-os/log"
)

func main() {
	logger := log.NewLog(
		log.WithLevel(log.InfoLevel),
		log.WithFields(log.Fields{
			"logger": "os",
		}),
		log.WithOutput(
			log.NewOutput(log.OutputName("/dev/stdout")),
		),
	)

	for i := 0; i < 100; i++ {
		logger.Info("This is a log message")
		time.Sleep(time.Second)
	}
}
