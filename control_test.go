package groundcontrol

import (
	"testing"
	"time"
	"log"
	"fmt"
)

func TestDevice(t *testing.T) {
	control := NewGroundControl("/dev/tty.usbmodem1411")

	err := control.Connect()
	if err {
		log.Fatal(err)
	}

	control.Display.Write("hello", "world")

	for _, led := range control.Leds {
		led.On()
		time.Sleep(1 * time.Second)
		led.Off()
	}

	control.Buzzer.Play(100, 100)

	for _, btn := range control.Switches {
		fmt.Println(btn.Active())
	}

	// press button 2 to stop the loop
	control.Buttons[1].Pushed(func() {
		control.Disconnect()
	})

	control.Switches[0].Pushed(func() {
		fmt.Println("Pushed the switch")
	})

	control.Switches[1].Released(func() {
		fmt.Println("Released the switch")
	})

	control.Loop(func() {
		level := control.Potentiometer.Level()
		switch {
		case level > 100:
			control.Leds[0].On()
		case level > 500:
			control.Leds[1].On()
		case level > 900:
			control.Leds[1].On()
		}

		states := "Btns: "
		for _, btn := range control.Buttons {
			var str string
			if btn.Active() {
				str = "1 "
			} else {
				str = "0 "
			}
			states += str
		}
		states += "\n"

		states += "Switches: "
		for _, btn := range control.Switches {
			var str string
			if btn.Active() {
				str = "1 "
			} else {
				str = "0 "
			}
			states += str
		}
		states += "\n"

		// TODO test buzzer
	})
}