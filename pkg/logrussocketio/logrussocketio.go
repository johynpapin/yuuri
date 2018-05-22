/*
Package logrussocketio is a hook for logrus to send logs via socket.io.
*/
package logrussocketio

import (
	"encoding/json"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type message struct {
	level   logrus.Level  `json:"level"`
	time    time.Time     `json:"time"`
	data    logrus.Fields `json:"data"`
	message string        `json:"message"`
}

// Hook to send logs via socket.io
type Hook struct {
	Socket socketio.Socket
}

// NewHook create a hook to be added to an instance of logger
func NewHook(socket socketio.Socket) *Hook {
	return &Hook{socket}
}

// Fire is used by logrus when a message is logged
func (hook *Hook) Fire(entry *logrus.Entry) error {
	return emitToSocket(hook.Socket, message{entry.Level, entry.Time, entry.Data, entry.Message})
}

// Levels is used by logrus to get the levels avaible for this hook
func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func emitToSocket(socket socketio.Socket, message interface{}) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return socket.Emit("log", messageJSON)
}
