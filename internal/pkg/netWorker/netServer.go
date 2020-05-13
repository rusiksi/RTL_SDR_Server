package netWorker

import (
	"RTL_SDR_Server/internal/pkg/pkgProcessor"
	"errors"
	"fmt"
	"net"
	"strings"
)

// Server defines the minimum contract our
// TCP and UDP server implementations must satisfy.
type Server interface {
	Run() error
	Close() error
}

// NewServer creates a new Server using given protocol
// and addr.
func NewServer(config * Config) (Server, error) {

	if config == nil {
		return nil, errors.New("Invalid config")
	}

	switch strings.ToLower(config.Protocol) {
	case "tcp":
		return &TCPServer{
			addr: config.Address,
		}, nil
		//case "udp":
		//	return &UDPServer{
		//		addr: addr,
		//	}, nil
	}
	return nil, errors.New("Invalid protocol given")
}

// TCPServer holds the structure of our TCP
// implementation.
type TCPServer struct {
	addr   string
	server net.Listener
}

// Run starts the TCP Server.
func (t *TCPServer) Run() (err error) {
	t.server, err = net.Listen("tcp", t.addr)

	if err != nil {
		return err
	}

	defer t.Close()

	return t.waitNewConnections()
}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

func (t *TCPServer) waitNewConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			break
		}

		go t.handleConnection(conn)
	}
	return
}

func (t *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	//TODO: логгирование
	name := conn.RemoteAddr().String()
	fmt.Printf("%+v connected\n", name)

	var proc pkgProcessor.IPkgProcessor = new(pkgProcessor.PkgProcessorImpl)
	err := proc.ReadData(conn)
	if err != nil {
		fmt.Println(err)
	}
	proc = nil
}

// UDPServer holds the necessary structure for our
// UDP server.
type UDPServer struct {
	addr   string
	server *net.UDPConn
}

//// Run starts the UDP server.
//func (u *UDPServer) Run() (err error) {
//	laddr, err := net.ResolveUDPAddr("udp", u.addr)
//	if err != nil {
//		return errors.New("could not resolve UDP addr")
//	}
//
//	u.server, err = net.ListenUDP("udp", laddr)
//	if err != nil {
//		return errors.New("could not listen on UDP")
//	}
//
//	return u.waitNewConnections()
//}

//func (u *UDPServer) waitNewConnections() error {
//	var err error
//	for {
//		buf := make([]byte, 2048)
//		n, conn, err := u.server.ReadFromUDP(buf)
//		if err != nil {
//			log.Println(err)
//			break
//		}
//		if conn == nil {
//			continue
//		}
//
//		go u.handleConnection(conn, buf[:n])
//	}
//	return err
//}

//func (u *UDPServer) handleConnection(addr *net.UDPAddr, cmd []byte) {
//	u.server.WriteToUDP([]byte(fmt.Sprintf("Request recieved: %s", cmd)), addr)
//}
//
//// Close ensures that the UDPServer is shut down gracefully.
//func (u *UDPServer) Close() error {
//	return u.server.Close()
//}
