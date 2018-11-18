package main

import (
	"encoding/json"
	"fmt"
)

type BcMessage struct {
	// Always without the MQTT_TOPIC_PREFIX (="node/") prefix
	topic  string
	value interface{}
}

func (message *BcMessage) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&message.topic, &message.value}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in IntBcMessage: %d != %d", g, e)
	}
	return nil
}

func (message BcMessage) String() string {
	return message.toBigClownMessage()
}

func (message BcMessage) toBigClownMessage() string {
	jsn,_ := json.Marshal(message.value)
        return "[\"" + message.topic + "\"," + string(jsn) + "]\n"
}

func (message BcMessage) Bytes() []byte {
	jsn,_ := json.Marshal(message.value)
	return jsn
}

