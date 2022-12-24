package numbertowa

import (
	"fmt"
	"testing"
)

var text string = "Ini nomer hpnya 082312276687 sama 888222"

func TestFindNumberInString(t *testing.T) {
	actual := []string{"082312276687", "888222"}

	result := FindNumberInString(text)
	if result[0] != actual[0] {
		t.Errorf("Error, actual first number is %s but you got %s", actual[0], result[0])
	}

	if result[1] != actual[1] {
		t.Errorf("Error, actual second number is %s but you got %s", actual[1], result[1])
	}
}

func TestNumberToPhone(t *testing.T) {
	numbers := FindNumberInString(text)

	result, err := NumberToPhone(numbers)
	if err != nil {
		t.Errorf("Error, %s", err.Error())
	}

	for r := range result {
		fmt.Println(result[r])
	}
}

func TestPhoneToUri(t *testing.T) {
	numbers := FindNumberInString(text)

	phone, err := NumberToPhone(numbers)
	if err != nil {
		t.Errorf("Error, %s", err.Error())
	}

	result, err := PhoneToUri(phone)
	if err != nil {
		t.Errorf("Error, %s", err.Error())
	}

	for r := range result {
		fmt.Println(result[r])
	}
}
