package groundcontrol

import (
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot"
	"time"
	"gobot.io/x/gobot/drivers/i2c"
)

// Talk to the Arduino-based Ground Control device via USB
// encapsulate a device connected to a /dev port
type GroundControl struct {
	Display       *Display
	Leds          []*Led
	Potentiometer *Potentiometer
	Buttons       []*Button
	Switches      []*Button
	Buzzer        *Buzzer

	initCallback    func()
	workLoop        func()
	changedCallback func(interface{})
	adaptor         *firmata.Adaptor
	robot           *gobot.Robot
}

func NewGroundControl(port string) *GroundControl {
	control := GroundControl{}

	firmataAdaptor := firmata.NewAdaptor(port)

	control.Display = &Display{
		driver: i2c.NewGroveLcdDriver(firmataAdaptor),
	}

	// 3 leds
	control.Leds = setupLeds(firmataAdaptor)

	// a potentiometer
	control.Potentiometer = &Potentiometer{adaptor: firmataAdaptor, pin:"0"}

	// a buzzer
	control.Buzzer = &Buzzer{
		driver: gpio.NewBuzzerDriver(firmataAdaptor, "3"),
	}

	// 4 toggle switches
	control.Switches = setupSwitches(firmataAdaptor)

	// 2 buttons
	control.Buttons = setupButtons(firmataAdaptor)

	// wire everything w/ the "robot" so they can be stopped all together
	control.robot = gobot.NewRobot(
		"Ground Control",
		[]gobot.Connection{firmataAdaptor},
		allDevices(control),
		nil,
	)

	return &control
}

func setupButtons(firmataAdaptor *firmata.Adaptor) []*Button {
	btnPorts := []string{"1", "2"}
	buttons := make([]*Button, len(btnPorts))
	for i := range buttons {
		buttons[i] = &Button{
			driver: gpio.NewButtonDriver(firmataAdaptor, btnPorts[i]),
			port:   btnPorts[i],
		}
	}
	return buttons
}

func setupSwitches(firmataAdaptor *firmata.Adaptor) []*Button {
	switchPorts := []string{"4", "5", "6", "7"}
	switches := make([]*Button, len(switchPorts))
	for i := range switches {
		switches[i] = &Button{
			driver: gpio.NewButtonDriver(firmataAdaptor, switchPorts[i]),
			port:   switchPorts[i],
		}
	}

	return switches
}

func setupLeds(firmataAdaptor *firmata.Adaptor) []*Led {
	ledPorts := []string{"9", "10", "11" }
	leds := make([]*Led, len(ledPorts))
	for i := range leds {
		leds[i] = &Led{
			driver: gpio.NewLedDriver(firmataAdaptor, ledPorts[i]),
		}
	}

	return leds
}

func (control *GroundControl) Loop(callback func()) {
	control.workLoop = callback
}

func (control *GroundControl) Changed(callback func(interface{})) {
	control.changedCallback = callback
}

func (control *GroundControl) Init(callback func()) {
	control.initCallback = callback
}

func (g *GroundControl) Connect() error {
	g.robot.Work = func() {
		initialized := false

		// initialize everything
		g.Display.Init()

		// append callbacks, if any
		if g.initCallback != nil {
			g.initCallback()
		}
		initialized = true

		//gobot.Every(1*time.Second, func() {
		//	for _, led := range g.Leds {
		//		if led.blinking {
		//			led.Toggle()
		//		}
		//	}
		//})

		// poll for states of buttons and knobs and toggles
		gobot.Every(10*time.Millisecond, func() {
			if initialized {
				for _, btn := range g.Buttons {
					g.notifyChanged(btn)
				}
				for _, btn := range g.Switches {
					g.notifyChanged(btn)
				}
				for _, led := range g.Leds {
					g.notifyChanged(led)
				}
				g.notifyChanged(g.Potentiometer)

				if g.workLoop != nil {
					g.workLoop()
				}
			}
		})

	}
	return g.robot.Start()
}

func (control *GroundControl) notifyChanged(c Changeable) {
	if c.Changed() && control.changedCallback != nil {
		control.changedCallback(c)
	}
}

func (g *GroundControl) Disconnect() error {
	return g.robot.Stop()
}

func allDevices(control GroundControl) []gobot.Device {
	devices := make([]gobot.Device, 0)
	for i := range control.Switches {
		devices = append(devices, control.Switches[i].driver)
	}
	for i := range control.Buttons {
		devices = append(devices, control.Buttons[i].driver)
	}
	for i := range control.Leds {
		devices = append(devices, control.Leds[i].driver)
	}
	devices = append(devices, control.Display.driver)
	devices = append(devices, control.Buzzer.driver)

	return devices
}
