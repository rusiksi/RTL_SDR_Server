package netWorker

import (
	"RTL_SDR_Server/internal/pkg/pkgProcessor"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

// Интерфейс TCP и UDP сервера
type Server interface {
	Run() error
	Close() error
}

// Фабричный метод для создания инстанса конкретного класса сервера
func NewServer(config *Config) (Server, error) {

	if config == nil {
		return nil, errors.New("Invalid config")
	}

	switch strings.ToLower(config.Protocol) {
	case "tcp":
		return &ServerTcpImpl{
			addr: config.Address,
		}, nil
	case "udp":
		return &ServerUdpImpl{
			addr: config.Address,
		}, nil
	}
	return nil, errors.New("Invalid protocol given")
}

// ServerTcpImpl - реализация интерфейса Server для TCP/IP подключений
type ServerTcpImpl struct {
	addr   string
	server net.Listener
}

// Запуск сервера.
// Операция блокирующая. Сервер в бесконечном цикле ждёт новых подключений
func (t *ServerTcpImpl) Run() (err error) {
	t.server, err = net.Listen("tcp", t.addr)

	if err != nil {
		return err
	}

	defer t.Close()

	return t.waitNewConnections()
}

// Закрытие сервера
func (t *ServerTcpImpl) Close() (err error) {
	return t.server.Close()
}

// функция ожидания новых подключений.
// Операция блокирующая. Сервер в бесконечном цикле ждёт новых подключений
func (t *ServerTcpImpl) waitNewConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			return err
		}
		// обработчик каждого нового подключения запускается в новой горутине
		go t.handleConnection(conn)
	}
	return nil
}

func (t *ServerTcpImpl) handleConnection(conn net.Conn) {
	defer conn.Close()

	name := conn.RemoteAddr().String()
	log.Printf("%+v connected to server\n", name)

	// каждое новое подключение получает свой экземпляр процессора пакетов
	var proc pkgProcessor.IPkgProcessor = new(pkgProcessor.PkgProcessorImpl)
	//конкретный инстанс процессора сам знает как правильно читать и распаковывать данные
	// Операция блокирующая, выполняется до тех пор,пока клиент не отключится или не возникнет ошибка передачи
	err := proc.ReadData(conn)
	if err != nil {
		log.Printf("Error in IPkgProcessor.ReadData = %v, close connection", err)
	}
	proc = nil
}

// реализация интерфейса Server для Udp протокола
type ServerUdpImpl struct {
	addr   string
	server *net.UDPConn
}

// Запуск службы для прослушивания udp
// Операция блокирующая.
func (u *ServerUdpImpl) Run() (err error) {
	laddr, err := net.ResolveUDPAddr("udp", u.addr)
	if err != nil {
		return errors.New("could not resolve UDP addr")
	}

	u.server, err = net.ListenUDP("udp", laddr)
	if err != nil {
		return errors.New("could not listen on UDP")
	}

	return u.waitNewConnections()
}

// Ожиданние новых подключений
// Операция блокирующая.
func (u *ServerUdpImpl) waitNewConnections() error {
	var err error
	for {
		buf := make([]byte, 2048)
		n, conn, err := u.server.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			break
		}
		if conn == nil {
			continue
		}

		go u.handleConnection(conn, buf[:n])
	}
	return err
}

//обработчик новых подключений. Выполняет обычный эхо - ответ
func (u *ServerUdpImpl) handleConnection(addr *net.UDPAddr, cmd []byte) {
	u.server.WriteToUDP([]byte(fmt.Sprintf("Request recieved: %s", cmd)), addr)
}

//Остановка сервиса работы с udp
func (u *ServerUdpImpl) Close() error {
	return u.server.Close()
}
