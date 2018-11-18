package main

import (
	"fmt"
	"reflect"
	"strings"
)

type BigClownTranslator struct {
	nodeIdToName map[string]string
}

func InitBigClownTranslator() BigClownTranslator {
	return BigClownTranslator{make(map[string]string)}
}

func (t *BigClownTranslator) UpdateByMessage(bcMsg BcMessage) {
	if strings.HasPrefix(bcMsg.topic, "$eeprom/alias/list/") {
		switch v := bcMsg.value.(type) {
		case map[string]interface{}:
			for k, v2 := range v {
				switch v3 := v2.(type) {
				case string:
					t.nodeIdToName[k] = v3
				}
			}
		default:
			fmt.Println("Type:", reflect.TypeOf(v))
		}
		fmt.Println("BigClownTranslator updated ", t.nodeIdToName)
	}
}

func (t *BigClownTranslator) FromMqttToSerial(input BcMessage) BcMessage {
	topicForSerial := strings.Replace(input.topic, "node/", "", 1)
	for k,v := range t.nodeIdToName {
		topicForSerial = strings.Replace(topicForSerial, v, k, 1)
	}
	return BcMessage{
		topicForSerial,
		input.value,
	}
}

func (t *BigClownTranslator) FromSerialToMqtt(input BcMessage) BcMessage {
	topicToMqtt := input.topic
	for k,v := range t.nodeIdToName {
		topicToMqtt = strings.Replace(topicToMqtt, k, v, 1)
	}

	if strings.HasPrefix(input.topic, "$") || strings.HasPrefix(input.topic, "/") {
		return BcMessage{topicToMqtt, input.value}
	} else {
		return BcMessage{"node/" + topicToMqtt, input.value}
	}
}


