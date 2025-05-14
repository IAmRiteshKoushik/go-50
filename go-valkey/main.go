package main

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

func main() {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()
	// Select DB1
	client.B().Select().Index(1)
	client.Do(ctx, client.B().Set().Key("test2").Value("data").Nx().Build()).Error()
	hm, _ := client.Do(ctx, client.B().Get().Key("test2").Build()).ToString()
	m, _ := client.Do(ctx, client.B().Keys().Pattern("*").Build()).AsStrSlice()

	fmt.Println(hm)
	fmt.Println(m)
}
