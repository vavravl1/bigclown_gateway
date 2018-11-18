package main

import (
	"time"
	"log"
	"strings"
	"os"
	"github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
)

type mqttListener struct {
	topicMatcher string
	callback     func(message BcMessage)
}

type MqttConnector struct {
	client   mqtt.Client
	handlers []mqttListener
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
	return MqttConnector{client, make([]mqttListener, 0)}
}

func (connector *MqttConnector) AddListener(topicSuffix string, onMessage func(message BcMessage)) {
	connector.handlers = append(connector.handlers, mqttListener{topicSuffix, onMessage})
	topics := map[string]byte{
	    "#": 0,
	    "$eeprom/#": 0,
	}
	if token := connector.client.SubscribeMultiple(topics, connector.mqttCallback); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (connector *MqttConnector) Publish(bcMsg BcMessage) {
        var topic string
	if strings.HasPrefix(bcMsg.topic, "$") {
	    topic = bcMsg.topic
	} else {
	    topic = MQTT_TOPIC_PREFIX + bcMsg.topic
	}
	connector.client.Publish(
		topic,
		0,
		false,
		bcMsg.Bytes(),
	)
}

func (connector *MqttConnector) mqttCallback(client mqtt.Client, msg mqtt.Message) {
	for _, handler := range connector.handlers {
		if strings.HasSuffix(msg.Topic(), handler.topicMatcher) {
			log.Print("Client received message " + string(msg.Payload()) + " from " + msg.Topic())
			if msg, err := createBcMessageWithRemovedNodePrefix(msg.Topic(), msg.Payload()); err == nil {
				handler.callback(msg)
			} else {
				log.Print("Unable to read message from mqtt: " + err.Error())
			}
		}
	}
}

func createBcMessageWithRemovedNodePrefix(topic string, value []byte) (BcMessage, error) {
	var tmp interface{}
	if err := json.Unmarshal(value, &tmp); err != nil {
		return BcMessage{}, err
	}
	return BcMessage{
		strings.Replace(topic, MQTT_TOPIC_PREFIX, "", 1),
		tmp,
	}, nil
}
