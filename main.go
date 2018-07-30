package am2301

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type Reading struct {
	Temperature      float64
	RelativeHumidity float64
}


func waitChange(pin rpio.Pin, mode rpio.Mode, timeout time.Duration)
{
	var voltage1, voltage2, voltage3 int
	start := time.Now()

	for time.Since(start) < timeout {
	/* Primitive low-pass filter */
	v1 = digitalRead(_pin_am2301);
	v2 = digitalRead(_pin_am2301);
	v3 = digitalRead(_pin_am2301);
	if (v1 == v2 && v2 == v3 && v3 == mode) {
	    return (micros() - now);
	}
    } while ((micros() - now) < tmo);
    return -1;
}


func GetReading(pin rpio.Pin) (Reading, error) {
	reading := Reading{}

	/* Leave it high for a while */
	pin.Output()
	pin.High()
	time.Sleep(100 * time.Microsecond)

	/* Set it low to give the start signal */
	pin.Low()
	time.Sleep(1000 * time.Microsecond)

    /* Now set the pin high to let the sensor start communicating */
	pin.High()
	pin.Input()




	return reading, nil
}

func GetTemperature(pin rpio.Pin) (float64, error) {

}

func GetRelativeHumidity(pin rpio.Pin) (float64, error) {

}

func main() {
	debug := false
	debugString := os.Getenv("DEBUG")
	if debugString != "" {
		debug = true
	}
	pinNumber := os.Getenv("PIN_NUMBER")
	if pinNumber == "" {
		log.Fatal("Please provide env var PIN_NUMBER")
	}
	rpio.Open()
	defer rpio.Close()
	pin := rpio.Pin(PIN_NUMBER)
	for trial_counter := 0; trial_counter < 10; trial_counter++ {
		reading, err := GetReading(pin)
		if err != nil && debug {
			log.Println(err)
		} else {
			log.Println(reading.Temperature)
			log.Println(reading.RelativeHumidity)
			break
		}
	}
}
