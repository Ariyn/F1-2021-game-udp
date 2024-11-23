package main

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	godotenv.Load()
}

func main() {
	//stdoutLogger := logger.NewStdoutClient()
	sqlClient, err := logger.NewSqlClient(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	//stdoutLogger

	listener, err := packet.NewListener(context.Background(), packet.DefaultNetwork, packet.DefaultAddress, sqlClient)
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
