package main

import (
	"time"
	"log"
	"strings"
	"os"
	"github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
)

type MqttConnector struct {
	client   mqtt.Client
	topicPrefix string
	callback func(message BcMessage)
}

func InitMqtt(topics []string, topicPrefix string, onMessage func(message BcMessage)) MqttConnector {
	opts := mqtt.NewClientOptions().
		AddBroker(os.Getenv("MQTT_BROKER_URL")).
		SetUsername(os.Getenv("MQTT_BROKER_USERNAME")).
		SetPassword(os.Getenv("MQTT_BROKER_PASSWORD"))
	opts.SetClientID("go-bigclown-gateway" + string(time.Now().Unix()))

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	result := MqttConnector{client, topicPrefix, onMessage}
        result.addListener(topics);

	return result
}

func (connector *MqttConnector) Publish(bcMsg BcMessage) {
        var topic string
	if strings.HasPrefix(bcMsg.topic, "$") || strings.HasPrefix(bcMsg.topic, "/") {
	    topic = bcMsg.topic
	} else {
	    topic = connector.topicPrefix + bcMsg.topic
	}
	connector.client.Publish(
		topic,
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
	if msg, err := connector.createBcMessageWithRemovedNodePrefix(msg.Topic(), msg.Payload()); err == nil {
		connector.callback(msg);
	} else {
		log.Print("Unable to read message from mqtt: " + err.Error())
	}
}

func (connector *MqttConnector)createBcMessageWithRemovedNodePrefix(topic string, value []byte) (BcMessage, error) {
	var tmp interface{}
	if err := json.Unmarshal(value, &tmp); err != nil {
		return BcMessage{}, err
	}
	return BcMessage{
		strings.Replace(topic, connector.topicPrefix, "", 1),
		tmp,
	}, nil
}
