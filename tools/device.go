package tools

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// reference : https://github.com/boy-hack/ksubdomain/blob/426067a5eb3ada94e365de7256d93ddf3c9d8e4a/core/device/device.go#L15

type SelfMac net.HardwareAddr

type EtherTable struct {
	SrcIp  net.IP  `yaml:"src_ip"`
	Device string  `yaml:"device"`
	SrcMac SelfMac `yaml:"src_mac"`
	DstMac SelfMac `yaml:"dst_mac"`
}

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func AutoGetDevices() (*EtherTable, error) {
	domain := CreateRandomString(4) + ".baidu.com"
	signal := make(chan *EtherTable)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取网络设备失败:%s\n", err.Error()))
	}
	data := make(map[string]net.IP)
	keys := []string{}
	for _, d := range devices {
		for _, address := range d.Addresses {
			ip := address.IP
			if ip.To4() != nil && !ip.IsLoopback() {
				data[d.Name] = ip
				keys = append(keys, d.Name)
			}
		}
	}
	ctx := context.Background()
	// 在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)
	for _, drviceName := range keys {
		go func(drviceName string, domain string, ctx context.Context) {
			var (
				snapshot_len int32         = 1024
				promiscuous  bool          = false
				timeout      time.Duration = -1 * time.Second
				handle       *pcap.Handle
			)
			var err error
			handle, err = pcap.OpenLive(
				drviceName,
				snapshot_len,
				promiscuous,
				timeout,
			)
			if err != nil {
				return
			}
			defer handle.Close()
			// Use the handle as a packet source to process all packets
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for {
				select {
				case <-ctx.Done():
					return
				default:
					packet, err := packetSource.NextPacket()
					if err != nil {
						continue
					}
					if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
						dns, _ := dnsLayer.(*layers.DNS)
						if !dns.QR {
							continue
						}
						for _, v := range dns.Questions {
							if string(v.Name) == domain {
								ethLayer := packet.Layer(layers.LayerTypeEthernet)
								if ethLayer != nil {
									eth := ethLayer.(*layers.Ethernet)
									etherTable := EtherTable{
										SrcIp:  data[drviceName],
										Device: drviceName,
										SrcMac: SelfMac(eth.DstMAC),
										DstMac: SelfMac(eth.SrcMAC),
									}
									signal <- &etherTable
									return
								}
							}
						}
					}
				}
			}
		}(drviceName, domain, ctx)
	}
	for {
		select {
		case c := <-signal:
			cancel()
			fmt.Print("\n")
			return c, nil
		default:
			_, _ = net.LookupHost(domain)
			time.Sleep(time.Second * 1)
		}
	}
}
