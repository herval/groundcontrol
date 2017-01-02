package main

func main() {
	driver := NewDevice("/dev/tty0")

	device := driver.Status()

	driver.Toggle(device.Leds[0])
}