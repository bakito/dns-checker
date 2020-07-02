package manualdns

// thanks to https://github.com/vishen/go-dnsquery

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

/*
	Example header:

	AA AA - ID
	01 00 - Query parameters (QR | Opcode | AA | TC | RD | RA | Z | ResponseCode)
	00 01 - Number of questions
	00 00 - Number of answers
	00 00 - Number of authority records
	00 00 - Number of additional records
*/

type dnsQuery struct {
	ID     uint16 // An arbitrary 16 bit request identifier (same id is used in the response)
	QR     bool   // A 1 bit flat specifying whether this message is a query (0) or a response (1)
	Opcode uint8  // A 4 bit fields that specifies the query type; 0 (standard), 1 (inverse), 2 (status), 4 (notify), 5 (update)

	AA           bool  // Authoritative answer
	TC           bool  // 1 bit flag specifying if the message has been truncated
	RD           bool  // 1 bit flag to specify if recursion is desired (if the DNS server we secnd out request to doesn't know the answer to our query, it can recursively ask other DNS servers)
	RA           bool  // Recursive available
	Z            uint8 // Reserved for future use
	ResponseCode uint8

	QDCount uint16 // Number of entries in the question section
	ANCount uint16 // Number of answers
	NSCount uint16 // Number of authorities
	ARCount uint16 // Number of additional records

	Questions []dnsQuestion
}

func (q dnsQuery) encode() []byte {

	q.QDCount = uint16(len(q.Questions))

	var buffer bytes.Buffer

	_ = binary.Write(&buffer, binary.BigEndian, q.ID)

	b2i := func(b bool) int {
		if b {
			return 1
		}

		return 0
	}

	queryParams1 := byte(b2i(q.QR)<<7 | int(q.Opcode)<<3 | b2i(q.AA)<<1 | b2i(q.RD))
	queryParams2 := byte(b2i(q.RA)<<7 | int(q.Z)<<4)

	_ = binary.Write(&buffer, binary.BigEndian, queryParams1)
	_ = binary.Write(&buffer, binary.BigEndian, queryParams2)
	_ = binary.Write(&buffer, binary.BigEndian, q.QDCount)
	_ = binary.Write(&buffer, binary.BigEndian, q.ANCount)
	_ = binary.Write(&buffer, binary.BigEndian, q.NSCount)
	_ = binary.Write(&buffer, binary.BigEndian, q.ARCount)

	for _, question := range q.Questions {
		buffer.Write(question.encode())
	}

	return buffer.Bytes()
}

/*
	Example Question:

	07 65 - 'example' has length 7, e
	78 61 - x, a
	6D 70 - m, p
	6C 65 - l, e
	03 63 - 'com' has length 3, c
	6F 6D - o, m
	00    - zero byte to end the QNAME
	00 01 - QTYPE
	00 01 - QCLASS

	76578616d706c6503636f6d0000010001
*/

type dnsQuestion struct {
	Domain string
	Type   uint16 // DNS Record type we are looking up; 1 (A record), 2 (authoritive name server)
	Class  uint16 // 1 (internet)
}

func (q dnsQuestion) encode() []byte {
	var buffer bytes.Buffer

	domainParts := strings.Split(q.Domain, ".")
	for _, part := range domainParts {
		if err := binary.Write(&buffer, binary.BigEndian, byte(len(part))); err != nil {
			log.Fatalf("Error binary.Write(..) for '%s': '%s'", part, err)
		}

		for _, c := range part {
			if err := binary.Write(&buffer, binary.BigEndian, uint8(c)); err != nil {
				log.Fatalf("Error binary.Write(..) for '%s'; '%c': '%s'", part, c, err)
			}
		}
	}

	_ = binary.Write(&buffer, binary.BigEndian, uint8(0))
	_ = binary.Write(&buffer, binary.BigEndian, q.Type)
	_ = binary.Write(&buffer, binary.BigEndian, q.Class)

	return buffer.Bytes()

}

func resolve(ctx context.Context, query []byte, dnsServer string) (byte, error) {
	// Setup a UDP connection
	var d net.Dialer
	conn, err := d.DialContext(ctx, "udp", dnsServer)
	if err != nil {
		return 0, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	_, _ = conn.Write(query)

	encodedAnswer := make([]byte, len(query))
	if _, err := bufio.NewReader(conn).Read(encodedAnswer); err != nil {
		return 0, err
	}

	return encodedAnswer[3] & 0xF, nil

}

func responseCode(responseCode byte) (string, error) {

	switch responseCode {
	case 0:
		return "Domain exists!", nil
	case 1:
		return "", errors.New("format error")
	case 2:
		return "", errors.New("server failure")
	case 3:
		return "", errors.New("non-existent domain")
	case 9:
		return "", errors.New("server not authoritative for zone")
	case 10:
		return "", errors.New("name not in zone")
	default:
		return "", fmt.Errorf("unmapped response code for '%d'", responseCode)
	}
}
