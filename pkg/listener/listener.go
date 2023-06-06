package listener

import (
	"fmt"
	"net"

	"github.com/IMaloney/snowcast/pkg/utils"
)

type UDPListener struct {
	conn     *net.UDPConn
	exitChan chan struct{}
}

// CreateUDPListener creates a udp listener
func CreateUDPListener(port string) (*UDPListener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return &UDPListener{
		conn:     conn,
		exitChan: make(chan struct{}, 1),
	}, nil
}

// Quit quits the UDP listener
func (l *UDPListener) Quit() {
	l.exitChan <- struct{}{}
}

func (l *UDPListener) Listen() {
	for {
		select {
		case <-l.exitChan:
			l.conn.Close()
			return
		default:
			buffer := make([]byte, utils.SONGCHUNK)
			bytesRead, err := l.conn.Read(buffer)
			if err != nil {
				continue
			}
			fmt.Printf("%s", buffer[:bytesRead])
		}

	}
}
