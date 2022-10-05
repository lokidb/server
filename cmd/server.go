package cmd

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/lokidb/engine"
	"github.com/lokidb/server/communication/grpc"
	"github.com/lokidb/server/communication/rest"
)

var (
	data_dir    = flag.String("storage_dir", "./", "The path where the data files will be created")
	cache_size  = flag.Int("cache_size", 50000, "Max number of keys to save in the RAM for fast reads")
	files_count = flag.Int("files_count", 10, "number of files to save the keys on (more files means better performance, max 1000)")

	rest_port = flag.Int("rest_port", 8080, "REST server port")
	res_host  = flag.String("rest_host", "127.0.0.1", "REST server host")
	run_rest  = flag.Bool("run_rest", true, "Serve REST API")

	grpc_port = flag.Int("grpc_port", 50051, "gRPC server port")
	grpc_host = flag.String("grpc_host", "127.0.0.1", "gRPC server host")
	run_grpc  = flag.Bool("run_grpc", true, "Serve gRPC API")
)

func Execute() {
	flag.Parse()
	engine := engine.New(*data_dir, *cache_size, *files_count)

	if !(*run_grpc || *run_rest) {
		fmt.Println("You must allow at least on of (REST, gRPC) to run.")
		os.Exit(1)
		return
	}

	var wg sync.WaitGroup

	if *run_rest {
		restServer := rest.NewServer(*res_host, *rest_port, &engine)

		wg.Add(1)
		go func() {
			restServer.Run()
			wg.Done()
		}()
	}

	if *run_grpc {
		grpcServer := grpc.NewServer(&engine, *grpc_host, *grpc_port)
		wg.Add(1)
		go func() {
			grpcServer.Run()
			wg.Done()
		}()
	}

	wg.Wait()
}
