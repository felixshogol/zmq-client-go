package zmqclient

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
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


func TestAddTunnelsRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 4 
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

	ip := int2ip(1234)
	t.Logf("ip1 :%s",ip.String())
	ip = int2ip(4321)
	t.Logf("ip2 :%s",ip.String())

	ip = int2ip(5678)
	t.Logf("ip3 :%s",ip.String())
	ip = int2ip(8765)
	t.Logf("ip4 :%s",ip.String())

	ip = int2ip(9012)
	t.Logf("ip5 :%s",ip.String())
	ip = int2ip(2109)
	t.Logf("ip6 :%s",ip.String())



	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}
    
	
	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	lengthSize := reflect.TypeOf(msg.Header.Length).Size()

	tunnels := len(msg.AddTunnelRequest.Tunnels)
	var dummy *zmqencdec.Tunnel
	tunnelL := unsafe.Sizeof(*dummy)

	l := int(tunnelL)*tunnels + int(cmdSize) + 8 //flowid + tunnels
	msg.Header.Length = uint16(l)
	

	request, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("ZMQ Encode failed. Err:%v", err)
	}

	t.Logf("Tunnels request:%s\n", hex.EncodeToString(request))

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

	t.Logf("Start response: cmd:%d len:%d flowid:%d ", msg.Header.Command,msg.Header.Length,msg.TunnelResponse.FlowId)

	assert.Equal(t, uint16(lengthSize), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.TunnelResponse.FlowId, "\nThe two FlowId should be the same.")
}

func IntToIPv4(i int) {
	panic("unimplemented")
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

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}