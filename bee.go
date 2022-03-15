package main

import (
	"flag"
	"fmt"
	"time"
	"os"
	"strconv"
	"net"
	"strings"
	"sync"
)

// Bee (version 1.0)
// Created by Angel Pimentel
//
// LinkedIn  https://www.linkedin.com/in/angel-pimentell
// GitHub    https://github.com/angelpimentell
// Twitter   https://twitter.com/angelpimentelll
//
// NOTES:
// Using too many 'threads' can lead to false negatives
// Using by udp protocol can lead to false negatives

// Global variables
var start time.Time
var ip_address string
var ports string
var routines int
var timeout time.Duration
var protocol []string

func arguments() {
	flag.StringVar(&ip_address, "i", "127.0.0.1", "a ip address")
	flag.StringVar(&ports, "p", "1-65535", "a list of ports")
	flag.IntVar(&routines, "t", 200, "routines to use")
	var tima = flag.Int("tm", 1, "timeout to use")
	boolpT := flag.Bool("pT", false, "scanning by tcp protocol")
	boolpU := flag.Bool("pU", false, "scanning by udp protocol")

	flag.Parse()

	timeout = time.Duration(*tima) * time.Second

	if *boolpT == true {
                protocol = append(protocol, "tcp")
        }

	if *boolpU == true {
                protocol = append(protocol, "udp")
        }

	if *boolpU == *boolpT {
		protocol = append(protocol, "tcp")
	}

}

func prepare_data() {

	// Wait gourutines
	var wg sync.WaitGroup

	// Print data target
	fmt.Println("[+] Target:", ip_address)
	fmt.Printf("[+] Ports: %v\n\n", ports)

	// List to divide
	list_ports := make([]string, 0, 350)

	// Commas
	for _, value_com := range strings.Split(ports, ",") {
		if strings.ContainsAny(value_com, "-") {

			data := strings.Split(value_com, "-")
			start, err := strconv.Atoi(data[0])

			if err != nil {
				panic(err)
			}

			finish, err := strconv.Atoi(data[1])
			if err != nil {
				panic(err)
			}

			// Stripes
			for y := start; y <= finish; y++ {
				list_ports = append(list_ports, strconv.Itoa(y))
			}
		} else {
			list_ports = append(list_ports, value_com)
		}

	}

	// Less gourutines than ports
	if len(list_ports) < routines {
		routines = len(list_ports)
	}

	jobs := len(list_ports) / routines

	for j := jobs; j <= len(list_ports); j=j+jobs {

		top := 0
		low := j - jobs

		// Last job
		if j + jobs > len(list_ports) {
			diff := len(list_ports) - j
			top = j + diff

		// Next job
		} else {
			top = j
		}

        	// Increment the WaitGroup counter
                wg.Add(1)

		// Start scan with specified ports
		start = time.Now()
		for _, p := range protocol {
			if p == "tcp" {
				go scan_ports_tcp(list_ports[low:top], &wg)
			} else if p == "udp"{
				go scan_ports_udp(list_ports[low:top], &wg)
			}
		}
	}
	wg.Wait()
}

func scan_ports_tcp(ports_list []string, wg *sync.WaitGroup) {

	defer wg.Done()
	for _, port := range ports_list {

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip_address, port), timeout)

		if err != nil {
			// Do nothing - Connection error
		} else {
			conn.Close()
	            	fmt.Printf("[!] OPEN [%s/TCP]\n", port)
        	}
	}
}


func scan_ports_udp(ports_list []string, wg *sync.WaitGroup) {

        defer wg.Done()
        for _, port := range ports_list {

    		conn, err := net.Dial("udp", ip_address + ":" + port)

		if err != nil {
		        // Do nothing - Connection error
    		} else {

			fmt.Fprintf(conn, "0")
			p :=  make([]byte, 1024)

			conn.SetReadDeadline(time.Now().Add(timeout)) // Timeout
    			_, err = conn.Read(p)

			if err == nil {
				fmt.Printf("[!] OPEN [%s/UDP]\n", port)
			} else {
        			// Do nothing - No response
			}

			conn.Close()
		}
	}
}
func presentation(){
	fmt.Println("\n" +
	"	██████╗ ███████╗███████╗\n" +
	"	██╔══██╗██╔════╝██╔════╝\n" +
	"	██████╔╝█████╗  █████╗\n" +
	"	██╔══██╗██╔══╝  ██╔══╝\n" +
	"	██████╔╝███████╗███████╗\n" +
	"	╚═════╝ ╚══════╝╚══════╝\n" +
	"-----------------------------------------\n" +
	"----------- Created by Kraken -----------\n" +
	"-----------------------------------------\n")
}

func checkTime() {
	now := time.Now()
	diff := now.Sub(start)
	fmt.Printf("\n[+] Time elapsed: %v\n\n", diff)
}

func main() {
	presentation()
	// Usage
	flag.Usage = func() {
  		fmt.Fprintf(os.Stderr, "\n[?] USAGE\n	bee -i <ip-address> -p <ports>\n\n[?] FLAGS\n	-i, IP ADDRESS\n	-p, LIST OF PORTS\n	-pT, SCAN BY TCP PROTOCOL\n	-pU, SCAN BY UDP PROTOCOL\n	-tm, TIMEOUT (Default 1s)\n	-t, THREADS (Default 200)\n\n[!] EXAMPLES\n	bee -i 10.0.0.1 -p 80,443\n	bee -i 10.0.0.1 -p 21,22,80-100\n	bee -i 10.0.0.1 -p 80-100 -pU -pT\n\n")
	}
	arguments()
	prepare_data()
	checkTime()
}
