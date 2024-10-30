package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/thansen0/localllmapiserver/protos"
)

func generate_llm_response(prompt string) string {
    fmt.Printf("Received prompt: %s\n", prompt)

    return "fake llama response"
}

func is_user(api_key string) bool {
    return true
}

func savePromptAnswer(prompt string, api_key string, answer string) {
    fmt.Println("Saving to log or DB")
}

// server is used to implement AskLLMQuestionServer
type server struct {
    pb.UnimplementedAskLLMQuestionServer
}

// PromptLLM implements AskLLMQuestion.PromptLLM
func (s *server) PromptLLM(ctx context.Context, request *pb.LLMInit) (*pb.LLMInference, error) {
    answerChannel := make(chan string)

    // generate response from llm 
    go func() {
        answerChannel <- generate_llm_response(request.Prompt)
    }()

    // check whether user is legit
    //      increment if true

    // prepare logging struct

    answer := <-answerChannel
    go savePromptAnswer(request.Prompt, request.ApiKey, answer)

    return &pb.LLMInference{Answer: answer}, nil
}

func main() {
    // Set up a listener on port 50051
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Create a new gRPC server
    grpcServer := grpc.NewServer()

    // Register the service with the server
    pb.RegisterAskLLMQuestionServer(grpcServer, &server{})

    log.Println("gRPC server is running on port 50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

