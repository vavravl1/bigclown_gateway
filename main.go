package main

import "log"

func main() {
	serial := InitSerial()

	topics := []string{
		"$eeprom/#",
		"node/+/+/+/+/set",
		"node/+/+/+/+/get",
		"/nodes/get",
		"/nodes/+",
		"/nodes/get",
		"/nodes/get",
		"/pairing-mode/#",
		"/info/get",
	}
	mqttConnector := InitMqtt(topics, "node/" , func(msg BcMessage) {
		log.Print("Sending msg to serial" + msg.String())
		serial.WriteSingleMessage(msg)
	})


	serial.ConsumeMessagesFromSerial(func(bcMsg BcMessage) {
		log.Print("Received msg from serial " + bcMsg.String())
		mqttConnector.Publish(bcMsg)
	})
}
