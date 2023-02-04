package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	fileName, port := handleArgs()
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip, ipNet, _ := net.ParseCIDR(scanner.Text())
		ones, bits := ipNet.Mask.Size()
		wg.Add(1 << uint(bits-ones))
		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			go scanHost(net.JoinHostPort(ip.String(), port), &wg)
		}
		wg.Wait()
	}
}

func handleArgs() (string, string) {
	var fileName string
	var port string
	flag.StringVar(&fileName, "f", "", "file name that contains network addresses as ip/cidr")
	flag.StringVar(&port, "p", "", "port number to scan")
	flag.Parse()
	if fileName == "" || port == "" {
		fmt.Println("file name and port number are required")
		os.Exit(1)
	}
	return fileName, port
}

func scanHost(address string, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err == nil && conn != nil {
		defer conn.Close()
		fmt.Println(address)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
