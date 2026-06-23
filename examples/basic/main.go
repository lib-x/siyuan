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

	version, err := client.System.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SiYuan version:", version)

	notebooks, err := client.Notebooks.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, notebook := range notebooks {
		fmt.Printf("%s\t%s\tclosed=%v\n", notebook.ID, notebook.Name, notebook.Closed)
	}
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
