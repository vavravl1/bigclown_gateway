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

	result := MqttConnector{client}

	return result
}

func (connector *MqttConnector) ConsumeMessageFromMqtt(topics string, onMessage func(message BcMessage)) {
	topicsArray := []string{topics}
	connector.ConsumeMessagesFromMqtt(topicsArray, onMessage)
}

func (connector *MqttConnector) ConsumeMessagesFromMqtt(topics []string, onMessage func(message BcMessage)) {
	filters := make(map[string]byte)
	for _, pref := range topics {
		filters[pref] = 0
	}
	if token := connector.client.SubscribeMultiple(filters, connector.createMqttCallback(onMessage)); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (connector *MqttConnector) Publish(bcMsg BcMessage) {
	connector.client.Publish(
		bcMsg.topic,
		0,
		false,
		bcMsg.Bytes(),
	)
}

func (connector *MqttConnector) RequestAliases() {
	for i := 0; i < 4; i++ {
		connector.Publish(BcMessage{
			"$eeprom/alias/list",
			i,
		})
	}
}

func (connector *MqttConnector) createMqttCallback(callback func(message BcMessage)) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		var tmp interface{}
		if err := json.Unmarshal(msg.Payload(), &tmp); err != nil {
			log.Print("Unable to read message from mqtt: " + err.Error())
		}

		bcMsg := BcMessage{msg.Topic(), tmp}
		callback(bcMsg);
	}
}

func (connector *MqttConnector) CreateBigClownTranslator() BigClownTranslator {
	return BigClownTranslator{}
}
