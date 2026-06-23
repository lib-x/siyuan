package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lib-x/siyuan"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <file>", os.Args[0])
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	name := filepath.Base(os.Args[1])
	result, err := client.Assets.Upload(ctx, siyuan.UploadAssetsRequest{
		AssetsDirPath: "/assets/",
		Files: []siyuan.UploadAssetFile{
			{Name: name, Reader: file},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.SuccMap[name])
}

func newClient() (*siyuan.Client, error) {
	opts := []siyuan.Option{
		siyuan.WithToken(os.Getenv("SIYUAN_TOKEN")),
	}
	if endpoint := os.Getenv("SIYUAN_ENDPOINT"); endpoint != "" {
		opts = append(opts, siyuan.WithEndpoint(endpoint))
	}
	return siyuan.New(opts...)
}
