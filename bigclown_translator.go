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
	return BcMessage{
		strings.Replace(input.topic, "node/", "", 1),
		input.value,
	}
}

func (t *BigClownTranslator) FromSerial(input BcMessage) BcMessage {
	if strings.HasPrefix(input.topic, "$") || strings.HasPrefix(input.topic, "/") {
		return input
	} else {
		return BcMessage{"node/" + input.topic, input.value}
	}
}


