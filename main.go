package main

import "log"

func main() {
	serial := InitSerial()
	mqttConnector := InitMqtt()

	mqttConnector.AddListener("/set", func(msg BcMessage) {
		serial.WriteSingleMessage(msg)
	})

	mqttConnector.AddListener("/nodes/get", func(msg BcMessage) {
		serial.WriteSingleMessage(msg)
	})

	mqttConnector.AddListener("/info/get", func(msg BcMessage) {
		serial.WriteSingleMessage(msg)
	})

	mqttConnector.AddListener("start", func(msg BcMessage) {
		serial.WriteSingleMessage(msg)
	})

	mqttConnector.AddListener("stop", func(msg BcMessage) {
		serial.WriteSingleMessage(msg)
	})

	serial.ConsumeMessagesFromSerial(func(bcMsg BcMessage) {
		log.Print("Received msg from serial " + bcMsg.String())
		mqttConnector.Publish(bcMsg)
	})
}
