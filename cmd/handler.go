package cmd

import (
	"github.com/SouvickSarkar20/goredis/persistence"
	"github.com/SouvickSarkar20/goredis/resp"
	"github.com/SouvickSarkar20/goredis/store"
	"strings"
)

type HandlerFunc func(*resp.Writer, *store.Store, resp.Value) error

var routes = map[string]HandlerFunc{
	"PING":      handlePing,
	"SET":       handleSet,
	"GET":       handleGet,
	"DEL":       handleDel,
	"LPUSH":     handleLPush,
	"LPOP":      handleLPop,
	"HGET":      handleHGet,
	"HSET":      handleHAdd,
	"HDEL":      handleHDel,
	"SMEMBERS":  handleSMembers,
	"SADD":      handleSAdd,
	"SISMEMBER": handleSIsMember,
	"SREM":      handleSRem,
}

var aofLogger *persistence.AOF

var mutatingCommands = map[string]struct{}{
	"SET":   {},
	"DEL":   {},
	"LPUSH": {},
	"LPOP":  {},
	"HSET":  {},
	"HDEL":  {},
	"SADD":  {},
	"SREM":  {},
}

// SetAOF injects the AOF logger into the command handler.
func SetAOF(aof *persistence.AOF) {
	aofLogger = aof
}

func Handle(writer *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) == 0 {
		return writer.WriteError("ERR empty command")
	}

	command := strings.ToUpper(input.Array[0].Str)

	handler, exists := routes[command]
	if !exists {
		return writer.WriteError("ERR unknown command " + command)
	}

	// Execute handler
	if err := handler(writer, db, input); err != nil {
		return err
	}

	// Log mutating commands to AOF (after successful execution)
	if _, shouldLog := mutatingCommands[command]; shouldLog && aofLogger != nil {
		args := make([]string, 0, len(input.Array))
		for _, v := range input.Array {
			args = append(args, v.Str)
		}
		if err := aofLogger.AppendCommand(args); err != nil {
			return writer.WriteError("ERR AOF append failed")
		}
	}

	return nil
}

func handlePing(w *resp.Writer, db *store.Store, args resp.Value) error {
	return w.WriteSimpleString("PONG")
}
