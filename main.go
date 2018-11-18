package main

import "log"

func main() {
	serial := InitSerial()
	mqttConnector := InitMqtt()

	suffixes := []string{
		"/set",
		"/nodes/get",
		"/info/get",
		"start",
		"stop",
		"$eeprom/alias/add",
		"$eeprom/alias/remove",
		"$eeprom/alias/list",
	}

	for _,s := range suffixes {
		mqttConnector.AddListener(s, func(msg BcMessage) {
		        log.Print("Sending msg to serial" + msg.String())
			serial.WriteSingleMessage(msg)
		})

	}

	serial.ConsumeMessagesFromSerial(func(bcMsg BcMessage) {
		log.Print("Received msg from serial " + bcMsg.String())
		mqttConnector.Publish(bcMsg)
	})
}
