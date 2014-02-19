package irc

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const (
	R = '→'
	W = '←'
)

type Socket struct {
	conn    net.Conn
	done    chan bool
	reader  *bufio.Reader
	receive chan string
	send    chan string
	writer  *bufio.Writer
}

func NewSocket(conn net.Conn) *Socket {
	socket := &Socket{
		conn:    conn,
		done:    make(chan bool, 1),
		reader:  bufio.NewReader(conn),
		receive: make(chan string, 16),
		send:    make(chan string, 16),
		writer:  bufio.NewWriter(conn),
	}

	go socket.readLines()
	go socket.writeLines()

	return socket
}

func (socket *Socket) String() string {
	return socket.conn.RemoteAddr().String()
}

func (socket *Socket) Close() {
	socket.done <- true
}

func (socket *Socket) Read() <-chan string {
	return socket.receive
}

func (socket *Socket) Write(lines ...string) {
	for _, line := range lines {
		socket.send <- line
	}
	return
}

func (socket *Socket) readLines() {
	for {
		line, err := socket.reader.ReadString('\n')
		if socket.isError(err, R) {
			break
		}

		line = strings.TrimRight(line, "\r\n")
		if len(line) == 0 {
			continue
		}
		if DEBUG_NET {
			log.Printf("%s → %s", socket, line)
		}

		socket.receive <- line
	}

	close(socket.receive)
	socket.Close()
}

func (socket *Socket) writeLines() {
	done := false
	for !done {
		select {
		case line := <-socket.send:
			if _, err := socket.writer.WriteString(line); socket.isError(err, W) {
				break
			}

			if err := socket.writer.Flush(); socket.isError(err, W) {
				break
			}
			if DEBUG_NET {
				log.Printf("%s ← %s", socket, line)
			}

		case done = <-socket.done:
			if DEBUG_NET {
				log.Printf("%s done", socket)
			}
			continue
		}
	}

	if done {
		socket.conn.Close()
	}

	if DEBUG_NET {
		log.Printf("%s closed", socket)
	}

	// read incoming messages and discard to avoid hangs
	for {
		select {
		case <-socket.send:
		case <-socket.done:
		}
	}
}

func (socket *Socket) isError(err error, dir rune) bool {
	if err != nil {
		if DEBUG_NET {
			log.Printf("%s %c error: %s", socket, dir, err)
		}
		return true
	}
	return false
}
