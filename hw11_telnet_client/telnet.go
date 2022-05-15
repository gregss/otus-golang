package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	ctx     context.Context
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *Client) Connect() (err error) {
	c.conn, err = (&net.Dialer{Timeout: c.timeout}).DialContext(c.ctx, "tcp", c.addr)
	return
}

func (c *Client) Close() (err error) {
	if c.conn == nil {
		return
	}
	err = c.conn.Close()
	return
}

func (c *Client) Send() (err error) {
	b := make([]byte, 100)
	n, err := c.in.Read(b)
	if err != nil {
		return
	}
	_, err = c.conn.Write(b[:n])
	return
}

func (c *Client) Receive() (err error) {
	reader := bufio.NewReader(c.conn)
	b := make([]byte, 100)
	n, _ := reader.Read(b)
	_, err = c.out.Write(b[:n])

	return
}

func NewTelnetClient(
	ctx context.Context,
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
) TelnetClient {
	return &Client{ctx, address, timeout, in, out, nil}
}
