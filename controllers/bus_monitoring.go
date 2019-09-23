package controllers

import (
	"github.com/Aleshus/nocand/models/device"
	"github.com/Aleshus/nocand/models/network"
	"github.com/Aleshus/nocand/models/rpi"

	"github.com/Aleshus/nocand/socket"
	"github.com/omzlo/clog"
	"time"
)

const (
	NO_BUS_RESET = false
	BUS_RESET    = true
)

var DeviceInfo *device.Info

func MilliAmpEstimation(c uint16) uint {
	var ma float64
	ma = 1000 * float64(c) / 4095 * 3.3 / 1120 * 2150
	return uint(ma)
}

func (nc *NocanNetworkController) RequestPowerStatusUpdate() {
	if network.DriverReady {
		ps, err := network.DriverUpdatePowerStatus()
		if err != nil {
			clog.Warning("Failed to read driver power status: %s", err)
		} else {
			clog.DebugX("Driver voltage=%.1f, current sense=%d (~ %d mA), reference voltage=%.2f, status(%x)=%s.", ps.Voltage, ps.CurrentSense, MilliAmpEstimation(ps.CurrentSense), ps.RefLevel, byte(ps.Status), ps.Status)
		}
		EventServer.Broadcast(socket.BusPowerStatusUpdateEvent, ps)
	}
}

func (nc *NocanNetworkController) RunPowerMonitor(interval time.Duration) {
	go func() {
		for {
			nc.RequestPowerStatusUpdate()
			time.Sleep(interval)
		}
	}()
}

func (nc *NocanNetworkController) Initialize(with_reset bool, spi_speed uint) error {
	di, err := network.DriverInitialize(with_reset, spi_speed)
	DeviceInfo = di
	return err
}

func (nc *NocanNetworkController) SetPower(power_on bool) {
	network.DriverSetPower(power_on)
	if power_on == false {
		Nodes.Clear()
	}
	EventServer.Broadcast(socket.BusPowerEvent, power_on)
}

func (nci *NocanNetworkController) SetCurrentLimit(limit uint16) {
	network.DriverSetCurrentLimit(limit)
	clog.DebugX("Driver current limit set to %d (~ %d mA)", limit, MilliAmpEstimation(limit))
}
