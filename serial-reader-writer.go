package main

import (
	"github.com/tarm/serial"
	"log"
	"encoding/json"
	"os"
)

type SerialReaderWriter struct {
	port *serial.Port
	bcTranslator BigClownTranslator
}

func InitSerial(bcTranslator BigClownTranslator) SerialReaderWriter {
	c := &serial.Config{Name: os.Getenv("BC_DEVICE"), Baud: BC_GATEWAY_DEVICE_BAUD_RATE}
	port, openSerialErr := serial.OpenPort(c)
	if openSerialErr != nil {
		log.Fatal(openSerialErr)
	}
	return SerialReaderWriter{port, bcTranslator}
}

func (readerWriter *SerialReaderWriter) WriteSingleMessage(message BcMessage) {
	toSerial := readerWriter.bcTranslator.FromMqttToSerial(message)
	_, serialWriteErr := readerWriter.port.Write([]byte(toSerial.toBigClownMessage()))
	if serialWriteErr != nil {
		log.Fatal(serialWriteErr)
	}
	log.Print("Message " + message.String() + " written to serial")
}

func (readerWriter *SerialReaderWriter) ConsumeMessagesFromSerial(callback func(message BcMessage)) {
	for ;; {
		line := readerWriter.readLine()
		var bcMsg BcMessage
		if parseBcMessageError := json.Unmarshal(line, &bcMsg); parseBcMessageError != nil {
			log.Print("Unable to parse message " + string(line) + " :" + parseBcMessageError.Error())
		} else {
			callback(readerWriter.bcTranslator.FromSerial(bcMsg))
		}
	}
}

func (readerWriter *SerialReaderWriter) readLine() []byte {
	singleCharBuffer := make([]byte, 1)
	result := make([]byte, 0)

	var bytesRead int
	var serialReadErr error

	for bytesRead, serialReadErr = readerWriter.port.Read(singleCharBuffer);
		serialReadErr == nil && bytesRead > 0 && singleCharBuffer[0] != '\n';
	bytesRead, serialReadErr = readerWriter.port.Read(singleCharBuffer) {
		result = append(result, singleCharBuffer[0])
	}

	if serialReadErr != nil {
		log.Fatal("Unable to read from serial: " + serialReadErr.Error())
	}

	return result
}
