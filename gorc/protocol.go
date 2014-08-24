package gorc

import (
	"encoding/json"
	"time"
)

const (
	LE string = "\n"
)

type Protocol struct {
	Schema string
	Date   string
	From   string
	Body   string
	SesId  string
}

func (proto *Protocol) IsQuit() bool {
	return proto.Body == "quit"
}

func (proto *Protocol) JsonStringify() string {
	json, err := json.Marshal(proto)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (proto *Protocol) Dispfy() string {
	return "--------------" + LE +
		"SendFrom:" + proto.From + LE +
		"Data:" + proto.Date + LE +
		"Message:" + proto.Body + LE +
		"--------------" + LE
}

func NewProtocol(prop map[string]interface{}) *Protocol {
	return &Protocol{
		From:  prop["From"].(string),
		Body:  prop["Body"].(string),
		SesId: prop["SesId"].(string),
		Date:  time.Now().Format(time.ANSIC),
	}
}

func NewProtocolFromString(stream string) *Protocol {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(stream), &dat); err != nil {
		panic(err)
	}

	return &Protocol{
		From:  dat["From"].(string),
		Body:  dat["Body"].(string),
		Date:  dat["Date"].(string),
		SesId: dat["SesId"].(string),
	}
}
