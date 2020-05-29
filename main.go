package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"flag"
)

func main() {
	addrPtr := flag.String("addr", "", "IP Address")
	pathPtr := flag.String("file", "", "Path to file")
	serverMode := flag.Bool("server", false, "Server mode")
	transferMode := flag.Bool("send", false, "Mode")
	flag.Parse()

	if *addrPtr == "" {
		fmt.Println("Usage: -addr=IPAddr -file=Filepath -server=true/false -send=true/false")
		return
	}

	fmt.Println("initializing...")
	if *transferMode && !fileExists(*pathPtr) {
		fmt.Printf("Couldn't find filepath named %s\n", *pathPtr)
		return
	}
	if *serverMode == true {
		initServer(*addrPtr, *pathPtr, *transferMode)
	} else {
		initClient(*addrPtr, *pathPtr, *transferMode)
	}
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func initServer(addr string, filepath string, transferMode bool) {
	fmt.Println("Starting listening...")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if transferMode == true {
			go handleSend(conn, filepath)
		} else {
			go handleReceive(conn)
		}
	}
}

func initClient(addr string, filepath string, transferMode bool) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	if transferMode == true {
		handleSend(conn, filepath)
	} else {
		handleReceive(conn)
	}
}
