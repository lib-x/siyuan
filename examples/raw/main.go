package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib-x/siyuan"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	var currentTime int64
	if err := client.Raw.Post(ctx, "/api/system/currentTime", nil, &currentTime); err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentTime)
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
