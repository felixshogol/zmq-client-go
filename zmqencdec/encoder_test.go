package zmqencdec

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"gotest.tools/assert"
)

// ///////////////////////////////////////////
// Encode tests
func TestEncodeStartRequest(t *testing.T) {

	expect := "000a0001000004d10000000a"

	msg := &Message{
		Header: MsgHeader{
			Command: ZMQ_CMD_START,
		},
		StartRequest: MsgStartRequest{
			FlowId:          1233,
			MetricsInterval: 10,
		},
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	l := unsafe.Sizeof(msg.StartRequest) + lengthSize
	msg.Header.Length = uint16(l)

	encoder := &ZmqEncoder{}

	bytes, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, msg.Header.Length+uint16(cmdSize), uint16(len(bytes)), "\nThe two length should be the same.")
	assert.Equal(t, expect, hex.EncodeToString(bytes), "\nThe two array should be the same.")
}

func TestEncodeStopRequest(t *testing.T) {

	expect := "00060002000004d1"

	msg := &Message{
		Header: MsgHeader{
			Command: ZMQ_CMD_STOP,
		},
		StopRequest: MsgStopRequest{
			FlowId: 1233,
		},
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	l := unsafe.Sizeof(msg.StopRequest) + cmdSize
	msg.Header.Length = uint16(l)

	encoder := &ZmqEncoder{}

	bytes, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, msg.Header.Length+uint16(lengthSize), uint16(len(bytes)), "\nThe two length should be the same.")
	assert.Equal(t, expect, hex.EncodeToString(bytes), "\nThe two array should be the same.")
}

func TestEncodeAddTunnelsRequest(t *testing.T) {

	expect := "003a0004000004d10000000300000001000003e9000004d2000010e100000002000003ea0000162e0000223d00000003000003eb000023340000083d"

	msg := &Message{
		Header: MsgHeader{
			Command: ZMQ_CMD_ADD_TUNNELS,
		},
		AddTunnelRequest: MsgAddTunnelsRequest{
			FlowId: 1233,
			Tunnels: []Tunnel{
				{
					TeidIn:  1,
					TeidOut: 1001,
					UeIpV4:  1234,
					UpfIpV4: 4321,
				},
				{
					TeidIn:  2,
					TeidOut: 1002,
					UeIpV4:  5678,
					UpfIpV4: 8765,
				},
				{
					TeidIn:  3,
					TeidOut: 1003,
					UeIpV4:  9012,
					UpfIpV4: 2109,
				},
			},
		},
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	tunnels := len(msg.AddTunnelRequest.Tunnels)
	var dummy *Tunnel
	tunnelL := unsafe.Sizeof(*dummy)

	l := int(tunnelL)*tunnels + int(cmdSize) + 8 //flowid + tunnels
	msg.Header.Length = uint16(l)

	encoder := &ZmqEncoder{}

	bytes, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, msg.Header.Length+uint16(lengthSize), uint16(len(bytes)), "\nThe two length should be the same.")
	assert.Equal(t, expect, hex.EncodeToString(bytes), "\nThe two array should be the same.")
}

func TestEncodeDelTunnelsRequest(t *testing.T) {

	expect := "00160005000004d100000003000003e9000003ea000003eb"

	msg := &Message{
		Header: MsgHeader{
			Command: ZMQ_CMD_DEL_TUNNELS,
		},
		DelTunnelsRequest: MsgDelTunnelsRequest{
			FlowId: 1233,
			Teids: []uint32{
				1001, 1002, 1003,
			},
		},
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	teids := len(msg.DelTunnelsRequest.Teids)

	var dummy *uint32
	teidl := unsafe.Sizeof(*dummy)

	l := int(teidl)*teids + int(cmdSize) + 8 // flowid, tieds
	msg.Header.Length = uint16(l)
	encoder := &ZmqEncoder{}

	bytes, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, msg.Header.Length+uint16(lengthSize), uint16(len(bytes)), "\nThe two length should be the same.")
	assert.Equal(t, expect, hex.EncodeToString(bytes), "\nThe two array should be the same.")
}

func TestEncodeDelAllTunnelsRequest(t *testing.T) {

	expect := "00060006000004d1"

	msg := &Message{
		Header: MsgHeader{
			Command: ZMQ_CMD_DEL_ALL_TUNNELS,
		},
		DelAllTunnelsRequest: MsgDelAllTunnelsRequest{
			FlowId: 1233,
		},
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	var dummy *uint32
	flowIdSize := unsafe.Sizeof(*dummy)

	l := int(flowIdSize) + int(cmdSize)
	msg.Header.Length = uint16(l)
	encoder := &ZmqEncoder{}

	bytes, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, msg.Header.Length+uint16(lengthSize), uint16(len(bytes)), "\nThe two length should be the same.")
	assert.Equal(t, expect, hex.EncodeToString(bytes), "\nThe two array should be the same.")
}

