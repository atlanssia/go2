package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

// 定义 SOCKS5 的协议格式的数据包。
// 解析客户端发来的数据包。
// 根据请求的地址和端口，转发请求到目标服务器。
// 处理目标服务器的响应，并将其发送回客户端。
// 支持多客户端并发。
// 实现请求追踪日志打印功能。

const (
	SOCKS5Version = 0x05
	SOCKS5Auth   = 0x01
	SOCKS5UnAuth = 0x02
	SOCKS5Reserved = 0x00000000
	SOCKS5ATYPIPv4 = 0x01
	SOCKS5ATYPIPv6 = 0x04
	SOCKS5ATYPDomain = 0x03
	SOCKS5ATYPFQDN = 0x0F
	SOCKS5ATYPUnknown = 0x10

	SOCKS5CommandConnect   = 0x01
	SOCKS5CommandBind      = 0x02
	SOCKS5CommandUDPAssoc  = 0x03

	SOCKS5ResponseSuccess = 0x00
	SOCKS5ResponseServerFailure = 0x01
	SOCKS5ResponseNetworkUnreachable = 0x02
	SOCKS5ResponseHostUnreachable = 0x03
	SOCKS5ResponseConnectionRefused = 0x04
	SOCKS5ResponseTTLExpired = 0x05
	SOCKS5ResponseCommandNotSupported = 0x07
	SOCKS5ResponseAddressTypeNotSupported = 0x08

	SOCKS5LogLevelError = 0x01
	SOCKS5LogLevelWarn = 0x02
	SOCKS5LogLevelInfo = 0x03
	SOCKS5LogLevelDebug = 0x04
)

type Socks5Request struct {
	Version    byte
	NMethods  byte
	Methods   []byte
	Command   byte
	Reserved  byte
	AddressType  byte
	Address    []byte
	Port       []byte
}

type Socks5Response struct {
	Version    byte
	Response  byte
	Reserved  byte
	AddressType  byte
	Address    []byte
	Port       []byte
}

type Socks5Server struct {
	Addr      string
	Port      int
	conns     map[string]net.Conn
	lock     sync.Mutex
	loglevel byte
}

func NewSocks5Server(addr string, port int, loglevel byte) *Socks5Server {
	return &Socks5Server{
		Addr:      addr,
		Port:      port,
		conns:    make(map[string]net.Conn),
		loglevel: loglevel,
	}
}

func (s *Socks5Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return err
	}
	log.Printf("socks5 server started on %s:%d", s.Addr, s.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}
		s.handleConnection(conn)
	}
}

func (s *Socks5Server) handleConnection(conn net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	id := uuid.New().String()
	s.conns[id] = conn
	log.Printf("new connection: %s", id)
	go s.handleSocks5(id, conn)
}

func (s *Socks5Server) handleSocks5(id string, conn net.Conn) {
	defer func() {
		s.lock.Lock()
		delete(s.conns, id)
		s.lock.Unlock()
		log.Printf("connection closed: %s", id)
		conn.Close()
	}()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error reading from client: %v", err)
			return
		}
		if s.loglevel >= SOCKS5LogLevelDebug {
			log.Printf("received from client: %s", hex(buf[:n]))
		}
		req := parseSocks5Request(buf[:n])
		if req == nil {
			log.Printf("invalid socks5 request")
			return
		}
		resp := s.handleRequest(req)
		if resp == nil {
			log.Printf("error handling request")
			return
		}
		if s.loglevel >= SOCKS5LogLevelInfo {
			log.Printf("sending to client: %s", hex(resp.Bytes()))
		}
		_, err = conn.Write(resp.Bytes())
		if err != nil {
			log.Printf("error writing to client: %v", err)
			return
		}
	}
}

func (s *Socks5Server) handleRequest(req *Socks5Request) *Socks5Response {
	// TODO: Implement the logic to handle the request and forward it to the target address
	// Return a response or nil if error occurs
}

func parseSocks5Request(buf []byte) *Socks5Request {
	// TODO: Implement the parsing logic for the SOCKS5 request
	// Return a Socks5Request object or nil if the request is invalid
}

func hex(b []byte) string {
	hex := "0x"
	for _, byte := range b {
		hex += fmt.Sprintf("%02x", byte)
	}
	return hex
}

func main() {
	server := NewSocks5Server("127.0.0.1", 1080, SOCKS5LogLevelInfo)
	server.Start()
}
