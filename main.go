package am2301

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio"
)

type Reading struct {
	Temperature      float64
	RelativeHumidity float64
}

// waitChange waited for the pin to change it's state from the given state, for a
// maximum time of timeout. If the state doesn't change, it returns an error. If
// the state changes withing the timeout, the amount of time it took for the
// change to occur is returned
func waitChange(pin rpio.Pin, mode rpio.State, timeout time.Duration) (time.Duration, error) {
	var voltage1, voltage2, voltage3 rpio.State
	start := time.Now()

	for time.Since(start) < timeout {
		/* Primitive low-pass filter */
		voltage1 = pin.Read()
		voltage2 = pin.Read()
		voltage3 = pin.Read()
		if voltage1 == voltage2 && voltage2 == voltage3 && voltage3 == mode {
			return time.Since(start), nil
		}
	}
	return 0, errors.Errorf("Timeout exceeded while waiting for change "+
		"on pin %d to mode %d", pin, mode)
}

func GetReading(pin rpio.Pin, mode rpio.Mode) (Reading, error) {
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
	if _, err := waitChange(pin, rpio.High, 100*time.Microsecond); err != nil {
		return reading, err
	}
	/* Wait for ACK */
	if _, err := waitChange(pin, rpio.Low, 100*time.Microsecond); err != nil {
		return reading, err
	}

	if _, err := waitChange(pin, rpio.High, 100*time.Microsecond); err != nil {
		return reading, err
	}

	/* When restarting, it looks like this lookfor start bit is not needed */
	if mode != 0 {
		/* Wait for the start bit */
		if _, err := waitChange(pin, rpio.Low, 200*time.Microsecond); err != nil {
			return reading, err
		}
		if _, err := waitChange(pin, rpio.High, 200*time.Microsecond); err != nil {
			return reading, err
		}
	}

	var reads [5]uint64
	for read_counter := 0; read_counter < 5; read_counter++ {
		for exponent, read := 7, 0; exponent >= 0; exponent-- {
			timeTilChange, err := waitChange(pin, rpio.Low, 500*time.Microsecond)
			if err != nil {
				return reading, err
			}

			readDigit := 0

			if timeTilChange >= 50*time.Microsecond {
				readDigit = 1
			}
			// read = read + (read_digit * 2^exponent)
			read = read | (read_digit << exponent)
			if _, err = waitChange(pin, rpio.High, 500*time.Microsecond); err != nil {
				return reading, err
			}
		}
		reads[read_counter] = uint64(read)
	}

	pin.Output()
	pin.High()

	/* Verify checksum */
	checksum := reads[0] + reads[1] + reads[2] + reads[3]
	if checksum != reads[4] {
		return reading, errors.New("Checksum check failed!")
	}

	reading.RelativeHumidity = float64((reads[0] << 8) | reads[1])
	reading.RelativeHumidity /= 10.0
	reading.Temperature = float64((reads[2] << 8) | reads[3])
	reading.Temperature /= 10.0

	if reading.RelativeHumidity > 100.0 || reading.RelativeHumidity < 0.0 {
		return reading, errors.New("Read relative humidity value non-sensical")
	}
	if reading.Temperature > 80.0 || reading.Temperature < -40.0 {
		return reading, errors.New("Read relative humidity value out of bounds of sensor")
	}
	return reading, nil
}

func GetTemperature(pin rpio.Pin) (float64, error) {
	return 0, nil
}

func GetRelativeHumidity(pin rpio.Pin) (float64, error) {
	return 0, nil
}

func main() {
	debug := false
	debugString := os.Getenv("DEBUG")
	if debugString != "" {
		debug = true
	}
	pinNumberString := os.Getenv("PIN_NUMBER")
	if pinNumberString == "" {
		log.Fatal("Please provide env var PIN_NUMBER")
	}
	pinNumber, err := strconv.Atoi(pinNumberString)
	if err != nil {
		log.Fatal(err)
	}
	rpio.Open()
	defer rpio.Close()
	pin := rpio.Pin(pinNumber)
	for trial_counter := 0; trial_counter < 10; trial_counter++ {
		reading, err := GetReading(pin, 1)
		if err != nil && debug {
			log.Println(err)
			time.Sleep(2 * time.Second)
		} else {
			log.Println(reading.Temperature)
			log.Println(reading.RelativeHumidity)
			break
		}
	}
}
