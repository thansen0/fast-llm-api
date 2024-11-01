package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"
    "bufio"

    "go.uber.org/zap"
    "gopkg.in/natefinch/lumberjack.v2"
    "go.uber.org/zap/zapcore"

    "google.golang.org/grpc"
    pb "github.com/thansen0/fast-llm-api/protos"

	"github.com/spf13/viper"
    c "github.com/thansen0/fast-llm-api/config"
)

var logger *zap.Logger
var conf c.Configurations
var userMap map[string]bool

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
        Compress:   false, // Compress rotated files
    })

    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable timestamp
    encoder := zapcore.NewJSONEncoder(encoderConfig)

    core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
    return zap.New(core)
}

func verify_user(api_key string) bool {
    var verified_user bool = userMap[api_key]

    if verified_user {
        return verified_user
    } else {
        // check postgres table
        var cur_uuid_validity bool = false

        // if true, add user to table
        if cur_uuid_validity {
            userMap[api_key] = cur_uuid_validity
            // sync to existing file
            writeUUIDMap(userMap)
        }
    }

    return verified_user

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

func initUUIDMap() map[string]bool {
    fmt.Println("File: "+conf.Verify.UuidFile)
    fmt.Println(conf.Server.Port)
    file, err := os.Open(conf.Verify.UuidFile)
    if err != nil {
        log.Fatalf("failed to open file: %v", err)
    }
    defer file.Close()

    var uuidMap map[string]bool = make(map[string]bool)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        uuid := scanner.Text()
        uuidMap[uuid] = true // Store UUID in the map
    }

    // Check for errors encountered during scanning
    if err := scanner.Err(); err != nil {
        log.Fatalf("error reading file: %v", err)
    }

    // Now uuidMap contains all UUIDs from the file
    // fmt.Println("Loaded UUIDs:", uuidMap)

    return uuidMap
}

func writeUUIDMap(uuidMap map[string]bool) error {
    file, err := os.Create(conf.Verify.UuidFile)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    // Initialize a writer to write each UUID to the file
    writer := bufio.NewWriter(file)
    for uuid := range uuidMap {
        _, err := writer.WriteString(uuid + "\n")
        if err != nil {
            return fmt.Errorf("failed to write UUID to file: %w", err)
        }
    }

    // Ensure all data is flushed to the file
    if err := writer.Flush(); err != nil {
        return fmt.Errorf("failed to flush data to file: %w", err)
    }

    fmt.Println("Closing UUID user map")
    return nil
}

func main() {
    // load config file with viper
	viper.SetConfigName("config")
	viper.AddConfigPath("config/")
	// viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

    // Set undefined variables
	// viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}

/*    fmt.Printf("Server Port: %d\n", conf.Server.Port)
    fmt.Printf("Verification UUID File: %s\n", conf.Verify.UuidFile)
    fmt.Printf("Verification URL: %s\n", conf.Verify.VerifyUrl) */

    // configure zap logger
    logger = initLogger()
    defer logger.Sync()

    // initialize user map
    userMap = initUUIDMap()

    // Set up a listener on port 50051
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Server.Port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Create a new gRPC server
    grpcServer := grpc.NewServer()

    // Register the service with the server
    pb.RegisterAskLLMQuestionServer(grpcServer, &server{})

    log.Println(fmt.Sprintf("gRPC server is running on port %d", conf.Server.Port))
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

