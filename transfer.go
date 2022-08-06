package main

import (
	"fmt"
	"os"
	"net"
	"log"
	"encoding/binary"
	"path/filepath"
	"bufio"
	"io"
)

func handleSend(conn net.Conn, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("File cannot be opened")
		return
	}
	defer func() {
        if err = f.Close(); err != nil {
            log.Fatal(err)
        }
    }()
	if !sendPath(conn, filepath) {
		fmt.Println("Transfer Error")
		return
	}
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

func handleReceive(conn net.Conn) {
	fmt.Println("Starting receiving...")
	filename := getFileName(conn)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
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
		f.Write(b[:bytesNum])
	}
	fmt.Println("Transfer complete")
}

func sendPath(conn net.Conn, Filepath string) bool {
	data_content := []byte(Filepath)
	var length uint32 = uint32(len(data_content))
	data_len := make([]byte, 4)
	binary.LittleEndian.PutUint32(data_len, length)

	_, err := conn.Write([]byte(data_len))
	if err != nil {
		return false
	}
	_, err = conn.Write(data_content)
	if err != nil {
		return false
	}
	return true
}

func getFileName(conn net.Conn) string {
	data_len := make([]byte, 4)
	var length uint32

	_, err := conn.Read(data_len) // Read uint32 from sender, which determines filepath len
	if err != nil {
		log.Fatal(err)
	}
	length = uint32(binary.LittleEndian.Uint32(data_len)) // Transform received data to uint32
	data_content := make([]byte, length)
	_, err = conn.Read(data_content) // Receive filepath from sender
	if err != nil {
		log.Fatal(err)
	}
	path := string(data_content[:length])
	filename := filepath.Base(path)
	return filename
}
