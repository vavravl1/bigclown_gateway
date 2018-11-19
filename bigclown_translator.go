package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type BigClownTranslator struct {
	nodeIdToName map[string]string
}

func InitBigClownTranslator() BigClownTranslator {
	aliases := make(map[string]string)
	if data, err := ioutil.ReadFile("./stored-aliases.json"); err != nil {
		panic(err)
	} else {
		json.Unmarshal(data, &aliases)
	}
	return BigClownTranslator{aliases}
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
		t.storeAlias()
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

func (t *BigClownTranslator) storeAlias() {
	if lastValuesToStore, err := json.MarshalIndent(t.nodeIdToName, "", "    "); err != nil {
		panic(err)
	} else {
		if err := ioutil.WriteFile("./stored-aliases.json", lastValuesToStore, 0600); err != nil {
			panic(err)
		}
	}
}
