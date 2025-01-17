package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
	//Params...
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  int    `json:"id"`
	//Results...
	//Errors...
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
	//Params...
}
