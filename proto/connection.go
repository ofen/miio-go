package proto

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

const (
	// DefaultReadBufferSize is maximum size of read buffer
	// (could be overrided)
	DefaultReadBufferSize = 4096
	proto                 = "udp"
)

// Conn is protocol connection
// Implements net.Conn interface https://pkg.go.dev/net#Conn
type Conn struct {
	Token          []byte
	ReadBufferSize int
	requestID      int
	conn           net.Conn
	keys           deviceKeys
}

// Dial connects to the device with given token.
// Token should be 32 characters in lenght (16 bytes).
//
// Example:
//   Dial("192.168.0.3:54321", "a0b1c2d3e4f5a0b1c2d3e4f5a0b1c2d3")
func Dial(addr string, token string) (Conn, error) {
	conn, err := net.Dial(proto, addr)
	if err != nil {
		return Conn{}, err
	}

	t, err := hex.DecodeString(token)
	if err != nil {
		return Conn{}, err
	}

	tokenLenght := len(t)
	if tokenLenght != 16 {
		return Conn{}, fmt.Errorf("expected 16 bytes token, got %d", tokenLenght)
	}

	keys := newDeviceKeys(t)

	return Conn{
		Token:          t,
		ReadBufferSize: DefaultReadBufferSize,
		conn:           conn,
		keys:           keys,
	}, nil
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (int, error) {
	buff := make([]byte, c.ReadBufferSize)
	n, err := c.conn.Read(buff)
	if err != nil {
		return n, err
	}

	copy(b, c.keys.decrypt(buff[32:]))

	return n - 32, nil
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (int, error) {
	h, err := c.handshake()
	if err != nil {
		return 0, err
	}

	encrypted := c.keys.encrypt(b)
	req := prepareRequest(c.Token, h, encrypted)

	return c.conn.Write(req)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

type handshakeResponse struct {
	deviceID    uint32
	serverStamp uint32
}

func parseHandshakeResponse(b []byte) (*handshakeResponse, error) {
	lenght := len(b)
	if lenght != 32 {
		return nil, fmt.Errorf("unable to parse handshake, expected 32 bytes, got %d", lenght)
	}
	deviceID := binary.BigEndian.Uint32(b[8:])
	serverStamp := binary.BigEndian.Uint32(b[12:])
	return &handshakeResponse{deviceID, serverStamp}, nil
}

func (c *Conn) handshake() (*handshakeResponse, error) {
	req := []byte{
		0x21, 0x31, 0x00, 0x20, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}

	_, err := c.conn.Write(req)
	if err != nil {
		return nil, err
	}

	resp := make([]byte, c.ReadBufferSize)
	n, err := c.conn.Read(resp)
	if err != nil {
		return nil, err
	}

	return parseHandshakeResponse(resp[:n])
}

func prepareRequest(token []byte, handshake *handshakeResponse, requestBody []byte) []byte {
	header := [32]byte{0x21, 0x31}
	binary.BigEndian.PutUint16(header[2:], uint16(32+len(requestBody)))
	binary.BigEndian.PutUint32(header[8:], handshake.deviceID)
	binary.BigEndian.PutUint32(header[12:], handshake.serverStamp)
	checksum := md5(header[:16], token, requestBody)
	copy(header[16:], checksum)
	return append(header[:], requestBody...)
}
