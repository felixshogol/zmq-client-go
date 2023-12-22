package zmqencdec

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"encoding/hex"

	"github.com/golang/glog"
)

type ZmqEncoder struct {
}

// Encode - encode Messages
func (enc *ZmqEncoder) Encode(msg *Message) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("Message nil")
	}
	glog.Infof("Encode ZMQ message:%v", msg)

	switch msg.Header.Command {
	case (ZMQ_CMD_START):
		return enc.encodeStartRequest(msg)
	case ZMQ_CMD_STOP:
		return enc.encodeStopRequest(msg)
	case ZMQ_CMD_ADD_TUNNELS:
		return enc.encodeAddTunnelRequest(msg)
	case ZMQ_CMD_DEL_TUNNELS:
		return enc.encodeDelTunnelRequest(msg)
	case ZMQ_CMD_DEL_ALL_TUNNELS:
		return enc.encodeDelAllTunnelsRequest(msg)
	case ZMQ_CMD_GET_INFO:
		return enc.encodeGetInfoRequest(msg)

	default:
		return nil, fmt.Errorf("Wrong message command [%d]", msg.Header.Command)
	}
}

func (enc *ZmqEncoder) encodeStartRequest(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, msg.Header.Length)
	binary.Write(buffer, binary.BigEndian, msg.Header.Command)
	binary.Write(buffer, binary.BigEndian, msg.StartRequest.FlowId)
	binary.Write(buffer, binary.BigEndian, msg.StartRequest.MetricsInterval)

	return buffer.Bytes(), nil
}

func (enc *ZmqEncoder) encodeStopRequest(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, msg.Header.Length)
	binary.Write(buffer, binary.BigEndian, msg.Header.Command)
	binary.Write(buffer, binary.BigEndian, msg.StopRequest.FlowId)

	return buffer.Bytes(), nil
}

func (enc *ZmqEncoder) encodeAddTunnelRequest(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, msg.Header.Length)
	binary.Write(buffer, binary.BigEndian, msg.Header.Command)
	binary.Write(buffer, binary.BigEndian, msg.AddTunnelRequest.FlowId)
	binary.Write(buffer, binary.BigEndian, uint32(len(msg.AddTunnelRequest.Tunnels)))
	for _, tunnel := range msg.AddTunnelRequest.Tunnels {
		binary.Write(buffer, binary.BigEndian, tunnel.TeidIn)
		binary.Write(buffer, binary.BigEndian, tunnel.TeidOut)
		binary.Write(buffer, binary.BigEndian, tunnel.UeIpV4)
		binary.Write(buffer, binary.BigEndian, tunnel.SrvIpV4)
	}

	return buffer.Bytes(), nil
}

func (enc *ZmqEncoder) encodeDelTunnelRequest(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, msg.Header.Length)
	binary.Write(buffer, binary.BigEndian, msg.Header.Command)
	binary.Write(buffer, binary.BigEndian, msg.DelTunnelsRequest.FlowId)
	binary.Write(buffer, binary.BigEndian, uint32(len(msg.DelTunnelsRequest.Teids)))
	for _, teid := range msg.DelTunnelsRequest.Teids {
		binary.Write(buffer, binary.BigEndian, teid)
	}

	return buffer.Bytes(), nil
}

func (enc *ZmqEncoder) encodeDelAllTunnelsRequest(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, msg.Header.Length)
	binary.Write(buffer, binary.BigEndian, msg.Header.Command)
	binary.Write(buffer, binary.BigEndian, msg.DelAllTunnelsRequest.FlowId)
	binary.Write(buffer, binary.BigEndian, uint32(0))// tunnels number

	return buffer.Bytes(), nil
}

func (enc *ZmqEncoder) encodeGetInfoRequest(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, msg.Header.Length)
	binary.Write(buffer, binary.BigEndian, msg.Header.Command)
	binary.Write(buffer, binary.BigEndian, msg.GetInfoRequest.FlowId)

	return buffer.Bytes(), nil
}

