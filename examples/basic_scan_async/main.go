package main

import (
	"fmt"
	"github.com/zan8in/masscan"
	"log"
)

// BESTPORTS "21,22,80,U:137,U:161,443,445,U:1900,3306,3389,U:5353,8080"
const BESTPORTS = `
	21,22,23,69,80,443,445,1433,1434,1521,1158,210,8080,8009,9080,9081,9090,7001,
	7002,4848,8983,1352,3306,6379,5432,27001,5000,4100,4200,11211,9200,9300,50010,50070,
	2181,2049,137,389,3389,5900,5901,5632,6000,25,465,110,995,109,143,993,53,67,68,161,
	512,513,514,873,8069,1090,1099,2375,161,135,139,1883,6666,6667,7777,8161,9000,9001,
	12345,27017,1080
	`
func main() {
	var (
		scannerResult []masscan.ScannerResult
		errorBytes  []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("146.56.202.100-146.56.202.200"),
		masscan.SetParamPorts("3306"),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(2000),
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

	if err := scanner.Wait(); err !=nil {
		panic(err)
	}

	fmt.Println("masscan result count : ", len(scannerResult))

}
