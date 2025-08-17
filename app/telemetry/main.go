package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	duckDBClient, err := logger.NewDuckDBClient("f1_2021_packets.duckdb")
	if err != nil {
		panic(err)
	}

	listener, err := packet.NewRawListener(context.Background(), packet.DefaultNetwork, packet.DefaultAddress, duckDBClient)
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