// Decode - decode Messages
func (enc *ZmqEncoder) Decode(bytesArray []byte) (*Message, error) {

	if bytesArray == nil {
		return nil, fmt.Errorf("bytesArray nil")
	}
	glog.Infof("Decode ZMQ message:%v", hex.EncodeToString(bytesArray))

	msg := &Message{}

	buffer := bytes.NewBuffer(bytesArray)
	binary.Read(buffer, binary.BigEndian, &msg.Header.Length)
	binary.Read(buffer, binary.BigEndian, &msg.Header.Command)

	switch msg.Header.Command {
	case (ZMQ_CMD_START):
		return msg, enc.decodeStartResponse(buffer, msg)
	case ZMQ_CMD_STOP:
		return msg, enc.decodeStopResponse(buffer, msg)
	case ZMQ_CMD_ADD_TUNNELS:
		return msg, enc.decodeAddTunnelResponse(buffer, msg)
	case ZMQ_CMD_DEL_TUNNELS:
		return msg, enc.decodeDelTunnelResponse(buffer, msg)
	case ZMQ_CMD_DEL_ALL_TUNNELS:
		return msg, enc.decodeDelAllTEIDsResponse(buffer, msg)
	case ZMQ_CMD_GET_INFO:
		return msg, enc.decodeGetInfoResponse(buffer, msg)
	case ZMQ_CMD_ERROR:
		return msg, enc.decodeErrorResponse(buffer, msg)
	case ZMQ_CMD_MSG_ERROR:
		return msg, enc.decodeMsgErrorResponse(buffer, msg)
	default:
		return nil, fmt.Errorf("Wrong message command [%d]", msg.Header.Command)
	}
}

func (enc *ZmqEncoder) decodeStartResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.StartResponse.FlowId)
	//read publisher
	bytes := buffer.Bytes()
	glog.Infof("bytes:%s", hex.EncodeToString(bytes))
	msg.StartResponse.Publisher = string(bytes)
	return nil
}

func (enc *ZmqEncoder) decodeStopResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.Response.FlowId)
	return nil
}
func (enc *ZmqEncoder) decodeAddTunnelResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.TunnelResponse.FlowId)
	binary.Read(buffer, binary.BigEndian, &msg.TunnelResponse.Tunnels)
	return nil
}
func (enc *ZmqEncoder) decodeDelTunnelResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.TunnelResponse.FlowId)
	binary.Read(buffer, binary.BigEndian, &msg.TunnelResponse.Tunnels)
	return nil
}

func (enc *ZmqEncoder) decodeDelAllTEIDsResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.TunnelResponse.FlowId)
	binary.Read(buffer, binary.BigEndian, &msg.TunnelResponse.Tunnels)
	return nil

}
func (enc *ZmqEncoder) decodeGetInfoResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.GetInfoResponse.FlowId)
	//read Version
	bytes := buffer.Bytes()
	glog.Infof("bytes:%s", hex.EncodeToString(bytes))
	msg.GetInfoResponse.Version = string(bytes)
	return nil

}

func (enc *ZmqEncoder) decodeErrorResponse(buffer *bytes.Buffer, msg *Message) error {
	bytes := buffer.Bytes()
	glog.Infof("bytes:%s", hex.EncodeToString(bytes))

	//read Error
	msg.MsgErrorResponse.Error = string(bytes)
	return nil
}

func (enc *ZmqEncoder) decodeMsgErrorResponse(buffer *bytes.Buffer, msg *Message) error {
	binary.Read(buffer, binary.BigEndian, &msg.MsgErrorResponse.FlowId)
	//read Error
	bytes := buffer.Bytes()
	glog.Infof("bytes:%s", hex.EncodeToString(bytes))
	msg.MsgErrorResponse.Error = string(bytes)
	return nil
}
