package persistence

import (
	"github.com/SouvickSarkar20/goredis/resp"
	"os"
	"sync"
)

type FsyncMode string

const (
	FsyncAlways   FsyncMode = "always"
	FsyncEverySec FsyncMode = "everysec" // treated as buffered sync-on-close for now
	FsyncNo       FsyncMode = "no"
)

type AOF struct {
	mu   sync.Mutex
	file *os.File
	mode FsyncMode
}

func OpenAOF(path string, mode FsyncMode) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &AOF{
		file: f,
		mode: mode,
	}, nil
}

func (a *AOF) AppendCommand(args []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	w := resp.NewWriter(a.file)

	if err := w.WriteArray(len(args)); err != nil {
		return err
	}
	for _, arg := range args {
		if err := w.WriteBulkString(arg); err != nil {
			return err
		}
	}

	if a.mode == FsyncAlways {
		return a.file.Sync()
	}
	return nil
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.mode == FsyncEverySec || a.mode == FsyncAlways {
		_ = a.file.Sync()
	}
	return a.file.Close()
}
