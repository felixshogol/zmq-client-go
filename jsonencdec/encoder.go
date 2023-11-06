package jsonencdec

import (
	"encoding/json"
	"zmqclient/zmqencdec"

	"github.com/golang/glog"
)

type JsonEncoder struct {
}

// Encode - encode Messages
func (enc *JsonEncoder) Encode(jsonData string) (*zmqencdec.Message, error) {
	msg := &zmqencdec.Message{}

	glog.Infof("Encode Json message:%v", jsonData)

	//decoding the json data and storing in the message map
	err := json.Unmarshal([]byte(jsonData), msg)

	//Checks whether the error is nil or not
	if err != nil {
		//Prints the error if not nil
		glog.Errorf("Error while decoding the data", err.Error())
		return nil, err
	}
	return msg, nil
}

func (enc *JsonEncoder) Decode(msg *zmqencdec.Message) (string, error) {

	j, err := json.Marshal(msg)
	if err != nil {
		glog.Errorf("Error: %s", err)
		return "", err
	}

	return string(j), nil
}
