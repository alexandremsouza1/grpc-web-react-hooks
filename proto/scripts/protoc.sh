#!/bin/sh

set -xe

protoc --version

# Loop through all .proto files in the proto directory
for PROTO_FILE in proto/*.proto; do
  if [ -f "$PROTO_FILE" ]; then
    # Extrair o nome do arquivo sem a extensão
    FILENAME=$(basename "$PROTO_FILE" .proto)
    ONLY_FILE=$($PROTO_FILE | cut -d'/' -f2)
    # Determinar as pastas de saída com base no nome do arquivo proto
    SERVER_OUTPUT_DIR="server/${ONLY_FILE}"
    CLIENT_OUTPUT_DIR="client/src/${FILENAME}"

    # Deletar as pastas de saída, se existirem
    rm -rf "${CLIENT_OUTPUT_DIR}"

    # Criar as pastas de saída, se não existirem
    mkdir -p "${SERVER_OUTPUT_DIR}"
    mkdir -p "${CLIENT_OUTPUT_DIR}"

    # Gerar os arquivos
    protoc --proto_path=proto "${PROTO_FILE}" \
      --go_out="${SERVER_OUTPUT_DIR}" \
      --go-grpc_out="${SERVER_OUTPUT_DIR}" \
      --js_out=import_style=commonjs:"${CLIENT_OUTPUT_DIR}" \
      --grpc-web_out=import_style=typescript,mode=grpcwebtext:"${CLIENT_OUTPUT_DIR}"
  fi
done