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
)

var (
    ec2Client      *ec2.Client
    lastInstanceId string
)


type Response struct {
    Message string `json:"system-response"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Разрешаем CORS
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

func main() {
    // Отдача статики из папки static при запросе к /
    http.Handle("/", http.FileServer(http.Dir("./static")))
    
    // API обработчик
    http.HandleFunc("/api/hello", helloHandler)

    fmt.Println("Server started at http://localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }

    cfg, err := config.LoadDefaultConfig(context.TODO(),
    	config.WithRegion(os.Getenv("AWS_REGION")),
    )
    if err != nil {
	log.Fatalf("Error of cfg AWS %v", err)
    }

    ec2Client = ec2.NewFromConfig(cfg)

    http.HandlerFunc("/api/aws/vm/create", handleCreateInstance)
    http.HandlerFunc("/api/aws/vm.terminate", handleTerminateInstance)

}

func handleCreateInstance(w http.ResponceWriter, r *http.Request){
    ami := os.Getenv("AMI_ID")
    keyName := os.Getenv("KEY_NAME")
    instanceType := "t3.micro"

    input := &ec2.RunInstanceInput{
	ImageId: aws.String(ami),
	InstanceType: ec2.InstanceType(instanceType),
	MinCount: aws.Int32(1),
	MaxCount: aws.Int32(1),
	KeyName: aws.String(keyName),
    } 

    result, err := ec2Client.RunInstance(context.TODO(), input)
    if err != nil || len(result.Instances) == 0 {
	http.Error(w, "Instance creation error: "+err.Error(), http.StatusInternalServerError)
	return
    }

    lastInstanceId = *result.Instance[0].InstanceId
    log.Println("Created new instance; ", lastInstanceId)

    resp:=map[string]string{
	"instanceId": lastInstanceId,
	"state":      string(result.Instances[0].State.Name),
    }
    json.NewEncoder(w).Encode(resp)
}

func handleTerminateInstance(w http.ResponseWriter, r *http.request){
    if lastInstanceid == "" {
	http.Error(w, "We not found active instance", http.StatusBadRequest)
	return
    }

    input := &ec2.TerminateInstanceInput{
	InstanceIds: []string{lastInstanceId},
    }

    result, err := ec2Client.TerminateInstances(context.TODO(), input)
    if err != nil {
	http.Error(w, "Error at termination: "+err.Error(), http.StatusInternalServerError)
	return
    }

    state := result.Terminatinginstance[0].CurrentState.Name
    log.Println("Instance is over: ", lastInstanceId)

    resp := map[string]string{
	"terminated": lastInstanceId,
	"state":      string(state),
    }
    json.NewEncoder(w).Encode(resp)

    lastInstanceId = ""
}

