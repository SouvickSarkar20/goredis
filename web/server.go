package web

import (
	"bytes"
	"encoding/json"
	"github.com/krishsinghhura/go-redis/cmd"
	"github.com/krishsinghhura/go-redis/resp"
	"github.com/krishsinghhura/go-redis/store"
	"net/http"
	"strings"
)

type commandRequest struct {
	Command string `json:"command"`
}

type commandResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func StartHTTP(addr string, db *store.Store) error {
	mux := http.NewServeMux()

	// /api/command — registered first, more specific
	mux.HandleFunc("/api/command", func(w http.ResponseWriter, r *http.Request) {
		handleCommand(w, r, db)
	})

	// / — catch-all, serves everything in web/dist/
	mux.Handle("/", http.FileServer(http.Dir("web/dist")))

	return http.ListenAndServe(addr, mux)
}

func servePlayground(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/playground.html")
}

func handleCommand(w http.ResponseWriter, r *http.Request, db *store.Store) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req commandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, commandResponse{Error: "invalid request body"})
		return
	}

	parts := strings.Fields(req.Command)
	if len(parts) == 0 {
		writeJSON(w, commandResponse{Error: "empty command"})
		return
	}

	// Build a fake RESP Value — same structure the TCP parser produces
	args := make([]resp.Value, len(parts))
	for i, p := range parts {
		args[i] = resp.Value{Typ: "BulkString", Str: p}
	}
	value := resp.Value{Typ: "Array", Array: args}

	// Write RESP response into a buffer (not a TCP connection)
	var buf bytes.Buffer
	respWriter := resp.NewWriter(&buf)

	// Execute the command — same as the TCP path
	if err := cmd.Handle(respWriter, db, value); err != nil {
		writeJSON(w, commandResponse{Error: err.Error()})
		return
	}

	// Parse the RESP bytes back into a Go value for JSON
	result := parseRESPToJSON(buf.Bytes())

	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, result)
}

func parseRESPToJSON(data []byte) commandResponse {
	if len(data) == 0 {
		return commandResponse{Result: nil}
	}

	switch data[0] {
	case '+': // Simple String: "+OK\r\n"
		return commandResponse{Result: strings.TrimRight(string(data[1:]), "\r\n")}
	case '-': // Error: "-ERR message\r\n"
		return commandResponse{Error: strings.TrimRight(string(data[1:]), "\r\n")}
	case ':': // Integer: ":42\r\n"
		numStr := strings.TrimRight(string(data[1:]), "\r\n")
		var n int64
		// Simple parse — production code would use strconv
		for _, c := range numStr {
			n = n*10 + int64(c-'0')
		}
		return commandResponse{Result: n}
	case '$': // Bulk String or Nil
		if data[1] == '-' { // $-1\r\n = nil
			return commandResponse{Result: nil}
		}
		// Parse the string value after the length line
		lines := bytes.SplitN(data, []byte("\r\n"), 3)
		if len(lines) >= 2 {
			return commandResponse{Result: string(lines[1])}
		}
	case '*': // Array
		p := resp.NewParser(bytes.NewReader(data))
		v, err := p.ParseOne()
		if err != nil {
			return commandResponse{Error: err.Error()}
		}
		result := make([]string, 0, len(v.Array))
		for _, item := range v.Array {
			result = append(result, item.Str)
		}
		return commandResponse{Result: result}
	}

	return commandResponse{Result: string(data)}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
