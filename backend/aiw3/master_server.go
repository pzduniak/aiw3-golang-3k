package aiw3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Result struct {
	IP1   uint8
	IP2   uint8
	IP3   uint8
	IP4   uint8
	Port1 uint8
	Port2 uint8
}

func QueryMasterServer(ipaddr string) []string {
	address, err := net.ResolveUDPAddr("udp", ipaddr)
	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, "\xFF\xFF\xFF\xFFgetservers IW4 61586 full empty\x00")

	var total string

	for {
		var packet [5000]byte
		n, err := conn.Read(packet[:])
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}

		total += string(packet[:n])

		if strings.Contains(total, "EOT") {
			break
		}
	}

	found := []string{}

	parts := strings.Split(total, "\\")
	for _, part := range parts[1:] {
		if !strings.Contains(part, "EOT") {
			buf := bytes.NewBuffer([]byte(part))

			var result Result
			binary.Read(buf, binary.BigEndian, &result)

			port := int(result.Port1)*256 + int(result.Port2)

			ip := strconv.Itoa(int(result.IP1)) + "." +
				strconv.Itoa(int(result.IP2)) + "." +
				strconv.Itoa(int(result.IP3)) + "." +
				strconv.Itoa(int(result.IP4)) + ":" +
				strconv.Itoa(port)

			found = append(found, ip)
		}
	}

	return found
}
