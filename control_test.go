package groundcontrol

import (
	"testing"
	"log"
	"fmt"
	"time"
)

func TestDevice(t *testing.T) {
	control := NewGroundControl("/dev/tty.usbmodem1411")

	control.Init(func() {
		control.Display.Write("hello\nworld")

		for _, led := range control.Leds {
			led.On()
			time.Sleep(1 * time.Second)
			led.Off()
		}

		for i := 100; i <= 5000; i++ {
			fmt.Println("Buzzing...")
			control.Buzzer.Play(float64(i), 2000)
			time.Sleep(2 * time.Second)
			i += 100
		}

		fmt.Println("Polling buttons...")
		for _, btn := range control.Switches {
			fmt.Println(btn.State())
		}
	})

	control.Switches[0].Pushed(func() {
		fmt.Println("Pushed the switch")
	})

	control.Switches[1].Released(func() {
		fmt.Println("Released the switch")
	})

	// press button 2 to stop the loop
	control.Buttons[1].Pushed(func() {
		control.Disconnect()
	})

	control.Changed(func(device interface{}) {
		fmt.Println(fmt.Sprintf("Changed: %+v", device))
	})

	control.Loop(func() {
		for _, led := range control.Leds {
			led.Off()
		}

		level := control.Potentiometer.Level()
		switch {
		case level > 900:
			control.Leds[0].On()
			control.Leds[1].On()
			control.Leds[2].On()
		case level > 600:
			control.Leds[0].On()
			control.Leds[1].On()
		case level > 300:
			control.Leds[0].On()
		}

		states := fmt.Sprintf("P %d ", level)

		states += "B "
		for _, btn := range control.Buttons {
			var str string
			if btn.active {
				str = "1 "
			} else {
				str = "0 "
			}
			states += str
		}
		states += "\n"

		states += "S "
		for _, btn := range control.Switches {
			var str string
			if btn.active {
				str = "1 "
			} else {
				str = "0 "
			}
			states += str
		}
		control.Display.Write(states)
	})

	err := control.Connect()
	if err != nil {
		fmt.Println("Make sure the board is connected before running this!")
		log.Fatal(err)
	}
}
