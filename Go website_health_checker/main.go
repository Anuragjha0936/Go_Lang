package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "health-checker",
		Usage: "A tiny tool that checks whether a website is running or down",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "domain",
				Aliases:  []string{"d"},
				Usage:    "Domain name to check",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Port number to check",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			port := cmd.String("port")
			if port == "" {
				port = "80"
			}

			domain := cmd.String("domain")
			status := check(domain, port)
			fmt.Println(status)
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func check(destination string, port string) string {
	address := destination + ":" + port
	timeout := 5 * time.Second
	con, err := net.DialTimeout("tcp", address, timeout)

	var status string
	if err != nil {
		status = fmt.Sprintf("[DOWN] %v is unreachable, \n Error:%v", destination, err)
	} else {
		defer con.Close()
		status = fmt.Sprintf("[UP] %v is reachable, \n From:%v\n To:%v", destination,
			con.LocalAddr(), con.RemoteAddr())
	}
	return status
}


