package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	pb "ex05_grpc/keyvalue"

	"google.golang.org/grpc"
)

func GetBodyTest() {
	// http.Get 내부에서는 DefaultClient.Get 을 호출한다.
	// DefaultClient는 타임아웃이 설정되지 않았다.
	// 타임아웃을 설정하려면 var client = &http.Client{Timeout: time.Second * 10} 과 같이 타임아웃을 설정한 클라이언트를 생성하고,
	// client.Get("")을 호출해야 한다.
	res, err := http.Get("http://example.com")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func PostJsonTest() {
	const json = `{"name": "Matt", "age": 44}`
	in := strings.NewReader(json)

	res, err := http.Post("http://example.com/upload", "text/json", in)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(message))
}

// keyvalue store
type Store struct {
	sync.RWMutex
	m map[string]string
}

// sentinel error
var ErrorNoSuchKey = errors.New("no such key")

func (s *Store) Put(key string, value string) error {
	s.Lock()
	s.m[key] = value
	s.Unlock()
	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.RLock()
	value, ok := s.m[key]
	s.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (s *Store) Delete(key string) error {
	s.Lock()
	delete(s.m, key)
	s.Unlock()
	return nil
}

func NewStore() *Store {
	return &Store{m: make(map[string]string)}
}

var store *Store = NewStore()

// grpc study
type server struct {
	pb.UnimplementedKeyValueServer
}

func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := store.Get(r.Key)

	return &pb.GetResponse{Value: value}, err
}

func (s *server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received PUT key=%v, value=%v", r.Key, r.Value)

	err := store.Put(r.Key, r.Value)

	return &pb.PutResponse{}, err
}

func main() {
	// grpc 서버 생성
	s := grpc.NewServer()
	// KeyValueServer 등록
	pb.RegisterKeyValueServer(s, &server{})

	// 요청 수신을 위한 포트 오픈
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 수신 포트로부터 연결 수신
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
