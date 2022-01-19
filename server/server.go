package main

import (
	"context"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	interceptor "gRPC_measurement_tool/interceptors"
)

const (
	port = ":50051"
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
	// for i := 0; i < 5; i++ {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return nil, status.Errorf(codes.DeadlineExceeded, "HelloworldService.SayHello DeadlineExceeded")
	// 	}

	// 	time.Sleep(10 * time.Second)
	// }

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
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor))
	srv := NewServer()

	health_pb.RegisterHealthServer(s, srv)
	reflection.Register(s)

	pb.RegisterGreeterServer(s, &server{}) // helloworld_grpc.pb.go 에 있음

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //grpc 서버 시작
		log.Fatalf("failed to serve: %v", err)
	}

}

func main() {
	// go GoRuntimeStats()
	ListenAndGrpcServer()

}
