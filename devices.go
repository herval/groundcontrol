package main


type Rgba struct {
	address string
	color   string
	state   string
}

type Led struct {
	address string
	state   string
}

type Lcd struct {
	address string
	state   string
	text    string
}

type Switch struct {
	address string
	state   string
}

type Button struct {
	address string
	state   string
}