package groundcontrol

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type Changeable interface {
	Changed() bool
	State() interface{}
}

type Rgba struct {
	address string
	color   string
	state   string
}

type Display struct {
	driver *i2c.GroveLcdDriver
	txt    string
}

func (d *Display) Init() {
	d.driver.Start()
	d.driver.Home()
	d.driver.Scroll(false)
	d.driver.Clear()
}

func (d *Display) Write(txt string) {
	if txt != d.txt {
		d.driver.Write(padRight(txt, " ", 32))
		d.txt = txt
	}
}

type Led struct {
	driver *gpio.LedDriver
	on     bool
}

func (l *Led) State() interface{} {
	return l.on
}

func (l *Led) On() {
	if !l.on {
		l.driver.On()
		l.on = true
	}

}

func (l *Led) Off() {
	if l.on {
		l.driver.Off()
		l.on = false
	}
}

func (l *Led) Changed() bool {
	return wasModified(l, func() {
		l.on = l.driver.State()
	})
}

type Potentiometer struct {
	pin     string
	level   int
	adaptor *firmata.Adaptor
}

func (p *Potentiometer) State() interface{} {
	return p.level
}

func (p *Potentiometer) Changed() bool {
	return wasModified(p, func() {
		level, err := p.adaptor.AnalogRead(p.pin)
		if err == nil {
			level = roundDown(level) // a bit less precision goes a long way
			p.level = level
		}
	})
}

func (p *Potentiometer) Level() int {
	return p.level
}

type Buzzer struct {
	driver *gpio.BuzzerDriver
}

func (b *Buzzer) Play(tone, duration float64) {
	b.driver.Tone(tone, duration)
}

type Button struct {
	driver *gpio.ButtonDriver
	port   string
	active bool
}

func (b *Button) State() interface{} {
	return b.active
}

func (b *Button) Changed() bool {
	return wasModified(b, func() {
		b.active = b.driver.Active
	})
}

func wasModified(b Changeable, modification func()) bool {
	prev := b.State()
	modification()
	return prev == b.State()
}

func (b *Button) Pushed(callback func()) {
	b.driver.On(gpio.ButtonPush, func(s interface{}) {
		b.active = true
		callback()
	})
}

func (b *Button) Released(callback func()) {
	b.driver.On(gpio.ButtonRelease, func(s interface{}) {
		b.active = false
		callback()
	})
}

func roundDown(n int) int {
	return n - n%10
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
