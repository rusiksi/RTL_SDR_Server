package netWorker

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

var tcpServer Server

func init() {

	config := NewConfig()
	config.Address = "127.0.0.1:62001"
	tcpServer, err := NewServer(config)

	if err != nil {
		log.Println("error starting TCP server")
		return
	}
	go func() {
		tcpServer.Run()
	}()
	time.Sleep(2 * time.Second)
}

func TestNETServer_Running(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections.
	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", "127.0.0.1:62001"},
		{"tcp", "127.0.0.1:62001"},
		{"tcp", "127.0.0.1:62001"},
	}
	for _, serv := range servers {
		con, err := net.Dial("tcp", serv.addr)

		if err != nil {
			t.Fatal("could not connect to server: ", err)
			return
		}
		if serv.addr != con.RemoteAddr().String() || serv.protocol != con.RemoteAddr().Network() {
			t.Fatalf("got %s->%s; want %s->%s",
				con.LocalAddr().Network(),
				con.RemoteAddr().Network(),
				serv.addr,
				serv.protocol)
		}

		con.Close()
	}
}

func TestNETServer_SingleClient(t *testing.T) {

	var tests = []struct {
		name string
		send []byte
		size int
		want int
	}{
		{"0 - size, 0 - frame", []byte{0x0, 0x0}, 2, 2},
		{"Byte package", []byte{
			0x2d, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xbe, 0xba, 0xfe, 0xca, 0x49, 0x6d, 0x69, 0x74,
			0x5f, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x17, 0x8c, 0xa1, 0x21, 0xac, 0xea, 0x20, 0x0e, 0x9f, 0x03, 0x0b, 0x0e, 0x72, 0x01, 0x00,
			0x00, 0x01, 0x00, 0x00, 0x00, 0xbe, 0xba, 0xad, 0xab, 0x49, 0x6d, 0x69, 0x74, 0x5f, 0x31, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17, 0x8c,
			0xa1, 0x21, 0xac, 0xea, 0x20, 0x0e, 0x9f, 0x03, 0x0b, 0x0e, 0x72, 0x01, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00,
		}, 98, 98},
	}

	type server struct {
		protocol string
		addr     string
	}

	serv := server{"tcp", "127.0.0.1:62001"}

	con, err := net.Dial("tcp", serv.addr)

	if err != nil {
		t.Fatal("could not connect to server: ", err)
		return
	}
	if serv.addr != con.RemoteAddr().String() || serv.protocol != con.RemoteAddr().Network() {
		t.Fatalf("got %s->%s; want %s->%s",
			con.LocalAddr().Network(),
			con.RemoteAddr().Network(),
			serv.addr,
			serv.protocol)
	}

	for _, tt := range tests {

		testname := fmt.Sprintf("%v", tt.name)
		t.Run(testname, func(t *testing.T) {
			ans, errSend := con.Write(tt.send)
			if errSend != nil {
				t.Errorf("error = %s", errSend)
			}
			if ans != tt.want {
				t.Errorf("Tast:%v error. Got %d, Want %d", tt.name, ans, tt.want)
			}
		})
	}
	con.Close()

}

func TestNETServer_SeveralClients(t *testing.T) {

	var tests = []struct {
		name string
		send []byte
		size int
		want int
	}{
		{"0 - size, 0 - frame", []byte{0x0, 0x0}, 2, 2},
		{"Byte package", []byte{
			0x2d, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xbe, 0xba, 0xfe, 0xca, 0x49, 0x6d, 0x69, 0x74,
			0x5f, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x17, 0x8c, 0xa1, 0x21, 0xac, 0xea, 0x20, 0x0e, 0x9f, 0x03, 0x0b, 0x0e, 0x72, 0x01, 0x00,
			0x00, 0x01, 0x00, 0x00, 0x00, 0xbe, 0xba, 0xad, 0xab, 0x49, 0x6d, 0x69, 0x74, 0x5f, 0x31, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17, 0x8c,
			0xa1, 0x21, 0xac, 0xea, 0x20, 0x0e, 0x9f, 0x03, 0x0b, 0x0e, 0x72, 0x01, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00,
		}, 98, 98},
	}

	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", "127.0.0.1:62001"},
		{"tcp", "127.0.0.1:62001"},
		{"tcp", "127.0.0.1:62001"},
		{"tcp", "127.0.0.1:62001"},
		{"tcp", "127.0.0.1:62001"},
	}
	for inx, serv := range servers {
		con, err := net.Dial("tcp", serv.addr)

		if err != nil {
			t.Fatal("could not connect to server: ", err)
			return
		}
		if serv.addr != con.RemoteAddr().String() || serv.protocol != con.RemoteAddr().Network() {
			t.Fatalf("got %s->%s; want %s->%s",
				con.LocalAddr().Network(),
				con.RemoteAddr().Network(),
				serv.addr,
				serv.protocol)
		}

		for _, tt := range tests {
			testname := fmt.Sprintf("%v", tt.name)
			t.Run(testname, func(t *testing.T) {
				ans, errSend := con.Write(tt.send)
				if errSend != nil {
					t.Errorf("error = %s", errSend)
				}
				if ans != tt.want {
					t.Errorf("Test %d:%v error. Got %d, Want %d", inx, tt.name, ans, tt.want)
				}
			})
		}
		con.Close()
	}
}
