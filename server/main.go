package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	messagerpb "github.com/okmttdhr/grpc-web-react-hooks/proto/messenger" // Importe o pacote gerado para os protos do serviço Messenger
	uploadpb "github.com/okmttdhr/grpc-web-react-hooks/proto/upload"      // Importe o pacote gerado para os protos do serviço Upload

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
	fileName := r.GetFileName()
	chunks := r.GetChunks()
	isLastChunk := r.GetIsLastChunk()

	// cria uma pasta com o nome do arquivo sem a extensão
	tempPath := fileName[:len(fileName)-4]
	err := os.MkdirAll(tempPath, os.ModePerm)
	if err != nil {
			log.Printf("Erro ao criar pasta: %v", err)
			return nil, err
	}

	// Combinar todos os chunks em um único slice de bytes
	var data []byte
	var chunkNum int32
	for _, chunk := range chunks {
			data = append(data, chunk.GetData()...)
			chunkNum = chunk.GetChunkNumber()
	}

	// Criar o nome do arquivo completo com base no número do chunk
	completeFileName := fmt.Sprintf("%s/%d_%s", tempPath, chunkNum, fileName)

	// Salvar os dados em um arquivo
	err = saveToFile(completeFileName, data)
	if err != nil {
			log.Printf("Erro ao salvar arquivo: %v", err)
			return nil, err
	}

	// Se for o último chunk, consolidar os arquivos e limpar a pasta temporária
	if isLastChunk {
			consolidatedFileName := fileName // Nome do arquivo consolidado
			err := consolidateFiles(tempPath, consolidatedFileName)
			if err != nil {
					log.Printf("Erro ao consolidar arquivos: %v", err)
					return nil, err
			}

			// Remover pasta temporária
			err = os.RemoveAll(tempPath)
			if err != nil {
					log.Printf("Erro ao remover pasta temporária: %v", err)
					return nil, err
			}
	}

	// Retornar uma mensagem de sucesso
	resp := &uploadpb.UploadResponse{
			Message: "Arquivo recebido e salvo com sucesso!",
	}
	return resp, nil
}

func saveToFile(fileName string, data []byte) error {
	// Criar um novo arquivo no servidor
	file, err := os.Create(fileName)
	if err != nil {
			return err
	}
	defer file.Close()

	// Escrever os dados no arquivo
	_, err = file.Write(data)
	if err != nil {
			return err
	}

	return nil
}

func consolidateFiles(tempPath, fileName string) error {
	// Abrir arquivo consolidado para escrita
	consolidatedFile, err := os.Create(fileName)
	if err != nil {
			return err
	}
	defer consolidatedFile.Close()

	// Listar todos os arquivos na pasta temporária
	fileInfos, err := ioutil.ReadDir(tempPath)
	if err != nil {
			return err
	}

	// Concatenar o conteúdo de cada arquivo no arquivo consolidado
	for _, fileInfo := range fileInfos {
			// Ignorar diretórios, se houver
			if fileInfo.IsDir() {
					continue
			}

			// Abrir arquivo para leitura
			filePath := filepath.Join(tempPath, fileInfo.Name())
			fileData, err := ioutil.ReadFile(filePath)
			if err != nil {
					return err
			}

			// Escrever conteúdo do arquivo no arquivo consolidado
			if _, err := consolidatedFile.Write(fileData); err != nil {
					return err
			}
	}

	return nil
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
