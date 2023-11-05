package encdec

type ZmqMessageType uint16

const (
	ZMQ_CMD_NONE ZmqMessageType = iota
	ZMQ_CMD_START
	ZMQ_CMD_STOP
	ZMQ_CMD_SHUTDOWN
	ZMQ_CMD_ADD_TUNNELS
	ZMQ_CMD_DEL_TUNNELS
	ZMQ_CMD_DEL_ALL_TUNNELS
	ZMQ_CMD_GET_INFO
	ZMQ_CMD_INVALID
)

type MsgHeader struct {
	Length  uint16
	Command ZmqMessageType
}

type MsgStartRequest struct {
	FlowId          uint32
	MetricsInterval uint32
}

type MsgStartResponse struct {
	FlowId    uint32
	Publicher string
}

type MsgResponse struct {
	FlowId uint32
}

type MsgTunnelResponse struct {
	FlowId  uint32
	Tunnels uint32
}

type MsgStopRequest struct {
	FlowId uint32
}

type Tunnel struct {
	TeidIn  uint32
	TeidOut uint32
	UeIpV4  uint32
	UpfIpV4 uint32
}

type MsgAddTunnelsRequest struct {
	FlowId  uint32
	Tunnels []Tunnel
}

type MsgDelTEIDsRequest struct {
	FlowId uint32
	Teids  []uint32
}

type MsgDelAllTEIDsRequest struct {
	FlowId uint32
}

type MsgGetInfoRequest struct {
	FlowId uint32
}

type Message struct {
	Header            MsgHeader
	StartRequest      MsgStartRequest
	StopRequest       MsgStopRequest
	StartResponse     MsgStartResponse
	Response          MsgResponse
	AddTunnelRequest  MsgAddTunnelsRequest
	DelTeidRequest    MsgDelTEIDsRequest
	DelAllTeidRequest MsgDelAllTEIDsRequest
	TunnelResponse    MsgTunnelResponse
	GetInfoRequest    MsgGetInfoRequest
}
