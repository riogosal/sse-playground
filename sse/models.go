package sse

type EventName string

const (
	InitEvent   EventName = "init"
	PingEvent   EventName = "ping"
	Information EventName = "information"
)

func Event(name EventName) []byte {
	return append([]byte("event: "), []byte(name)...)
}

func Data(data []byte) []byte {
	return append([]byte("data: "), data...)
}
