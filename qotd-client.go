package main

import (
	"fmt"
	"net"
	"os"
  "strings"

	"github.com/codegangsta/cli"
  "github.com/armon/mdns"
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
    serverAddress := ""
    if (len(c.Args()) == 0) {
      entriesCh := make(chan *mdns.ServiceEntry, 4)
      go func() string {
        for entry := range entriesCh {
          hostname := strings.Split(entry.Name, ".")[0] + ".local"
          addrs, _ := net.LookupIP(hostname)
          for addr := range addrs {
            ip := addrs[addr]
            if ip.To4() != nil { 
              str := fmt.Sprintf("%v:%v", ip, entry.Port)
              fmt.Println(str)
            }
          }
        }
        return "Test and go"
      }()

      mdns.Lookup("_qotd._tcp", entriesCh)
      close(entriesCh)
      os.Exit(0)
    } else if len(c.Args()) == 2 {
      serverAddress = c.Args()[0] + ":" + c.Args()[1]
		} else {
			fmt.Println("Usage <host> <port>")
			os.Exit(1)
    }

		udp := c.Bool("udp")
		tcp := c.Bool("tcp") || (!udp)

		var message = ""
		if tcp {
			message = connectOverTCP(serverAddress)
		} else {
			message = connectOverUDP(serverAddress)
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

func connectOverUDP(servAddr string) string {
	udpAddr, err := net.ResolveUDPAddr("udp", servAddr)
	if err != nil {
		println("Error Resolving UDP Address:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)

	buffer := make([]byte, 1500)
	conn.Write([]byte("\n"))
	conn.Read(buffer[0:])
	conn.Close()

	return string(buffer)
}
