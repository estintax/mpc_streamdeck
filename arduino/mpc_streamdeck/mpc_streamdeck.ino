#include <Key.h>
#include <Keypad.h>

byte rowKeys[1] = {2};
byte colKeys[4] = {4, 3, 6, 5};

char keys[1][4] = {
  {'1', '2', '3', '4'}
};

Keypad keypad = Keypad(makeKeymap(keys), rowKeys, colKeys, 1, 4);

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
}

void loop() {
  // put your main code here, to run repeatedly:
  char key = keypad.getKey();
  if(key) {
    Serial.println(key);
  }
}
