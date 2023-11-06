package zmqclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-zeromq/zmq4"
	zmq "github.com/go-zeromq/zmq4"
	"github.com/golang/glog"
)

type ZmqPacketHandler func(msg *zmq4.Msg)

type ZmqClientOptions struct {
	Host string
	Port int
}

type ZmqClient struct {
	connected    bool
	options      *ZmqClientOptions
	socket       zmq.Socket
	handler      ZmqPacketHandler
	listenerExit chan bool
}

func NewZmqClient(options *ZmqClientOptions) *ZmqClient {
	return &ZmqClient{
		options: options,
	}
}
func (c *ZmqClient) WithHandler(handler ZmqPacketHandler) error {
	c.handler = handler
	c.listenerExit = make(chan bool)

	if c.handler != nil {
		go func() {
			c.listen()
		}()
	}
	return nil
}

func (client *ZmqClient) Connect() error {
	ctx := context.Background()
	socket := zmq.NewReq(ctx, zmq.WithDialerRetry(time.Second))

	// todo after connect send info command to get dfxp version & information
	glog.Infof("connecting perf zmq %s:%d", client.options.Host, client.options.Port)
	if err := socket.Dial(fmt.Sprintf("tcp://%s:%d", client.options.Host, client.options.Port)); err != nil {
		return err
	}

	client.socket = socket
	client.connected = true
	return nil
}

func (client *ZmqClient) isConnected() bool {
	return client.connected
}

//Send - send packet to server
func (client *ZmqClient) Send(packet []byte) error {
	msg := zmq4.NewMsg(packet)
	if err := client.socket.Send(msg); err != nil {
		return fmt.Errorf("sending: %w", err)
	}
	return nil
}

func (client *ZmqClient) SendAndReceiveWithTimeout(packet []byte, to int) ([]byte, int, error) {
	if err := client.Send(packet); err != nil {
		return nil, 0, err
	}
	return client.ReceiveWithTimeout(to)
}


func (client *ZmqClient) ReceiveWithTimeout(to int) ([]byte, int, error) {
	select {
	case <-time.After(time.Duration(to) * time.Second):
		return nil, -1, nil
	default:
		// Wait for reply.
		r, err := client.socket.Recv()
		if err != nil {
			return nil,-1,fmt.Errorf("receiving: %w", err)
		}
		return r.Bytes(),0,nil
	}
}

func (client *ZmqClient) Close() error {
	if client.socket != nil {
		if err := client.socket.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ///////////////////////////////////////////////////////////
// Local API
// ///////////////////////////////////////////////////////////
// listen - listen zmq and return received a complete message
func (c *ZmqClient) listen() {
	for {
		select {
		case <-c.listenerExit:
			glog.Infoln("sctp client listen quit")
			return
		default:
			// Wait for message.
			if msg, err := c.socket.Recv(); err == nil {
				go func(msg *zmq4.Msg) {
					c.handler(msg)
				}(&msg)
			} else {
				glog.Errorf("zmq client recv error:%v", err)
			}
		}
	}
}

// //////////////////////////////////////////////////////////
func ZmqClientExe() error {
	ctx := context.Background()
	socket := zmq.NewReq(ctx, zmq.WithDialerRetry(time.Second))
	defer socket.Close()

	fmt.Printf("Connecting to hello world server...")
	if err := socket.Dial("tcp://127.0.0.1:5555"); err != nil {
		return fmt.Errorf("dialing: %w", err)
	}

	for i := 0; i < 10; i++ {
		// Send hello.
		m := zmq.NewMsgString("hello")
		fmt.Println("sending ", m)
		if err := socket.Send(m); err != nil {
			return fmt.Errorf("sending: %w", err)
		}

		// Wait for reply.
		r, err := socket.Recv()
		if err != nil {
			return fmt.Errorf("receiving: %w", err)
		}
		fmt.Println("received ", r.String())
	}
	return nil
}
