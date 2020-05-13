package netWorker

import (
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
)

var tcp, udp Server
var wg sync.WaitGroup
var pkgChn chan []byte

func init() {
	// Start the new server
	tcp, err := NewServer("tcp", "127.0.0.1:62001")

	pkgChn = make(chan []byte)
	if err != nil {
		log.Println("error starting TCP server")
		return
	}
	wg.Add(1)
	go func() {
		tcp.Run(&wg, pkgChn)
	}()
	wg.Wait()
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
		send []byte
		size int
		want int
	}{
		{[]byte{0x0, 0x0}, 0, 0},
		{[]byte{0x1, 0x2}, 2, 2},
		{[]byte{}, 0, 0},
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

		testname := fmt.Sprintf("%d,%d", tt.size, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans, errSend := con.Write(tt.send)
			if errSend != nil {
				t.Errorf("error = %s", errSend)
			}
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
	fmt.Println(<-pkgChn)
	con.Close()

}

func TestNETServer_SeveralClients(t *testing.T) {

	var tests = []struct {
		send []byte
		size int
		want int
	}{
		{[]byte{0x0, 0x0}, 0, 0},
		{[]byte{}, 0, 0},
		{[]byte{
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

		for _, tt := range tests {

			testname := fmt.Sprintf("%d,%d", tt.size, tt.want)
			t.Run(testname, func(t *testing.T) {
				ans, errSend := con.Write(tt.send)
				if errSend != nil {
					t.Errorf("error = %s", errSend)
				}
				if ans != tt.want {
					t.Errorf("got %d, want %d", ans, tt.want)
				}
			})
		}
		fmt.Println(<-pkgChn)
		con.Close()
	}
}
