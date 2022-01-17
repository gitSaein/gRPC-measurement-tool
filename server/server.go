package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"

	pb "gRPC_measurement_tool/protos"
	health_pb "gRPC_measurement_tool/protos/health"
	proto "gRPC_measurement_tool/protos/health"

	"github.com/shirou/gopsutil/cpu"
	"github.com/vys/go-humanize"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":443"
)

var listener *net.Listener

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type Server struct {
	proto           string
	addr            string
	networkListener net.Listener
	mu              sync.Mutex
	statusMap       map[string]health_pb.HealthCheckResponse_ServingStatus
}

func GoRuntimeStats() {
	m := &runtime.MemStats{}
	for {
		percentage, err := cpu.Percent(0, true)

		if err != nil {
			continue
		}
		for idx, cpupercent := range percentage {
			log.Print("Current CPU " + strconv.Itoa(idx) + " utilization:" + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%")
		}

		log.Println("# goroutines: ", runtime.NumGoroutine())
		log.Println("CPU            : ", runtime.NumCPU())
		runtime.ReadMemStats(m)
		log.Println("Memory Acquired: ", humanize.Bytes(m.Sys))
		log.Println("Memory Used    : ", humanize.Bytes(m.Alloc))
		log.Println("# malloc       : ", m.Mallocs)
		log.Println("# free         : ", m.Frees)
		log.Println("GC enabled     : ", m.EnableGC)
		log.Println("# GC           : ", m.NumGC)
		log.Println("Last GC time   : ", m.LastGC)
		log.Println("Next GC        : ", humanize.Bytes(m.NextGC))
		log.Println("--------------------------------------------")
		time.Sleep(2 * time.Second)

		//runtime.GC()
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	// deadline
	for i := 0; i < 5; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, status.Errorf(codes.DeadlineExceeded, "HelloworldService.SayHello DeadlineExceeded")
		}

		time.Sleep(10 * time.Second)
	}

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *Server) Check(ctx context.Context, in *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.Service == "" {
		// check the server overall health status.
		return &proto.HealthCheckResponse{
			Status: proto.HealthCheckResponse_SERVING,
		}, nil
	}
	if status, ok := s.statusMap[in.Service]; ok {
		return &proto.HealthCheckResponse{
			Status: status,
		}, nil
	}

	return nil, status.Error(codes.NotFound, "unknown service")
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		statusMap: make(map[string]proto.HealthCheckResponse_ServingStatus),
	}
}

func ConnHandler(conn net.Conn) {
	recvBuf := make([]byte, 4096) // receive buffer: 4kB
	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Printf("fail to receive data; err: %v", err)
			return
		}
		if 0 < n {
			data := recvBuf[:n]
			log.Println(string(data))
		}
	}
}

func ListenAndGrpcServer() {
	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair("../cert2/server.pem", "../cert2/server.key")
	if err != nil {
		log.Fatalf("LoadX509KeyPair failed to listen: %v", err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}

func main() {
	// go GoRuntimeStats()
	ListenAndGrpcServer()

}
