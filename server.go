package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4443"
	CONN_TYPE = "tcp"
)

type quoteRequest struct {
	stock  string
	userID string
}

type quoteResponse struct {
	quote     float32
	stock     string
	userID    string
	timestamp int64
	cyrptokey string
}

func main() {
	// If we don't provide a seed for rand then it hehaves as if
	// we ran Seed(1). It's not safe to run this in concurrent
	// code so I'm doing it here!
	rand.Seed(time.Now().Unix())

	// Accept incoming connections
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on", CONN_HOST+":"+CONN_PORT)

	// Send active connections off for handing
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(2)
		}

		//logs an incoming message
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		// Use concurrent goroutines to serve connections
		go generateQuote(conn)
	}
}

func generateQuote(conn net.Conn) {
	// Read incoming data into buffer
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		conn.Close()
		return
	}

	req, err := parseReq(buff)
	if err != nil {
		// bail on the connection if it has a malformed request
		fmt.Println("Error parsing request:", err.Error())
		conn.Close()
		return
	}

	// use request to generate values for response
	resp := makeResp(req)

	// Send back the quote
	conn.Write([]byte(resp.ToCSVString()))

	// Don't need this anymore
	conn.Close()
}

func parseReq(buff []byte) (quoteRequest, error) {
	// convert to a string for easier processing
	buffSize := bytes.Index(buff, []byte{0})
	if buffSize <= 1 {
		// Probably a request with no body
		return quoteRequest{}, errors.New("Missing request?")
	}

	// Break out the request into individual arguments
	request := string(buff[:buffSize-1])
	requestParts := strings.Split(request, ",")
	if len(requestParts) != 2 {
		return quoteRequest{}, errors.New("Wrong number of arguments")
	}

	// Naïvely assume arguments are in right order and format.
	return quoteRequest{stock: requestParts[0], userID: requestParts[1]}, nil
}

func makeResp(req quoteRequest) quoteResponse {
	// Only use the first 3 char of a stock
	var truncatedStock string
	if stockLen := len(req.stock); stockLen < 3 {
		truncatedStock = req.stock[:stockLen]
	} else {
		truncatedStock = req.stock[:3]
	}

	// Send back current server time
	nowUnix := time.Now().Unix()

	// The cryptokey will be base64(stock + user + now)
	// FIXME: The output from this doesn't change as often as
	// it should. Lots of eHl6dG9t77+9. Should probably learn what
	// this does at some point.
	seed := req.stock + req.userID + string(nowUnix)
	cryptokey := base64.StdEncoding.EncodeToString([]byte(seed))

	return quoteResponse{
		quote:     1000 * rand.Float32(),
		stock:     strings.ToUpper(truncatedStock),
		userID:    req.userID,
		timestamp: nowUnix,
		cyrptokey: cryptokey,
	}
}

func (resp *quoteResponse) ToCSVString() string {
	s := []string{
		fmt.Sprintf("%.2f", resp.quote),
		resp.stock,
		resp.userID,
		fmt.Sprintf("%d", resp.timestamp),
		resp.cyrptokey,
	}

	return strings.Join(s, ",")
}
