package main

import (
	"gobot.io/x/gobot/platforms/firmata"
	"fmt"
	"errors"
)

// Talk to the Arduino-based Ground Control device via USB

// encapsulate a device connected to a /dev port
type Device struct {
	port    string
	adaptor *firmata.Adaptor
}

func NewDevice(port string) Device {
	firmataAdaptor := firmata.NewAdaptor(port)

	return Device{
		port:    port,
		adaptor: firmataAdaptor,
	}
}

func (d *Device) Toggle(t interface{}) error {
	switch v := t.(type) {
	case Led:
		//d.adaptor.
	default:
		return errors.New("Togglable device not found")
	}

	return nil
}


// read the current board status
type DeviceStatus struct {
	Buttons  []Button
	Switches []Switch
	Leds     []Led
	Rgbas    []Rgba
}

func (d *Device) Status() DeviceStatus {
	return DeviceStatus{}
}
