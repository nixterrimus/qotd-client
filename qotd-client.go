package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "QOTD Client"
	app.Usage = "QOTD client"
	app.Flags = []cli.Flag{
		cli.BoolFlag{"tcp", "Access over TCP [default]"},
		cli.BoolFlag{"udp", "Access over UDP"},
	}

	app.Action = func(c *cli.Context) {
		if len(c.Args()) != 2 {
			fmt.Println("Client requires <host> <port>")
			os.Exit(1)
		}
		udp := c.Bool("udp")
		tcp := c.Bool("tcp") || (!udp)

		var message = ""
		if tcp {
			message = connectOverTCP(c.Args()[0] + ":" + c.Args()[1])
		} else {
			message = "UDP isn't supported, yet"
		}
		fmt.Println(message)
	}
	app.Run(os.Args)
}

func connectOverTCP(servAddr string) string {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)

	if err != nil {
		println("Server Read failed:", err.Error())
		os.Exit(1)
	}

	conn.Close()

	return string(reply)
}
