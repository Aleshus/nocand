package network

import "C"
import (
	"fmt"
	"github.com/Aleshus/nocand/models/can"
	"github.com/Aleshus/nocand/models/device"
	"time"
)

var CanTxChannel chan (can.Frame)
var CanRxChannel chan (can.Frame)
var DriverReady = false

func DriverUpdatePowerStatus() (*device.PowerStatus, error) {

}

func DriverSetPower(powered bool) error {

}

func DriverSetCurrentLimit(limit uint16) error {

}

func DriverSendCanFrame(frame can.Frame) error {

	return nil
}

func DriverInitialize(reset bool, speed uint) (*device.Info, error) {

	return info, nil
}
