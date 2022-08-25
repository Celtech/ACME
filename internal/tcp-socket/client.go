package tcp_socket

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type TCPSocket struct {
	con     net.Conn
	Address string
	Port    int
}

func (tcp *TCPSocket) connect() error {
	var err error
	tcp.con, err = net.Dial("tcp", fmt.Sprintf("%s:%d", tcp.Address, tcp.Port))
	log.Infof("Connected to %s:%d", tcp.Address, tcp.Port)

	if err != nil {
		return err
	}

	return nil
}

func (tcp *TCPSocket) close() error {
	err := tcp.con.Close()
	if err != nil {
		return err
	}

	log.Infof("Closed connection to %s:%d", tcp.Address, tcp.Port)
	return nil
}

func (tcp *TCPSocket) Write(command string) error {
	err := tcp.connect()
	if err != nil {
		return fmt.Errorf("error connecting to %s:%d - \n\t%v", tcp.Address, tcp.Port, err)
	}
	defer tcp.close()

	log.Infof("Attempting to write message to %s:%d: %s", tcp.Address, tcp.Port, command)
	go tcp.reader()

	if _, err := tcp.con.Write([]byte(command + "\n")); err != nil {
		return fmt.Errorf("error writing command: \n\t%v", err)
	}

	log.Infof("Message written to %s:%d: %s", tcp.Address, tcp.Port, command)
	time.Sleep(time.Second * 2)
	return nil
}

func (tcp *TCPSocket) reader() {
	connbuf := bufio.NewReader(tcp.con)
	// Read the first byte and set the underlying buffer
	b, _ := connbuf.ReadByte()
	if connbuf.Buffered() > 0 {
		var msgData []byte
		msgData = append(msgData, b)
		for connbuf.Buffered() > 0 {
			// read byte by byte until the buffered data is not empty
			b, err := connbuf.ReadByte()
			if err == nil {
				msgData = append(msgData, b)
			} else {
				log.Errorf("unreadable caracter: %x", b)
			}
		}

		log.Infof(string(msgData[:]))
	}
}
