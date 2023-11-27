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
	Host            = "192.168.1.173" //"192.168.1.227"
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

	t.Logf("Start response: cmd:%d length:%d", msg.Header.Command, msg.Header.Length)
	if msg.Header.Command == zmqencdec.ZMQ_CMD_ERROR {
		t.Fatalf("Received error:%s", msg.ErrorResponse.Error)
	}
	if msg.Header.Command == zmqencdec.ZMQ_CMD_MSG_ERROR {
		t.Fatalf("Received message error. flowid:%d error:%s", msg.MsgErrorResponse.FlowId, msg.MsgErrorResponse.Error)
	}

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

	t.Logf("Start response: cmd:%d len:%d", msg.Header.Command, msg.Header.Length)
	if msg.Header.Command == zmqencdec.ZMQ_CMD_ERROR {
		t.Fatalf("Received error:%s", msg.ErrorResponse.Error)
	}
	if msg.Header.Command == zmqencdec.ZMQ_CMD_MSG_ERROR {
		t.Fatalf("Received message error. flowid:%d error:%s", msg.MsgErrorResponse.FlowId, msg.MsgErrorResponse.Error)
	}

	expLength := 4 + 2 // flowid + + h.cmd

	assert.Equal(t, uint32(1234), msg.Response.FlowId, "\nThe two FlowId should be the same.")
	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")

}

func TestAddTunnelsRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 4 
		},
		"AddJsonTunnelRequest": {
			"FlowId":1234,
			"JsonTunnels" : [
               {
				"TeidIn":1,
				"TeidOut":1001,
				"UeIpV4":"10.10.10.1",
				"UpfIpV4": "12.12.12.1"
			   },
			   {
				"TeidIn":2,
				"TeidOut":1002,
				"UeIpV4":"10.10.10.2",
				"UpfIpV4": "12.12.12.1"
			   },
			   {
				"TeidIn":3,
				"TeidOut":1003,
				"UeIpV4":"10.10.10.3",
				"UpfIpV4": "12.12.12.1"
			   }
			] 
		}
	}
	`

	ip := int2ip(0xc0a801af)
	t.Logf("ip1 :%s", ip.String())
	ip = int2ip(4321)
	t.Logf("ip2 :%s", ip.String())

	ip = int2ip(5678)
	t.Logf("ip3 :%s", ip.String())
	ip = int2ip(8765)
	t.Logf("ip4 :%s", ip.String())

	ip = int2ip(9012)
	t.Logf("ip5 :%s", ip.String())
	ip = int2ip(2109)
	t.Logf("ip6 :%s", ip.String())

	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()

	tunnels := len(msg.AddJsonTunnelRequest.JsonTunnels)

	tunnelL := calculateTunnelsLen(t, msg.AddJsonTunnelRequest.JsonTunnels)
	t.Logf("tunnelL1:%d", tunnelL)
	l := int(tunnelL)*tunnels + int(cmdSize) + 8 //flowid + tunnels
	msg.Header.Length = uint16(l)

	err = encodeJsonAddTunnelsMessage(t,msg)

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
	t.Logf("Start response: cmd:%d len:%d flowid:%d ", msg.Header.Command, msg.Header.Length, msg.TunnelResponse.FlowId)

	expLength := 4 + 4 + 2 // flowid + tunnels number + h.cmd

	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.TunnelResponse.FlowId, "\nThe two FlowId should be the same.")
}

func TestDelTunnelsRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 5 
		},
		"DelTunnelsRequest": {
			"FlowId":1234,
			"Teids" : [
				1,2,3
			] 
		}
	}
	`

	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	expLength := cmdSize + 4 + 4 // flowId ,tunnels

	teids := len(msg.DelTunnelsRequest.Teids)
	var dummy *uint32
	teidl := unsafe.Sizeof(*dummy)

	l := int(teidl)*teids + int(cmdSize) + 8 // flowid, tieds
	msg.Header.Length = uint16(l)
	request, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("ZMQ Encode failed. Err:%v", err)
	}

	t.Logf("Delete Tunnels request:%s\n", hex.EncodeToString(request))

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

	if msg.Header.Command == zmqencdec.ZMQ_CMD_ERROR {
		t.Fatalf("Received error:%s", msg.ErrorResponse.Error)
	}
	t.Logf("Start response: cmd:%d len:%d", msg.Header.Command, msg.Header.Length)
	if msg.Header.Command == zmqencdec.ZMQ_CMD_MSG_ERROR {
		t.Fatalf("Received message error. flowid:%d error:%s", msg.MsgErrorResponse.FlowId, msg.MsgErrorResponse.Error)
	}

	t.Logf("Start response: cmd:%d len:%d flowid:%d ", msg.Header.Command, msg.Header.Length, msg.TunnelResponse.FlowId)

	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.TunnelResponse.FlowId, "\nThe two FlowId should be the same.")
}

func TestDelAllTunnelsRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 10,
		  "Command": 6 
		},
		"DelAllTunnelsRequest": {
			"FlowId":1234
		}
	}
	`

	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()
	expLength := cmdSize + 4 + 4 // flowId ,tunnels

	teids := len(msg.DelTunnelsRequest.Teids)
	var dummy *uint32
	teidl := unsafe.Sizeof(*dummy)

	l := int(teidl)*teids + int(cmdSize) + 8 // flowid, tieds
	msg.Header.Length = uint16(l)

	request, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("ZMQ Encode failed. Err:%v", err)
	}

	t.Logf("Delete Tunnels request:%s\n", hex.EncodeToString(request))
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

	if msg.Header.Command == zmqencdec.ZMQ_CMD_ERROR {
		t.Fatalf("Received error:%s", msg.ErrorResponse.Error)
	}
	t.Logf("Start response: cmd:%d len:%d", msg.Header.Command, msg.Header.Length)
	if msg.Header.Command == zmqencdec.ZMQ_CMD_MSG_ERROR {
		t.Fatalf("Received message error. flowid:%d error:%s", msg.MsgErrorResponse.FlowId, msg.MsgErrorResponse.Error)
	}

	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.TunnelResponse.FlowId, "\nThe two FlowId should be the same.")
}

func TestGetInfosRequest(t *testing.T) {

	jsonData := `
	{
		"Header":{
          "Length": 0,
		  "Command": 7 
		},
		"GetInfoRequest": {
			"FlowId":1234
		}
	}
	`
	expVersion := "1.0.0"
	expLength := len(expVersion) + 4 + 2 // flowId + h.cmd
	msg, err := jsonEncoder.Encode(jsonData)
	if err != nil {
		t.Fatalf("Json Encode failed. Err:%v", err)
	}

	cmdSize := reflect.TypeOf(msg.Header.Command).Size()

	l := int(cmdSize) + 4 // flowId
	msg.Header.Length = uint16(l)

	request, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("ZMQ Encode failed. Err:%v", err)
	}

	t.Logf("get info request:%s\n", hex.EncodeToString(request))
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
	t.Logf("Version:%s\n", msg.GetInfoResponse.Version)

	t.Logf("Start response: cmd:%d len:%d flowid:%d ", msg.Header.Command, msg.Header.Length, msg.GetInfoResponse.FlowId)
	assert.Equal(t, uint16(expLength), msg.Header.Length, "\nThe two h.length should be the same.")
	assert.Equal(t, uint32(1234), msg.GetInfoResponse.FlowId, "\nThe two FlowId should be the same.")
	assert.Equal(t, expVersion, msg.GetInfoResponse.Version, "\nThe two Vesrion should be the same.")
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

func calculateTunnelsLen(t *testing.T, tunnels []zmqencdec.JsonTunnel) int {
	var l int
	for _, tunnel := range tunnels {
		l += len(tunnel.UeIpV4) + len(tunnel.UpfIpV4) + 4 + 4
		t.Logf("l:%d", l)
	}
	return l
}
func encodeJsonAddTunnelsMessage(t *testing.T, msg *zmqencdec.Message) error {
	jsonTunnels := msg.AddJsonTunnelRequest.JsonTunnels
	var l int
	elems := len(jsonTunnels)
	msg.AddTunnelRequest.Tunnels = make([]zmqencdec.Tunnel,elems)
	msg.AddTunnelRequest.FlowId = msg.AddJsonTunnelRequest.FlowId

	for idx, jsontunnel := range jsonTunnels {
		ueIpV4str := jsontunnel.UeIpV4
		upfIpV4 := jsontunnel.UpfIpV4
		ipAddress := net.ParseIP(ueIpV4str)
		msg.AddTunnelRequest.Tunnels[idx].UeIpV4 = ip2int(ipAddress)
		ipAddress = net.ParseIP(upfIpV4)
		msg.AddTunnelRequest.Tunnels[idx].UpfIpV4 = ip2int(ipAddress)
		msg.AddTunnelRequest.Tunnels[idx].TeidIn = msg.AddJsonTunnelRequest.JsonTunnels[idx].TeidIn
		msg.AddTunnelRequest.Tunnels[idx].TeidOut = msg.AddJsonTunnelRequest.JsonTunnels[idx].TeidOut
		t.Logf("UeIpV4:%x UpfIpV4:%x", msg.AddTunnelRequest.Tunnels[idx].UeIpV4,msg.AddTunnelRequest.Tunnels[idx].UpfIpV4)
		l += len(jsontunnel.UeIpV4) + len(jsontunnel.UpfIpV4) + 4 + 4
		t.Logf("l:%d", l)
	}
	return nil
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}
