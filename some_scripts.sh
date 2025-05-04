# установка протобаф
sudo apt install -y protobuf-compiler
# плагин для генерации кода
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# Добавьте пути к установленным плагинам в переменную окружения PATH
export PATH="$PATH:$(go env GOPATH)/bin"
# зависимости для работы с gRPC и Protobuf:
go get google.golang.org/protobuf
go get google.golang.org/grpc

# генерация прото дял сервера
svs@svs-pc:~/work/imageRrocessing$ protoc --go_out=./server --go-grpc_out=./server ./service.proto