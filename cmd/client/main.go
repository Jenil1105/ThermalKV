package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial(
		"tcp",
		"localhost:8080",
	)

	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connected to ThermalKV")

	serverReader := bufio.NewReader(conn)
	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("ThermalKV> ")

		text, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Println("Input error: ", err)
			return
		}
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		conn.Write([]byte(text + "\n"))
		response, err := serverReader.ReadString('\n')

		if err != nil {
			fmt.Println("Server disconnected")
			return
		}

		fmt.Print(response)

		if strings.ToUpper(text) == "EXIT" {
			return
		}

	}

}
