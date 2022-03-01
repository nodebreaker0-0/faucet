package main

import (
	//"encoding/json"
	"fmt"

	"github.com/joho/godotenv"
	//"github.com/tendermint/tmlibs/bech32"

	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var chain string
var amountATOM string

var key string
var node string
var publicUrl string

//var already []string

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, "=", value)
		return value
	} else {
		log.Fatal("Error loading environment variable: ", key)
		return ""
	}
}

func main() {
	err := godotenv.Load(".env_cosmos")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chain = getEnv("FAUCET_CHAIN")
	amountATOM = getEnv("FAUCET_AMOUNT_ATOM")

	key = getEnv("FAUCET_KEY")
	node = getEnv("FAUCET_NODE")
	publicUrl = getEnv("FAUCET_PUBLIC_URL")

	http.HandleFunc("/", getCoinsHandler)

	if err := http.ListenAndServe(publicUrl, nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func goExecute(command string) (cmd *exec.Cmd, pipeIn io.WriteCloser, pipeOut io.ReadCloser) {
	cmd = getCmd(command)
	pipeIn, _ = cmd.StdinPipe()
	pipeOut, _ = cmd.StdoutPipe()
	go cmd.Start()
	time.Sleep(2 * time.Second)
	return cmd, pipeIn, pipeOut
}

func getCmd(command string) *exec.Cmd {
	// split command into command and args
	split := strings.Split(command, " ")

	var cmd *exec.Cmd
	if len(split) == 1 {
		cmd = exec.Command(split[0])
	} else {
		cmd = exec.Command(split[0], split[1:]...)
	}

	return cmd
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getCoinsHandler(w http.ResponseWriter, request *http.Request) {
	enableCors(&w)
	query := request.URL.Query()
	address := query.Get("address")

	//for _, addres := range already {
	//	if address == addres {
	//		fmt.Fprintf(w, "You can only make 1 faucet request per account.")
	//		return
	//	}
	//}
	//already = append(already, address)

	sendFaucet := fmt.Sprintf("gaiad tx bank send %v %v %v --chain-id=%v -y --home /root/.gaia --node %v --keyring-backend test",
		key, address, amountATOM, chain, node)
	fmt.Println(sendFaucet)
	fmt.Println(time.Now().UTC().Format(time.RFC3339), address, "[1]")
	goExecute(sendFaucet)
	fmt.Fprintf(w, "Your faucet request has been processed successfully. Please check your wallet :)")
}
