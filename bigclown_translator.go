package main

import "strings"

type BigClownTranslator struct {
	nodeIdToName map[string]string
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