func TestEncodeGetInfoRequest(t *testing.T) {

	expect := "00060007000004d1"

	msg := &Message{
		Header: MsgHeader{
			Command: ZMQ_CMD_GET_INFO,
		},
		GetInfoRequest: MsgGetInfoRequest{
			FlowId: 1233,
		},
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	var dummy *uint32
	flowIdSize := unsafe.Sizeof(*dummy)

	l := int(flowIdSize) + int(cmdSize)
	msg.Header.Length = uint16(l)
	encoder := &ZmqEncoder{}

	bytes, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, msg.Header.Length+uint16(lengthSize), uint16(len(bytes)), "\nThe two length should be the same.")
	assert.Equal(t, expect, hex.EncodeToString(bytes), "\nThe two array should be the same.")
}

// ///////////////////////////////////////////
// Decode tests
func TestDecodeStartResponse(t *testing.T) {
	str := "00110001000004d16c6f63616c3a3539303031"
	expect := &Message{
		Header: MsgHeader{
			Length:  uint16(17),
			Command: ZMQ_CMD_START,
		},
		StartResponse: MsgStartResponse{
			FlowId:    uint32(1233),
			Publisher: "local:59001",
		},
	}

	v := []byte("local:59001")
	fmt.Printf("publisher %s\n", hex.EncodeToString(v))

	encoder := &ZmqEncoder{}
	bytes, _ := hex.DecodeString(str)
	msg, err := encoder.Decode(bytes)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, expect.Header, msg.Header, "\nThe two ZMQ message header should be the same.")
	assert.Equal(t, expect.StartResponse, msg.StartResponse, "\nThe two ZMQ message header should be the same.")
}

func TestDecodeGetInfoResponse(t *testing.T) {
	str := "00110007000004d1646678702076312e31"
	expect := &Message{
		Header: MsgHeader{
			Length:  uint16(17),
			Command: ZMQ_CMD_GET_INFO,
		},
		GetInfoResponse: MsgGetInfoResponse{
			FlowId:  uint32(1233),
			Version: "dfxp v1.1",
		},
	}

	v := []byte("dfxp v1.1")
	fmt.Printf("version %s\n", hex.EncodeToString(v))

	encoder := &ZmqEncoder{}
	bytes, _ := hex.DecodeString(str)
	msg, err := encoder.Decode(bytes)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, expect.Header, msg.Header, "\nThe two ZMQ message header should be the same.")
	assert.Equal(t, expect.GetInfoResponse, msg.GetInfoResponse, "\nThe two ZMQ message header should be the same.")
}

func TestDecodeStopResponse(t *testing.T) {
	str := "00060002000004d1"
	expect := &Message{
		Header: MsgHeader{
			Length:  uint16(6),
			Command: ZMQ_CMD_STOP,
		},
		Response: MsgResponse{
			FlowId: uint32(1233),
		},
	}

	encoder := &ZmqEncoder{}
	bytes, _ := hex.DecodeString(str)
	msg, err := encoder.Decode(bytes)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, expect.Header, msg.Header, "\nThe two ZMQ message header should be the same.")
	assert.Equal(t, expect.Response, msg.Response, "\nThe two ZMQ message header should be the same.")
}

func TestDecodeDeleteTunnelsResponse(t *testing.T) {
	str := "000a0005000004d100000006"
	expect := &Message{
		Header: MsgHeader{
			Length:  uint16(10),
			Command: ZMQ_CMD_DEL_TUNNELS,
		},
		TunnelResponse: MsgTunnelResponse{
			FlowId:  uint32(1233),
			Tunnels: uint32(6),
		},
	}

	encoder := &ZmqEncoder{}
	bytes, _ := hex.DecodeString(str)
	msg, err := encoder.Decode(bytes)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, expect.Header, msg.Header, "\nThe two ZMQ message header should be the same.")
	assert.Equal(t, expect.TunnelResponse, msg.TunnelResponse, "\nThe two ZMQ message header should be the same.")
}

func TestDecodeDeleteAllTunnelsResponse(t *testing.T) {
	str := "000a0006000004d100000006"
	expect := &Message{
		Header: MsgHeader{
			Length:  uint16(10),
			Command: ZMQ_CMD_DEL_ALL_TUNNELS,
		},
		TunnelResponse: MsgTunnelResponse{
			FlowId:  uint32(1233),
			Tunnels: uint32(6),
		},
	}

	encoder := &ZmqEncoder{}
	bytes, _ := hex.DecodeString(str)
	msg, err := encoder.Decode(bytes)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, expect.Header, msg.Header, "\nThe two ZMQ message header should be the same.")
	assert.Equal(t, expect.TunnelResponse, msg.TunnelResponse, "\nThe two ZMQ message header should be the same.")
}

func TestDecodeAddTunnelssResponse(t *testing.T) {
	str := "000a0004000004d100000006"
	expect := &Message{
		Header: MsgHeader{
			Length:  uint16(10),
			Command: ZMQ_CMD_ADD_TUNNELS,
		},
		TunnelResponse: MsgTunnelResponse{
			FlowId:  uint32(1233),
			Tunnels: uint32(6),
		},
	}

	encoder := &ZmqEncoder{}
	bytes, _ := hex.DecodeString(str)
	msg, err := encoder.Decode(bytes)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}
	assert.Equal(t, expect.Header, msg.Header, "\nThe two ZMQ message header should be the same.")
	assert.Equal(t, expect.TunnelResponse, msg.TunnelResponse, "\nThe two ZMQ message header should be the same.")
}
