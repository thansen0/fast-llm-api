package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "go.uber.org/zap"
    "gopkg.in/natefinch/lumberjack.v2"
    "go.uber.org/zap/zapcore"

    "google.golang.org/grpc"
    pb "github.com/thansen0/localllmapiserver/protos"
)

var logger *zap.Logger

// server is used to implement AskLLMQuestionServer
type server struct {
    pb.UnimplementedAskLLMQuestionServer
}

func initLogger() *zap.Logger {
    writeSyncer := zapcore.AddSync(&lumberjack.Logger{
        Filename:   "./logs/apirequests.log",
        MaxSize:    10, // Megabytes before rotating
        MaxBackups: 3,  // Number of old files to retain
        MaxAge:     28, // Days to retain old log files
        Compress:   true, // Compress rotated files
    })

    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable timestamp
    encoder := zapcore.NewJSONEncoder(encoderConfig)

    core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
    return zap.New(core)
}

func verify_user(api_key string) bool {
    return true
}

func savePromptAnswer(prompt string, api_key string, answer string) {
    logger.Info(fmt.Sprintf("%s, %s, %s", api_key, prompt, answer))
}

func generate_llm_response(prompt string) string {
    fmt.Printf("Received prompt: %s\n", prompt)

    return "fake llama response"
}

// PromptLLM implements AskLLMQuestion.PromptLLM
func (s *server) PromptLLM(ctx context.Context, request *pb.LLMInit) (*pb.LLMInference, error) {
    answerChannel := make(chan string)

    // generate response from llm 
    go func() {
        answerChannel <- generate_llm_response(request.Prompt)
    }()

    if verify_user(request.ApiKey) {
        answer := <-answerChannel
        go savePromptAnswer(request.Prompt, request.ApiKey, answer)

        return &pb.LLMInference{Answer: answer}, nil
    } else {
        answer := "401 Error, Check your API key"
        return &pb.LLMInference{Answer: answer}, nil
    }
}

func main() {
    // configure zap logger
    logger = initLogger()
    defer logger.Sync()

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

