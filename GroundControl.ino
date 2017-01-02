//  _ ___ _______     ___ ___ ___  ___ _   _ ___ _____ ___ 
// / |_  )__ /   \   / __|_ _| _ \/ __| | | |_ _|_   _/ __| 
// | |/ / |_ \ |) | | (__ | ||   / (__| |_| || |  | | \__ \ 
// |_/___|___/___/   \___|___|_|_\\___|\___/|___| |_| |___/ 
// 
// Ground Control
// 
// Made by Herval Freire
// License: GPL 3.0
// Downloaded from: https://circuits.io/circuits/3611955-ground-control

#include <LiquidCrystal.h>


//---- LED outputs
int latchPin = 12; //Pin connected to ST_CP of 74HC595
int clockPin = 13; //Pin connected to SH_CP of 74HC595
int dataPin = 11; ////Pin connected to DS of 74HC595

void updateOutputs() {
    // count from 0 to 255 and display the number 
  // on the LEDs
  for (int numberToDisplay = 0; numberToDisplay < 256; numberToDisplay++) {
    // take the latchPin low so 
    // the LEDs don't change while you're sending in bits:
    digitalWrite(latchPin, LOW);
    // shift out the bits:
    shiftOut(dataPin, clockPin, MSBFIRST, numberToDisplay); 

    //take the latch pin high so the LEDs will light up:
    digitalWrite(latchPin, HIGH);
    // pause before next value:
    delay(1000);
  }
}

void setupOutputs() {
  //set pins to output so you can control the shift register
  pinMode(latchPin, OUTPUT);
  pinMode(clockPin, OUTPUT);
  pinMode(dataPin, OUTPUT);
}

//---- LED outputs



//---- LCD display
int lcdRs = 8;
int lcdEnable = 9;

int lcdD4 = 5;
int lcdD5 = 4;
int lcdD6 = 3;
int lcdD7 = 2;
//int lcdWiper = 3;
LiquidCrystal lcd(lcdRs, lcdEnable, lcdD4, lcdD5, lcdD6, lcdD7);

void setupLcd() {
  // set up the LCD's number of columns and rows:
  lcd.begin(16, 2);
  // Print a message to the LCD.
  lcd.print("hello, world!");
}

void updateDisplay() {
  // Turn off the display:
  lcd.noDisplay();
  delay(500);
  // Turn on the display:
  lcd.display();
  delay(500);
}
//---- LCD display



void setup() {
  setupOutputs();
  setupLcd();
}

void loop() {
  updateOutputs();
  updateDisplay();
}
