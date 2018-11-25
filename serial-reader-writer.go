package main

import (
	"encoding/json"
	"github.com/tarm/serial"
	"log"
	"os"
	"sync"
)

type SerialReaderWriter struct {
	port         *serial.Port
	bcTranslator BigClownTranslator
	callbacks    []func(message BcMessage)
	callbackMux *sync.Mutex
}

func InitSerial() *SerialReaderWriter {
	c := &serial.Config{Name: os.Getenv("BC_DEVICE"), Baud: BC_GATEWAY_DEVICE_BAUD_RATE}
	port, openSerialErr := serial.OpenPort(c)
	if openSerialErr != nil {
		log.Fatal(openSerialErr)
	}
	rslt := SerialReaderWriter{
		port,
		InitBigClownTranslator(),
		make([]func(message BcMessage), 0),
		&sync.Mutex{},
	}
	rslt.ConsumeMessagesFromSerial(rslt.bcTranslator.UpdateByMessage)
	go rslt.readingLoop()
	return &rslt
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
	readerWriter.callbackMux.Lock()
	defer readerWriter.callbackMux.Unlock()
	readerWriter.callbacks = append(readerWriter.callbacks, callback)
}

func (readerWriter *SerialReaderWriter) readingLoop() {
	for ; ; {
		line := readerWriter.readLine()
		var bcMsg BcMessage
		if parseBcMessageError := json.Unmarshal(line, &bcMsg); parseBcMessageError != nil {
			log.Panic("Unable to parse message " + string(line) + " :" + parseBcMessageError.Error())
		} else {
			readerWriter.interateOverCallbacks(readerWriter.bcTranslator.FromSerialToMqtt(bcMsg))
		}
	}
}

func (readerWriter *SerialReaderWriter) interateOverCallbacks(bcMsg BcMessage) {
	readerWriter.callbackMux.Lock()
	defer readerWriter.callbackMux.Unlock()
	for _, c := range readerWriter.callbacks {
		c(bcMsg)
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
