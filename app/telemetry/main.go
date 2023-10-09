package main

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
)

func main() {
	stdoutLogger := logger.NewStdoutClient()
	redisLogger := logger.NewRedisClient(context.Background(), "localhost:6379", 0)

	listener, err := packet.NewListener(context.Background(), packet.DefaultNetwork, packet.DefaultAddress, stdoutLogger, redisLogger)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	fmt.Println("monitor start")
	defer fmt.Println("monitor ended")

	if err := listener.Run(); err != nil {
		log.Fatal(err)
	}
}
