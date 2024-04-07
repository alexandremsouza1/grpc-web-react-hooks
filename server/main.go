package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"time"

	messagerpb "github.com/okmttdhr/grpc-web-react-hooks/proto/messenger"// Importe o pacote gerado para os protos do serviço Messenger
	uploadpb "github.com/okmttdhr/grpc-web-react-hooks/proto/upload"       // Importe o pacote gerado para os protos do serviço Upload

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9090"
)

// Implementação do servidor gRPC
type server struct {
	messagerpb.UnimplementedMessengerServer
	requestsMessager []*messagerpb.MessageRequest
	uploadpb.UnimplementedUploadServer
	requestUpload []*uploadpb.UploadRequest
}

// Métodos para o serviço Messenger
func (s *server) GetMessages(_ *empty.Empty, stream messagerpb.Messenger_GetMessagesServer) error {
	for _, r := range s.requestsMessager {
		if err := stream.Send(&messagerpb.MessageResponse{Message: r.GetMessage()}); err != nil {
			return err
		}
	}

	previousCount := len(s.requestsMessager)

	for {
		currentCount := len(s.requestsMessager)
		if previousCount < currentCount && currentCount > 0 {
			r := s.requestsMessager[currentCount-1]
			log.Printf("Sent: %v", r.GetMessage())
			if err := stream.Send(&messagerpb.MessageResponse{Message: r.GetMessage()}); err != nil {
				return err
			}
		}
		previousCount = currentCount
	}
}

func (s *server) CreateMessage(ctx context.Context, r *messagerpb.MessageRequest) (*messagerpb.MessageResponse, error) {
	log.Printf("Received: %v", r.GetMessage())
	newR := &messagerpb.MessageRequest{Message: r.GetMessage() + ": " + time.Now().Format("2006-01-02 15:04:05")}
	s.requestsMessager = append(s.requestsMessager, newR)
	return &messagerpb.MessageResponse{Message: r.GetMessage()}, nil
}

func (s *server) UploadFile(ctx context.Context, r *uploadpb.UploadRequest) (*uploadpb.UploadResponse, error) {

	// Exemplo simples: apenas imprime as informações recebidas
	fmt.Printf("Recebido UploadRequest para o arquivo: %s, tamanho: %d bytes\n", r.GetFileName(), r.GetFileSize())

	// Lógica para processar os chunks do arquivo, se necessário
	for _, chunk := range r.GetChunks() {
		// Processar cada chunk do arquivo
		fmt.Printf("Chunk %d recebido: %d bytes\n", chunk.GetChunkNumber(), len(chunk.GetData()))
	}

	// Retornar uma mensagem de sucesso
	resp := &uploadpb.UploadResponse{
		Message: "Arquivo recebido com sucesso!",
	}
	return resp, nil
}


func main() {
	// Cria um listener TCP na porta especificada
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Cria um novo servidor gRPC
	s := grpc.NewServer()

	// Registra a implementação do serviço Messenger no servidor gRPC
	messagerpb.RegisterMessengerServer(s, &server{})

	// Registra a implementação do serviço Upload no servidor gRPC
	uploadpb.RegisterUploadServer(s, &server{})

	// Registra a reflexão do serviço no servidor gRPC (para uso com ferramentas como o gRPCurl)
	reflection.Register(s)

	// Inicia o servidor gRPC
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
