# masscan
![Just I love it](https://github.com/zan8in/masscan/blob/main/assets/golang.png)
Masscan is a golang library to run masscan scans, parse scan results. 
# What is masscan
Masscan is an Internet-scale port scanner. It can scan the entire Internet in under 5 minutes, transmitting 10 million packets per second, from a single machine.

Its usage (parameters, output) is similar to nmap, the most famous port scanner. When in doubt, try one of those features -- features that support widespread scanning of many machines are supported, while in-depth scanning of single machines aren't.

Internally, it uses asynchronous transmission, similar to port scanners like scanrand, unicornscan, and ZMap. It's more flexible, allowing arbitrary port and address ranges.
# Installation
```
go get github.com/zan8in/masscan
```
to install the package
```
import "github.com/zan8in/masscan"
```
# TODO
- Support more parameters 
- Add asynchronous scan
# Simple example
```go
package main

import (
	"fmt"
	"log"

	"github.com/zan8in/masscan"
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
```
The program above outputs:
```
/usr/bin/masscan 146.56.202.100/24 -p 80 --wait=0 --rate=10000 -oJ -
Host: 146.56.202.15 Port: 80
Host: 146.56.202.251 Port: 80
Host: 146.56.202.112 Port: 80
...
...
Host: 146.56.202.17 Port: 80
Host: 146.56.202.209 Port: 80
Host: 146.56.202.190 Port: 80
Host: 146.56.202.222 Port: 80
Host: 146.56.202.207 Port: 80
hosts len :  37
```
# The development soul comes from
[Ullaakut](https://github.com/Ullaakut/nmap)

# Special thanks 
李雪松 XueSong Lee
