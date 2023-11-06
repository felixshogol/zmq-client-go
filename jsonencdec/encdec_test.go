package jsonencdec

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"zmqclient/zmqencdec"

	"gotest.tools/assert"
)

var jsonEncoder JsonEncoder

// ///////////////////////////////////////////
// Encode tests
func TestJsonEncodeStartRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 1 
		},
		"StartRequest": {
			"FlowId":1234,
			"MetricsInterval": 5
		}
	}

	`
	expMsg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(10),
			Command: zmqencdec.ZMQ_CMD_START,
		},
		StartRequest: zmqencdec.MsgStartRequest{
			FlowId:          1234,
			MetricsInterval: 5,
		},
	}
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}

	assert.Equal(t, expMsg.Header, msg.Header, "\nThe two message header should be the same.")
	assert.Equal(t, expMsg.StartRequest, msg.StartRequest, "\nThe two messages should be the same.")

}

func TestJsonEncodeStopRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 1 
		},
		"StopRequest": {
			"FlowId":1234
		}
	}

	`
	expMsg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(10),
			Command: zmqencdec.ZMQ_CMD_START,
		},
		StopRequest: zmqencdec.MsgStopRequest{
			FlowId: 1234,
		},
	}
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}

	assert.Equal(t, expMsg.Header, msg.Header, "\nThe two message header should be the same.")
	assert.Equal(t, expMsg.StopRequest, msg.StopRequest, "\nThe two messages should be the same.")

}

func TestJsonEncodeAddTunnelstRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 1 
		},
		"AddTunnelRequest": {
			"FlowId":1234,
			"Tunnels" : [
               {
				"TeidIn":1,
				"TeidOut":1001,
				"UeIpV4":1234,
				"UpfIpV4": 4321
			   },
			   {
				"TeidIn":2,
				"TeidOut":1002,
				"UeIpV4":5678,
				"UpfIpV4": 8765
			   },
			   {
				"TeidIn":3,
				"TeidOut":1003,
				"UeIpV4":9012,
				"UpfIpV4": 2109
			   }
			] 
		}
	}

	`
	expMsg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(10),
			Command: zmqencdec.ZMQ_CMD_START,
		},
		AddTunnelRequest: zmqencdec.MsgAddTunnelsRequest{
			FlowId: 1234,
			Tunnels: []zmqencdec.Tunnel{
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
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}

	assert.DeepEqual(t, expMsg.Header, msg.Header)
	assert.DeepEqual(t, expMsg.AddTunnelRequest, msg.AddTunnelRequest)

}

func TestJsonEncodeDeleteTunnelstRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 5 
		},
		"DelTunnelsRequest": {
			"FlowId":1234,
			"Teids" : [
				1001,1002,1003
			] 
		}
	}

	`
	expMsg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(10),
			Command: zmqencdec.ZMQ_CMD_DEL_TUNNELS,
		},
		DelTunnelsRequest: zmqencdec.MsgDelTunnelsRequest{
			FlowId: 1234,
			Teids: []uint32{
				1001, 1002, 1003,
			},
		},
	}
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}

	assert.DeepEqual(t, expMsg.Header, msg.Header)
	assert.DeepEqual(t, expMsg.DelTunnelsRequest, msg.DelTunnelsRequest)
}

func TestJsonEncodeDeleteAllTunnelstRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 6 
		},
		"DelTunnelsRequest": {
			"FlowId":1234
			
		}
	}

	`
	expMsg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(10),
			Command: zmqencdec.ZMQ_CMD_DEL_ALL_TUNNELS,
		},
		DelTunnelsRequest: zmqencdec.MsgDelTunnelsRequest{
			FlowId: 1234,
		},
	}
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}

	assert.DeepEqual(t, expMsg.Header, msg.Header)
	assert.DeepEqual(t, expMsg.DelAllTunnelsRequest, msg.DelAllTunnelsRequest)
}

func TestJsonEncodeGetInfoRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 7 
		},
		"GetInfoRequest": {
			"FlowId":1234
		}
	}
	`
	expMsg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(10),
			Command: zmqencdec.ZMQ_CMD_GET_INFO,
		},
		GetInfoRequest: zmqencdec.MsgGetInfoRequest{
			FlowId: 1234,
		},
	}
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Encode failed. Err:%v", err)
	}

	assert.DeepEqual(t, expMsg.Header, msg.Header)
	assert.DeepEqual(t, expMsg.GetInfoRequest, msg.GetInfoRequest)
}

/////////////////////////////////////////////
// Decode tests

