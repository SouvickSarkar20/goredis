package resp

import (
	"io"
	"strconv"
)

// Writer is responsible for formatting Go values into RESP (Redis Serialization Protocol)
// and writing them to the given connection over the network.
type Writer struct {
	writer io.Writer
}

// NewWriter creates a new Writer wrapping an io.Writer (like a net.Conn).
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// WriteSimpleString writes a RESP Simple String.
// In RESP, a Simple String starts with '+' and ends with '\r\n'.
// For example, if s is "OK", this writes "+OK\r\n".
func (w *Writer) WriteSimpleString(s string) error {
	// First convert the Go string back into a byte array
	// then we just append the prefix and the expected suffix
	bytes := []byte("+" + s + "\r\n")

	// Write these bytes to the underlying connection
	_, err := w.writer.Write(bytes)
	return err
}

// WriteError writes a RESP Error.
// In RESP, an Error starts with '-' and ends with '\r\n'.
// For example, if err is "ERR unknown command", this writes "-ERR unknown command\r\n".
func (w *Writer) WriteError(errMessage string) error {
	bytes := []byte("-" + errMessage + "\r\n")
	_, err := w.writer.Write(bytes)
	return err
}

// WriteInteger writes a RESP Integer.
// In RESP, an Integer starts with ':' and ends with '\r\n'.
func (w *Writer) WriteInteger(i int64) error {
	s := strconv.FormatInt(i, 10)
	bytes := []byte(":" + s + "\r\n")
	_, err := w.writer.Write(bytes)
	return err
}

// WriteBulkString writes a RESP Bulk String.
// It starts with '$' followed by the length of the string, a '\r\n',
// the actual string data, and a final '\r\n'.
func (w *Writer) WriteBulkString(s string) error {
	length := strconv.Itoa(len(s))
	bytes := []byte("$" + length + "\r\n" + s + "\r\n")
	_, err := w.writer.Write(bytes)
	return err
}

// WriteBulkStringNil writes a RESP Nil Bulk String.
// Used when a key does not exist. It sends "$-1\r\n".
func (w *Writer) WriteBulkStringNil() error {
	bytes := []byte("$-1\r\n")
	_, err := w.writer.Write(bytes)
	return err
}

// WriteArray writes a RESP Array prefix.
// It starts with '*' followed by the number of elements in the array, and '\r\n'.
func (w *Writer) WriteArray(length int) error {
	s := strconv.Itoa(length)
	bytes := []byte("*" + s + "\r\n")
	_, err := w.writer.Write(bytes)
	return err
}
