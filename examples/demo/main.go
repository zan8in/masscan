package main

import (
	"fmt"
	"github.com/zan8in/masscan"
	"log"
)

// Example
func main() {
	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("146.56.202.100/24"),
		masscan.SetParamPorts("80"),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(10000),
	)
	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	scanResult, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("masscan encountered an error: %v", err)
	}

	if scanResult != nil {
		for i, v := range scanResult.Hosts {
			fmt.Printf("Host: %s Port: %v\n", v.IP, scanResult.Ports[i].Port)
		}
		fmt.Println("hosts len : ", len(scanResult.Hosts))
	}

}
