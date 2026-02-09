package main

import (
	"flag"
	"log"
	"triple-s/internal/cmd"
	"triple-s/internal/server"
)

func main() {
	port := flag.Int("port", 8080, "")
	dir := flag.String("dir", "data", "")
	help := flag.Bool("help", false, "")
	flag.Parse()
	if err := cmd.ArgsRoute(*port, *dir, *help); err != nil {
		log.Fatalf("Братья, у нас проблемы тут с %v", err)
		return
	}

	log.Fatal(server.StartServer(*port))
}
