package main

import "log"

func main() {
	mqttConnector := InitMqtt()
	serial := InitSerial(mqttConnector.CreateBigClownTranslator())

	topics := []string{
		"$eeprom/#",
		"node/+/+/+/+/set",
		"node/+/+/+/+/get",
		"/nodes/+",
		"/pairing-mode/#",
		"/info/get",
	}

	mqttConnector.ConsumeMessagesFromMqtt(topics, func(msg BcMessage) {
		log.Print("Sending msg to serial" + msg.String())
		serial.WriteSingleMessage(msg)
	})

	serial.ConsumeMessagesFromSerial(func(bcMsg BcMessage) {
		log.Print("Received msg from serial " + bcMsg.String())
		mqttConnector.Publish(bcMsg)
	})
}
