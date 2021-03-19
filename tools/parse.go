package tools

import (
	"encoding/json"
)

// ParseJson Parse takes a byte array of masscan json data and unmarshals it into a
// MasscanResult struct.
func ParseJson(content []byte) (*MasscanResult, error) {
	//fmt.Println(string(content))
	var m []Hosts

	err := json.Unmarshal(content, &m)
	if err != nil {
		return nil, err
	}

	var result MasscanResult
	for i := range m {
		//fmt.Printf("%s\t", m[i].IP)
		//fmt.Printf("%d\n", m[i].Ports[0].Port)
		result.Ports = append(result.Ports,
			Ports{Port: m[i].Ports[0].Port,
				Proto:  m[i].Ports[0].Proto,
				Status: m[i].Ports[0].Status,
				Reason: m[i].Ports[0].Reason,
				TTL:    m[i].Ports[0].TTL,
			})
		result.Hosts = append(result.Hosts,
			Hosts{IP: m[i].IP,
				Timestamp: m[i].Timestamp,
			})

	}

	return &result, nil
}

