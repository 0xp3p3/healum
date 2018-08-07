package main

import (
	"fmt"
	example "github.com/micro/examples/server/proto/example"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-os/trace"
	"github.com/micro/go-plugins/trace/zipkin"
	"golang.org/x/net/context"
	"time"
)

// publishes a message
func pub() {
	msg := client.NewPublication("topic.go.micro.srv.example", &example.Message{
		Say: "This is a publication",
	})

	// create context with metadata
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	// publish message
	if err := client.Publish(ctx, msg); err != nil {
		fmt.Println("pub err: ", err)
		return
	}

	fmt.Printf("Published: %v\n", msg)
}

func call(i int) {
	// Create new request to service go.micro.srv.example, method Example.Call
	req := client.NewRequest("go.micro.srv.example", "Example.Call", &example.Request{
		Name: "John",
	})

	// create context with metadata
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	rsp := &example.Response{}

	// Call service
	if err := client.Call(ctx, req, rsp); err != nil {
		fmt.Println("call err: ", err, rsp)
		return
	}

	fmt.Println("Call:", i, "rsp:", rsp.Msg)
}

func stream(i int) {
	// Create new request to service go.micro.srv.example, method Example.Call
	// Request can be empty as its actually ignored and merely used to call the handler
	req := client.NewRequest("go.micro.srv.example", "Example.Stream", &example.StreamingRequest{})

	stream, err := client.Stream(context.Background(), req)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if err := stream.Send(&example.StreamingRequest{Count: int64(i)}); err != nil {
		fmt.Println("err:", err)
		return
	}
	for stream.Error() == nil {
		rsp := &example.StreamingResponse{}
		err := stream.Recv(rsp)
		if err != nil {
			fmt.Println("recv err", err)
			break
		}
		fmt.Println("Stream: rsp:", rsp.Count)
	}

	if stream.Error() != nil {
		fmt.Println("stream err:", err)
		return
	}

	if err := stream.Close(); err != nil {
		fmt.Println("stream close err:", err)
	}
}

func pingPong(i int) {
	// Create new request to service go.micro.srv.example, method Example.Call
	// Request can be empty as its actually ignored and merely used to call the handler
	req := client.NewRequest("go.micro.srv.example", "Example.PingPong", &example.StreamingRequest{})

	stream, err := client.Stream(context.Background(), req)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	for j := 0; j < i; j++ {
		if err := stream.Send(&example.Ping{Stroke: int64(j + 1)}); err != nil {
			fmt.Println("err:", err)
			return
		}
		rsp := &example.Pong{}
		err := stream.Recv(rsp)
		if err != nil {
			fmt.Println("recv err", err)
			break
		}
		fmt.Printf("Sent ping %v got pong %v\n", j+1, rsp.Stroke)
	}

	if stream.Error() != nil {
		fmt.Println("stream err:", err)
		return
	}

	if err := stream.Close(); err != nil {
		fmt.Println("stream close err:", err)
	}
}

func main() {
	cmd.Init()

	t := zipkin.NewTrace(
		trace.Collectors("192.168.99.100:9092"),
	)
	defer t.Close()

	client.DefaultClient = client.NewClient(
		client.Wrap(
			trace.ClientWrapper(t, nil),
		),
	)

	fmt.Println("\n--- Traced Call example ---\n")
	for i := 0; i < 10; i++ {
		call(i)
	}

	/*
		fmt.Println("\n--- Streamer example ---\n")
		stream(10)

		fmt.Println("\n--- Ping Pong example ---\n")
		pingPong(10)

		fmt.Println("\n--- Publisher example ---\n")
		pub()
	*/
	<-time.After(time.Second * 10)
}
