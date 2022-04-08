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
# Dependencies
- `go` (> `1.10`)
- masscan (> `1.3.0`)
  - [download windows masscan](https://github.com/zan8in/masscan/blob/main/bin/masscan-win/masscan.exe)
  - [download linux masscan](https://github.com/zan8in/masscan/blob/main/bin/masscan-linux/masscan)


# Table of content
- masscan run timeout setting
- Get process ID (PID)
- Select network interface automatically
# TODO
- Support more parameters
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
# Async scan example
```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zan8in/masscan"
)

func main() {
	context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("60.10.116.10"),
		masscan.SetParamPorts("80"),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(50),
		masscan.WithContext(context),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		panic(err)
	}

	stdout := scanner.GetStdout()

	stderr := scanner.GetStderr()

	go func() {
		for stdout.Scan() {
			srs := masscan.ParseResult(stdout.Bytes())
			fmt.Println(srs.IP, srs.Port)
			scannerResult = append(scannerResult, srs)
		}
	}()

	go func() {
		for stderr.Scan() {
			fmt.Println("err: ", stderr.Text())
			errorBytes = append(errorBytes, stderr.Bytes()...)
		}
	}()

	if err := scanner.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("masscan result count : ", len(scannerResult), " PID : ", scanner.GetPid())

}
```
The program above outputs:
```
C:\masscan\masscan.exe 146.56.202.100-146.56.202.200 -p 3306 --wait=0 --rate=2000
err:  Starting masscan 1.3.2 (http://bit.ly/14GZzcT) at 2021-03-19 14:52:27 GMT
err:  Initiating SYN Stealth Scan
err:  Scanning 101 hosts [1 port/host]
146.56.202.115 3306
146.56.202.190 3306
146.56.202.188 3306
146.56.202.125 3306
146.56.202.185 3306
146.56.202.117 3306
146.56.202.112 3306
146.56.202.161 3306
146.56.202.165 3306
146.56.202.166 3306
                                                                             
masscan result count :  10

Process finished with exit code 0
```
# The development soul comes from
[Ullaakut](https://github.com/Ullaakut/nmap)

# Special thanks 
李雪松 XueSong Lee
