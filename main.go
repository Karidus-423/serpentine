package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"serpentine/analysis"
	"serpentine/lsp"
	"serpentine/rpc"
)

func main() {
	logger := getLogger("/home/kennett/personal/serpentine/log.txt")
	logger.Println("Started Logger.")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Error produced.")
		}
		handleMessage(logger, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
	logger.Printf("Recived message with method %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("From InitializeRequest: Unable to parse %s", err)
		}
		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		//Reply...
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Sent Initialize response successfully.")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf(
				"From DidOpenTextDocumentNotification: Unable to parse %s", err)
		}
		logger.Printf("Opened: %s \n Contents: %s",
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text)

		state.OpenDocument(
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf(
				"From TextDocumentDidChangeNotification: Unable to parse %s", err)
		}

		logger.Printf("Opened: %s \n Contents: %s",
			request.Params.TextDocument.URI,
			request.Params.ContentChanges)

		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("No valid filename given")
	}

	return log.New(logfile, "[serpentine] ", log.Ldate|log.Ltime|log.Lshortfile)
}
