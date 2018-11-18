package main

import (
	"encoding/json"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

type MqttConnector struct {
	client   mqtt.Client
	callback func(message BcMessage)
}

func InitMqtt() MqttConnector {
	opts := mqtt.NewClientOptions().
		AddBroker(os.Getenv("MQTT_BROKER_URL")).
		SetUsername(os.Getenv("MQTT_BROKER_USERNAME")).
		SetPassword(os.Getenv("MQTT_BROKER_PASSWORD"))
	opts.SetClientID("go-bigclown-gateway" + string(time.Now().Unix()))

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	result := MqttConnector{client, nil}

	return result
}

func (connector *MqttConnector) ConsumeMessagesFromMqtt(topics []string, onMessage func(message BcMessage)) {
	connector.callback = onMessage
	connector.addListener(topics)
}

func (connector *MqttConnector) Publish(bcMsg BcMessage) {
	connector.client.Publish(
		bcMsg.topic,
		0,
		false,
		bcMsg.Bytes(),
	)
}

func (connector *MqttConnector) addListener(_topics []string) {
	topics := make(map[string]byte)
	for _, pref := range _topics {
	    topics[pref] = 0
	}
	if token := connector.client.SubscribeMultiple(topics, connector.mqttCallback); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (connector *MqttConnector) mqttCallback(client mqtt.Client, msg mqtt.Message) {
	var tmp interface{}
	if err := json.Unmarshal(msg.Payload(), &tmp); err != nil {
		log.Print("Unable to read message from mqtt: " + err.Error())
	}

	bcMsg := BcMessage{msg.Topic(), tmp}
	connector.callback(bcMsg);
}

func (connector *MqttConnector) CreateBigClownTranslator() BigClownTranslator {
	return BigClownTranslator{}
}
