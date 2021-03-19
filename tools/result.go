package tools

type (
	// MasscanResult masscan output struct
	// eg:
	//[
	//{   "ip": "192.168.88.120",   "timestamp": "1614306482", "ports": [ {"port": 80, "proto": "tcp", "status": "open", "reason": "syn-ack", "ttl": 51} ] }
	//]
	MasscanResult struct {
		Hosts []Hosts `json:"hosts"`
		Ports []Ports `json:"ports"`
	}

	// Hosts  masscan hosts output struct
	Hosts struct {
		IP        string  `json:"ip"`
		Ports     []Ports `json:"ports"`
		Timestamp string  `json:"timestamp"`
	}

	// Ports  masscan ports output struct
	Ports struct {
		Port   int    `json:"port"`
		Proto  string `json:"proto"`
		Status string `json:"status"`
		Reason string `json:"reason"`
		TTL    int    `json:"ttl"`
	}

)
