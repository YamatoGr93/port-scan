package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func scanPort(protocol, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	hostname := "localhost"
	ip, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println("Error getting IP address:", err)
		return
	}
	ipAddress := ip[0].String()

	fmt.Println("Starting port scan...")

	// Create reports directory if it doesn't exist
	if _, err := os.Stat("reports"); os.IsNotExist(err) {
		err = os.Mkdir("reports", 0755)
		if err != nil {
			fmt.Println("Error creating reports directory:", err)
			return
		}
	}

	// Create results file with current date and time in the filename
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("reports/results_%s.txt", currentTime)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write header information
	file.WriteString(fmt.Sprintf("Target: %s (%s)\n", hostname, ipAddress))
	file.WriteString(fmt.Sprintf("Report Date and Time: %s\n", currentTime))

	totalPorts := 65535
	openPorts := 0
	closedPorts := 0
	openPortsList := ""

	for port := 1; port <= totalPorts; port++ {
		if scanPort("tcp", hostname, port) {
			openPorts++
			result := fmt.Sprintf("Port %d is open\n", port)
			fmt.Print(result)
			file.WriteString(result)
			openPortsList += fmt.Sprintf("Port %d\n", port)
		} else {
			closedPorts++
		}
	}

	// Write summary information
	file.WriteString(fmt.Sprintf("\nTotal ports scanned: %d\n", totalPorts))
	file.WriteString(fmt.Sprintf("Open ports: %d\n", openPorts))
	file.WriteString(fmt.Sprintf("Closed ports: %d\n", closedPorts))

	// Write open ports list
	if openPorts > 0 {
		file.WriteString("\nOpen ports found:\n")
		file.WriteString(openPortsList)
	}

	fmt.Println("Port scan completed.")
}
