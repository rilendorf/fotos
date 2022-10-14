#include "DigiKeyboard.h"

const int button = 2;

void setup() {
  pinMode(button, INPUT_PULLUP); 
}

void loop() {
  // something something legacy crap
  DigiKeyboard.sendKeyStroke(0);

  int s = digitalRead(button);

  if (s == LOW) {
    DigiKeyboard.println(""); // US Keystyle, idc cuz enter is same on most keyboards
    DigiKeyboard.delay(1000);
  }
}
