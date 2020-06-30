package main

import (
	"fmt"
	"net"

	wol "github.com/sabhiram/go-wol"
	log "github.com/sirupsen/logrus"
)

const WOL_PORT = 9

func GetOutboundUDPAddr() (*net.UDPAddr, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr, nil
}

func SendWakeOnLAN(macAddr string) error {
	localAddr, err := GetOutboundUDPAddr()
	if err != nil {
		return err
	}

	bcastAddr := fmt.Sprintf("%v:%v", "255.255.255.255", WOL_PORT)
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := wol.New(macAddr)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		return fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}

	if err != nil {
		return err
	}

	log.Infof("Magic packet sent successfully to %s", macAddr)
	return nil
}
