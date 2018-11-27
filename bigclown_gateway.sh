#!/bin/bash

export MQTT_BROKER_URL='tcp://localhost:1883'
export MQTT_BROKER_USERNAME='XXX'
export MQTT_BROKER_PASSWORD='XXX'
export BC_DEVICE='/dev/ttyACM0'

cd /home/vlada/services/bigclown_gateway
#until /home/vlada/services/bigclown_gateway/bigclown_gateway >> service.log 2>&1; do
until ./bigclown_gateway; do
    echo "Service 'bigclown_gateway' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done