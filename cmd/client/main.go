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

		for {
			response, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Server Disconnected", err)
				return
			}
			response = strings.TrimSpace(response)

			if response == "__END_RESPONSE__" {
				break
			}
			fmt.Println(response)

		}

		if strings.ToUpper(text) == "EXIT" {
			return
		}

	}

}
