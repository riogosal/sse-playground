package sse

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func SetSSEHeaders(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Disable buffering for nginx
	c.Header("Transfer-Encoding", "chunked")
}

func OpenStream(c *gin.Context, signal <-chan struct{}, err_ch <-chan error) {
	SetSSEHeaders(c)
	if err := SignalInit(c); err != nil {
		fmt.Println("Error streaming data", err)
		return
	}
	for {
		select {
		case <-signal:
			if err := SignalNewData(c); err != nil {
				fmt.Println("Error streaming data", err)
				return
			}

		// Pinger
		case time := <-time.After(15 * time.Second):
			if err := Ping(c, time); err != nil {
				fmt.Println("Error streaming time", err)
				return
			}

		case <-c.Writer.CloseNotify():
			fmt.Println("Client closed connection")
			return

		case err := <-err_ch:
			fmt.Println("Error streaming data", err)
			return
		}

	}
}

func OpenStringStream(c *gin.Context, signal <-chan string, err_ch <-chan error) {
	SetSSEHeaders(c)
	if err := SignalInit(c); err != nil {
		fmt.Println("Error streaming data", err)
		return
	}
	for {
		select {
		case data := <-signal:
			if err := StreamData(c, Event(Information), Data([]byte(data))); err != nil {
				fmt.Println("Error streaming data", err)
				return
			}

		// Pinger
		case time := <-time.After(15 * time.Second):
			if err := Ping(c, time); err != nil {
				fmt.Println("Error streaming time", err)
				return
			}

		case <-c.Writer.CloseNotify():
			fmt.Println("Client closed connection")
			return

		case err := <-err_ch:
			fmt.Println("Error streaming data", err)
			return
		}

	}
}

func SignalInit(c *gin.Context) error {
	bytes, err := jsoniter.Marshal("init")
	if err != nil {
		return err
	}
	return StreamData(c, Event(InitEvent), Data(bytes))
}

func Ping(c *gin.Context, timestamp time.Time) error {
	bytes, err := jsoniter.Marshal(timestamp.UnixMilli())
	if err != nil {
		return err
	}

	return StreamData(c, Event(PingEvent), Data(bytes))
}

func SignalData(c *gin.Context, data any) error {
	bytes, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}
	return StreamData(c, Event(Information), Data(bytes))
}

func SignalNewData(c *gin.Context) error {
	return SignalData(c, "new data")
}

func StreamData(c *gin.Context, data ...[]byte) error {
	finalBytes := []byte{}
	for _, bytes := range data {
		finalBytes = append(finalBytes, bytes...)
		finalBytes = append(finalBytes, []byte("\n")...)
	}

	// append newline to bytes
	finalBytes = append(finalBytes, []byte("\n\n")...)
	if _, err := c.Writer.Write(finalBytes); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}