func TestJsonDecodeStartResponse(t *testing.T) {
	expJson := `"{'Header':{'Length':17,'Command':1},'StartRequest':{'FlowId':0,'MetricsInterval':0},'StopRequest':{'FlowId':0},'StartResponse':{'FlowId':1233,'Publisher':'local:59001'},'Response':{'FlowId':0},'AddTunnelRequest':{'FlowId':0,'Tunnels':null},'DelTunnelsRequest':{'FlowId':0,'Teids':null},'DelAllTunnelsRequest':{'FlowId':0},'TunnelResponse':{'FlowId':0,'Tunnels':0},'GetInfoRequest':{'FlowId':0},'GetInfoResponse':{'FlowId':0,'Version':''}}"`
	msg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(17),
			Command: zmqencdec.ZMQ_CMD_START,
		},
		StartResponse: zmqencdec.MsgStartResponse{
			FlowId:    uint32(1233),
			Publisher: "local:59001",
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Decode failed:%s", err)
		return
	}
	escapedMarshal := strings.ReplaceAll(string(jsonMsg), `"`, `'`)
	jsonToString := fmt.Sprintf(`"%s"`, escapedMarshal)
	assert.Equal(t, expJson, jsonToString, "\nThe two jsons string should be the same.")
}

func TestJsonDecodeResponse(t *testing.T) {
	expJson := `"{'Header':{'Length':17,'Command':2},'StartRequest':{'FlowId':0,'MetricsInterval':0},'StopRequest':{'FlowId':0},'StartResponse':{'FlowId':0,'Publisher':''},'Response':{'FlowId':1233},'AddTunnelRequest':{'FlowId':0,'Tunnels':null},'DelTunnelsRequest':{'FlowId':0,'Teids':null},'DelAllTunnelsRequest':{'FlowId':0},'TunnelResponse':{'FlowId':0,'Tunnels':0},'GetInfoRequest':{'FlowId':0},'GetInfoResponse':{'FlowId':0,'Version':''}}"`
	msg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(17),
			Command: zmqencdec.ZMQ_CMD_STOP,
		},
		Response: zmqencdec.MsgResponse{
			FlowId:    uint32(1233),
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Decode failed:%s", err)
		return
	}
	escapedMarshal := strings.ReplaceAll(string(jsonMsg), `"`, `'`)
	jsonToString := fmt.Sprintf(`"%s"`, escapedMarshal)
	assert.Equal(t, expJson, jsonToString, "\nThe two jsons string should be the same.")
}

func TestJsonDecodeTunnelResponse(t *testing.T) {
	expJson := `"{'Header':{'Length':17,'Command':4},'StartRequest':{'FlowId':0,'MetricsInterval':0},'StopRequest':{'FlowId':0},'StartResponse':{'FlowId':0,'Publisher':''},'Response':{'FlowId':0},'AddTunnelRequest':{'FlowId':0,'Tunnels':null},'DelTunnelsRequest':{'FlowId':0,'Teids':null},'DelAllTunnelsRequest':{'FlowId':0},'TunnelResponse':{'FlowId':1233,'Tunnels':10},'GetInfoRequest':{'FlowId':0},'GetInfoResponse':{'FlowId':0,'Version':''}}"`
	msg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(17),
			Command: zmqencdec.ZMQ_CMD_ADD_TUNNELS,
		},
		TunnelResponse: zmqencdec.MsgTunnelResponse{
			FlowId:    uint32(1233),
			Tunnels: 10,
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Decode failed:%s", err)
		return
	}
	escapedMarshal := strings.ReplaceAll(string(jsonMsg), `"`, `'`)
	jsonToString := fmt.Sprintf(`"%s"`, escapedMarshal)
	assert.Equal(t, expJson, jsonToString, "\nThe two jsons string should be the same.")
}

func TestJsonDecodeGetInfoResponse(t *testing.T) {
	expJson := `"{'Header':{'Length':17,'Command':7},'StartRequest':{'FlowId':0,'MetricsInterval':0},'StopRequest':{'FlowId':0},'StartResponse':{'FlowId':0,'Publisher':''},'Response':{'FlowId':0},'AddTunnelRequest':{'FlowId':0,'Tunnels':null},'DelTunnelsRequest':{'FlowId':0,'Teids':null},'DelAllTunnelsRequest':{'FlowId':0},'TunnelResponse':{'FlowId':0,'Tunnels':0},'GetInfoRequest':{'FlowId':0},'GetInfoResponse':{'FlowId':1233,'Version':'dfxp v1.1'}}"`
	msg := &zmqencdec.Message{
		Header: zmqencdec.MsgHeader{
			Length:  uint16(17),
			Command: zmqencdec.ZMQ_CMD_GET_INFO,
		},
		GetInfoResponse: zmqencdec.MsgGetInfoResponse{
			FlowId:    uint32(1233),
			Version: "dfxp v1.1",
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Decode failed:%s", err)
		return
	}
	escapedMarshal := strings.ReplaceAll(string(jsonMsg), `"`, `'`)
	jsonToString := fmt.Sprintf(`"%s"`, escapedMarshal)
	assert.Equal(t, expJson, jsonToString, "\nThe two jsons string should be the same.")
}


