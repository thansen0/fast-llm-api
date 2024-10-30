package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "github.com/thansen0/localllmapiserver/protos"
)

func main() {
    // Set up a connection to the server
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewAskLLMQuestionClient(conn)

    // Contact the server and print out the response
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    req := &pb.LLMInit{
        ApiKey: "b150b0f9-235e-4f87-91db-d2b45da98a68",
        Prompt: "What is the capital of France?",
    }
    res, err := client.PromptLLM(ctx, req)
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("LLM answer: %s", res.Answer)
}

