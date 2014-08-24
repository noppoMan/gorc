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
}

func (proto *Protocol) Initialize(prop map[string]interface{}) *Protocol {
	proto.From = prop["From"].(string)
	proto.Body = prop["Body"].(string)
	proto.Date = time.Now().Format(time.ANSIC)
	return proto
}

func (proto *Protocol) Parse(stream string) *Protocol {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(stream), &dat); err != nil {
		panic(err)
	}

	proto.From = dat["From"].(string)
	proto.Date = dat["Date"].(string)
	proto.Body = dat["Body"].(string)
	return proto
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
