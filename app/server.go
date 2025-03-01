package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	dir        = flag.String("dir", "/tmp/redis-files", "Redis file directory ")
	dbfilename = flag.String("dbfilename", "dump.rdb", "RDB filename")
)

func main() {
	flag.Parse()
	cfg = &Config{
		Dir:        "/tmp/redis-files",
		DbFilename: "dump.rdb",
	}
	cfg.Dir = *dir
	cfg.DbFilename = *dbfilename
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	log.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	go ReduceKeyTTL()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Failed to make connection :", err.Error())
		}
		go handleConns(conn)
	}
}

func handleConns(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	_, err := reader.Read(buffer)
	if err != nil {
		log.Println("ERROR READING FROM CONNECTION")
	}
	resp := Parse(&conn, buffer)
	conn.Write([]byte(resp))
}
