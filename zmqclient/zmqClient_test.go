package zmqclient

import (
	"encoding/hex"
	"strconv"
	"testing"
	"zmqclient/jsonencdec"
	"zmqclient/zmqencdec"

	"gotest.tools/assert"
)

var jsonEncoder jsonencdec.JsonEncoder
var encoder zmqencdec.ZmqEncoder
var options ClientOptions

const (
	Host            = "192.168.1.227"
	ServerPort      = 5555
	PublisherPort   = 5557
	MetricsInterval = 5
)

func TestStartRequest(t *testing.T) {

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
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}

	request, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("ZMQ Encode failed. Err:%v", err)
	}

	client := buildZmqClient(t)

	t.Logf("Connect to server\n")
	err = client.Connect(options.To)
	if err != nil {
		t.Fatalf("Connect failed. Err:%v", err)
	}
	t.Logf("SendAndReceiveWithTimeout\n")
	response, err := client.SendAndReceiveWithTimeout(request, options.To)
	if err != nil {
		t.Fatalf("SendAndReceiveWithTimeout failed. Err:%v", err)
	}

	t.Logf("Response:%s\n", hex.EncodeToString(response))

	msg, err = encoder.Decode(response)

	if err != nil {
		t.Fatalf("Decode response failed. Err:%v", err)
	}

	t.Logf("Start response: cmd:%d length:%d", msg.Header.Command,msg.Header.Length)

	t.Logf("Start response: flowid:%d publisher:%s", msg.StartResponse.FlowId, msg.StartResponse.Publisher)
	expLength := len(msg.StartResponse.Publisher) + 4 + 2 // flowid + h.cmd

	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.StartResponse.FlowId, "\nThe two FlowId should be the same.")
	assert.Equal(t, "tcp://"+Host+":"+strconv.Itoa(PublisherPort), msg.StartResponse.Publisher, "\nThe two FlowId should be the same.")
}

func TestStopRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 6,
		  "Command": 2 
		},
		"StopRequest": {
			"FlowId":1234
		}
	}
	`
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}

	request, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("ZMQ Encode failed. Err:%v", err)
	}

	client := buildZmqClient(t)

	t.Logf("Connect to server\n")
	err = client.Connect(options.To)
	if err != nil {
		t.Fatalf("Connect failed. Err:%v", err)
	}
	t.Logf("SendAndReceiveWithTimeout\n")
	response, err := client.SendAndReceiveWithTimeout(request, options.To)
	if err != nil {
		t.Fatalf("SendAndReceiveWithTimeout failed. Err:%v", err)
	}

	t.Logf("Response:%s\n", hex.EncodeToString(response))

	msg, err = encoder.Decode(response)

	if err != nil {
		t.Fatalf("Decode response failed. Err:%v", err)
	}

	t.Logf("Start response: cmd:%d len:%d flowid:%d ", msg.Header.Command,msg.Header.Length,msg.Response.FlowId)
	expLength :=  4 + 2 // flowid + h.cmd

	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.Response.FlowId, "\nThe two FlowId should be the same.")
}

// ////////////////////////////////////////////////
// Local functions
// ////////////////////////////////////////////////
func buildZmqClient(t *testing.T) *ZmqClient {

	options.Host = Host
	options.Port = ServerPort
	options.To = MetricsInterval // seconds
	client := NewZmqClient(&options)

	return client
}
