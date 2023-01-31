package migration

import (
	"fmt"
	"io/ioutil"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/lxc/lxd/shared"
)

// ProtoRecv gets a protobuf message from a websocket.
func ProtoRecv(ws *websocket.Conn, msg proto.Message) error {
	if ws == nil {
		return fmt.Errorf("Empty websocket connection")
	}

	mt, r, err := ws.NextReader()
	if err != nil {
		return err
	}

	if mt != websocket.BinaryMessage {
		return fmt.Errorf("Only binary messages allowed")
	}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return err
	}

	return nil
}

// ProtoSend sends a protobuf message over a websocket.
func ProtoSend(ws *websocket.Conn, msg proto.Message) error {
	if ws == nil {
		return fmt.Errorf("Empty websocket connection")
	}

	w, err := ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	defer func() { _ = w.Close() }()

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	err = shared.WriteAll(w, data)
	if err != nil {
		return err
	}

	return w.Close()
}

// ProtoSendControl sends a migration control message over a websocket.
func ProtoSendControl(ws *websocket.Conn, err error) {
	message := ""
	if err != nil {
		message = err.Error()
	}

	msg := MigrationControl{
		Success: proto.Bool(err == nil),
		Message: proto.String(message),
	}

	_ = ProtoSend(ws, &msg)
}
