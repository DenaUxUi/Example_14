package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var (
    ec2Client      *ec2.Client
    lastInstanceId string
)

type Response struct {
    Message string `json:"system-response"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    name := r.URL.Query().Get("name")
    if name == "" {
        http.Error(w, "Missing 'name' parameter in query", http.StatusBadRequest)
        return
    }

    response := Response{
        Message: fmt.Sprintf("Hello, %s!", name),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleCreateInstance(w http.ResponseWriter, r *http.Request) {
    ami := os.Getenv("AMI_ID")
    keyName := os.Getenv("KEY_NAME")
    // instanceType := "t3.micro"

    input := &ec2.RunInstancesInput{
        ImageId:      aws.String(ami),
        InstanceType: ec2Types.InstanceType("t3.micro"),
        MinCount:     aws.Int32(1),
        MaxCount:     aws.Int32(1),
        KeyName:      aws.String(keyName),
    }

    result, err := ec2Client.RunInstances(context.TODO(), input)
    if err != nil || len(result.Instances) == 0 {
	log.Printf("RunInstances failed: %v\n", err)
        http.Error(w, "Instance creation error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    lastInstanceId = *result.Instances[0].InstanceId
    log.Println("Created new instance:", lastInstanceId)

    resp := map[string]string{
        "instanceId": lastInstanceId,
        "state":      string(result.Instances[0].State.Name),
    }
    json.NewEncoder(w).Encode(resp)
}

func handleTerminateInstance(w http.ResponseWriter, r *http.Request) {
    if lastInstanceId == "" {
        http.Error(w, "We not found active instance", http.StatusBadRequest)
        return
    }

    input := &ec2.TerminateInstancesInput{
        InstanceIds: []string{lastInstanceId},
    }

    result, err := ec2Client.TerminateInstances(context.TODO(), input)
    if err != nil {
        http.Error(w, "Error at termination: "+err.Error(), http.StatusInternalServerError)
        return
    }

    state := result.TerminatingInstances[0].CurrentState.Name
    log.Println("Instance terminated:", lastInstanceId)

    resp := map[string]string{
        "terminated": lastInstanceId,
        "state":      string(state),
    }
    json.NewEncoder(w).Encode(resp)

    lastInstanceId = ""
}

func main() {
    // 1. AWS SDK init
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(os.Getenv("AWS_REGION")),
    )
    if err != nil {
        log.Fatalf("Error loading AWS config: %v", err)
    }

    ec2Client = ec2.NewFromConfig(cfg)

    // 2. Handlers
    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/api/hello", helloHandler)
    http.HandleFunc("/api/aws/vm/create", handleCreateInstance)
    http.HandleFunc("/api/aws/vm/terminate", handleTerminateInstance)

    // 3. Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Println("Server started at http://localhost:" + port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

