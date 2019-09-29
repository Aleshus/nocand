package network

import "C"
import (
	"flag"
	"github.com/Aleshus/nocand/clog"
	"github.com/Aleshus/nocand/models/can"
	"github.com/Aleshus/nocand/models/device"
	"golang.org/x/sys/unix"
	"net"
	"time"
)

var fd = 0
var CanTxChannel chan (can.Frame)
var CanRxChannel chan (can.Frame)
var DriverReady = false

func DriverUpdatePowerStatus() (*device.PowerStatus, error) {
	clog.Info("Update power status is not supported")
	status := &device.PowerStatus{}
	return status, nil
}

func DriverSetPower(powered bool) error {
	clog.Info("Set power is not supported")
	return nil
}

func DriverSetCurrentLimit(limit uint16) error {
	clog.Info("Set current limit is not supported")
	return nil
}

func DriverSendCanFrame(frame can.Frame) error {
	CanTxChannel <- frame
	return nil
}

func DriverInitialize(reset bool, speed uint) (*device.Info, error) {
	canIface := flag.String("canif", "can0", "The CAN interface")

	flag.Parse()
	iface, err := net.InterfaceByName(*canIface)
	if err != nil {
		clog.Fatal("%s", err)
	}

	fd, err = unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		clog.Fatal("%s", err)
	}

	addr := &unix.SockaddrCAN{Ifindex: iface.Index}

	err = unix.Bind(fd, addr)
	if err != nil {
		clog.Fatal("%s", err)
	}

	info := &device.Info{}

	return info, nil
}

func init() {
	CanTxChannel = make(chan (can.Frame), 32)
	CanRxChannel = make(chan (can.Frame), 32)

	go func() {
		for {
			frame := <-CanTxChannel
			start := time.Now()
			for C.digitalReadTx() == 0 {
				now := time.Now()
				for C.digitalReadTx() == 0 && time.Since(now).Seconds() < 3 {
				}
				if C.digitalReadTx() == 0 {
					clog.Warning("Microcontroller transmission has been blocking for more than %d seconds on frame %s.", int(time.Since(start).Seconds()), frame)
				}
			}

			buf := make([]byte, can.FRAME_LEN)

			if err := can.EncodeFrame(&frame, buf[:]); err != nil {
				clog.Error("Failed to encode CAN frame - %s", err)
			}

			_, err := unix.Write(fd, buf)
			if err != nil {
				clog.Error("Failed to send CAN frame - %s", err)
			}

			clog.DebugXX("SEND FRAME %s", frame)

		}
	}()

	go func() {

		for {
			buf := make([]byte, can.FRAME_LEN)
			n, err := unix.Read(fd, buf)
			if err != nil {
				clog.Error("Failed to read frame - %s", err)
			}
			if n != can.FRAME_LEN {
				clog.Error("Wrong frame length - %d", n)
			}

			frame, err := can.DecodeFrame(buf[:])
			clog.DebugXX("RECV FRAME %s", frame)
			CanRxChannel <- *frame
		}
	}()
}
