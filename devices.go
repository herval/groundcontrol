package groundcontrol

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/gpio"
	"fmt"
	"gobot.io/x/gobot/platforms/firmata"
)

type Rgba struct {
	address string
	color   string
	state   string
}

type Display struct {
	driver *i2c.GroveLcdDriver
	line1  string
	line2  string
}

func (d Display) Init() {
	d.driver.Start()
	d.driver.Home()
	d.driver.Scroll(false)
	d.driver.Clear()
}

func (d Display) Write(line1 string, line2 string) {
	d.driver.Write(
		fmt.Sprintf("%s\n%s", line1, line2),
	)
	d.line1 = line1
	d.line2 = line2
}

type Led struct {
	driver  *gpio.LedDriver
}

func (l Led) On() {
	l.driver.On()
}

func (l Led) Off() {
	l.driver.Off()
}

type Potentiometer struct {
	pin     string
	adaptor *firmata.Adaptor
}

func (p Potentiometer) Level() int {
	level, _ := p.adaptor.AnalogRead(p.pin)
	return roundDown(level) // a bit less precision goes a long way
}

type Buzzer struct {
	driver *gpio.BuzzerDriver
}

func (b Buzzer) Play(tone, duration float64) {
	b.driver.Tone(tone, duration)
}

type Button struct {
	driver *gpio.ButtonDriver
	port   string
}

func (b Button) Active() bool {
	return b.driver.Active
}

func (b Button) Pushed(callback func()) {
	b.driver.On(gpio.ButtonPush, func(s interface{}) {
		callback()
	})
}

func (b Button) Released(callback func()) {
	b.driver.On(gpio.ButtonRelease, func(s interface{}) {
		callback()
	})
}

func roundDown(n int) int {
	return n - n % 10
}
