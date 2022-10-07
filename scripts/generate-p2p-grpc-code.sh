
# Run from lokidb root

cd pkg/p2p_grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative spec.proto