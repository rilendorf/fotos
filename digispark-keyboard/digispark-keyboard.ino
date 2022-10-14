#include "DigiKeyboard.h"

// Button PIN is p2 (its not used for anything else)
const int button = 2;

void setup() {
  // Set pin mode
  pinMode(button, INPUT_PULLUP); 
}

void loop() {
  // something something legacy crap
  DigiKeyboard.sendKeyStroke(0);

  // read state
  int s = digitalRead(button);

  // If is shorted to GND (btn pressed)
  if (s == LOW) {
	DigiKeyboard.println(""); // US Keystyle, idc cuz enter is same on most keyboards
	DigiKeyboard.delay(1000); // debounce
  }
}
