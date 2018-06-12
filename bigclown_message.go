package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	switch v := message.value.(type) {
	case float64:
		return "[\"" + message.topic + "\"," + strconv.FormatFloat(v, 'f', 2, 64) + "]\n"
	case bool:
		return "[\"" + message.topic + "\"," + strconv.FormatBool(v)  + "]\n"
	case string:
		return "[\"" + message.topic + "\",\"" + v  + "\"]\n"
	default:
		//return "[\"" + message.topic + "\",\"unknown{" + reflect.TypeOf(message.value).String() + "}\"]\n"
		return "[\"" + message.topic + "\",\"unknown\"]\n"
	}
}

func (message BcMessage) Bytes() []byte {
	switch v := message.value.(type) {
	case float64:
		return []byte(strconv.FormatFloat(v, 'f', 2, 64))
	case bool:
		if v {
			return []byte("true")
		} else {
			return []byte("false")
		}
	case string:
		return []byte("\"" + v + "\"")
	default:
        return []byte("unknown")
	}
}
