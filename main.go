package main

import (
	"log"
	"math"
	"time"
)

func main() {
	mqttConnector := InitMqtt()
	serial := InitSerial()

	topics := []string{
		"$eeprom/#", // Needs to be here BigClownTranslator in SerialReaderWriter
		"node/+/+/+/+/set",
		"node/+/+/+/+/get",
		"/nodes/+",
		"/pairing-mode/start",
		"/pairing-mode/stop",
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

	<-time.After(time.Duration(1000))
	mqttConnector.RequestAliases()
	<-time.After(time.Duration(math.MaxInt64))
}
