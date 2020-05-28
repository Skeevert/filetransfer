package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"bufio"
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func handleConnection(conn net.Conn, Filepath string) {

	fmt.Println("Handling incoming connection...")
	f, err := os.Open(Filepath)
	if err != nil {
		fmt.Println("File cannot be opened")
		return
	}
	defer func() {
        if err = f.Close(); err != nil {
            log.Fatal(err)
        }
    }()
    r := bufio.NewReader(f)
    b := make([]byte, 1024)
    for {
        bytesNum, err := r.Read(b)
        if err != nil {
            break
        }
		fmt.Printf("Sending %d bytes...\n", bytesNum)
		_, err = conn.Write(b[:bytesNum])
		if err != nil {
			log.Fatal(err)
		}
	}
	conn.Close()
	fmt.Println("Transfer complete")
}

func main() {
	fmt.Println("initializing server part...")

	var IPAddress string
	var Filepath string
	fmt.Print("Enter IP address: ")
	fmt.Scan(&IPAddress)

	fmt.Print("Enter Filepath to transfer: ")
	fmt.Scan(&Filepath)

	if fileExists(Filepath) == false {
		fmt.Printf("Couldn't find filepath named %s", Filepath)
		return
	}

	fmt.Println("Starting listening...")
	ln, err := net.Listen("tcp", IPAddress)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, Filepath)
	}
}