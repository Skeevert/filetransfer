package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"io"
)

func handleConnection(conn net.Conn, f *os.File) {

	fmt.Println("Starting receiving...")
	defer func() {
        if err := f.Close(); err != nil {
            log.Fatal(err)
        }
    }()
    b := make([]byte, 1024)
    for {
        bytesNum, err := conn.Read(b)
		if err != nil {
            if err != io.EOF {
                fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println("Got", bytesNum, "bytes")
		if bytesNum == 0 {
			break
		}
		f.Write(b)
	}
	fmt.Println("Transfer complete")
}

func main() {
	fmt.Println("initializing client part...")

	var IPAddress string
	var Filepath string
	fmt.Print("Enter IP address: ")
	fmt.Scan(&IPAddress)

	fmt.Print("Enter Filepath to save incoming data: ")
	fmt.Scan(&Filepath)
	f, err := os.Create(Filepath)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", IPAddress)
	if err != nil {
		log.Fatal(err)
	}
	handleConnection(conn, f)
}