# GoLang package for working with [Mertech QR-Pay](https://mertech.ru/displej-qr-kodov-mertech-qr-pay-red/) via serial port

## Supported commands:

| methods          |
|------------------|
| Connect          |
| Disconnect       |
| ScreenClear      |
| CheckFirmware    |
| ShowPicOk        |
| ShowPicFalse     |
| EnableBluetooth  |
| DisableBluetooth |
| ScreenOn         |
| ScreenOff        |
| ShowQr           |

## usage
```
var err error

mertech := NewMertechQr(&serial.Config{
    Name:        "COM3",
    Baud:        SpeedBaud,         // 115200
    ReadTimeout: time.Second,
    Size:        DataBits,          // 8 byte
})

err = mertech.Connect()

bytes, err := mertech.ShowQr()

err = mertech.Disconnect()
```

## Example use
```
package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	qr "github.com/sergw3x/mertech_qr"
	"go.bug.st/serial"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var err error

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Port"})
	for i, port := range ports {
		table.Append([]string{strconv.Itoa(i + 1), port})
	}
	table.Render()

	portNum := 0

	fmt.Println("for exit enter 0")
	fmt.Print("Select serial port to connect: ")
	_, _ = fmt.Scan(&portNum)

	if portNum == 0 {
		fmt.Println("\rExit. Bye")
		os.Exit(0)
	}

	if portNum > len(ports) || portNum < 0 {
		fmt.Println("\rWrong port num. Bye")
		os.Exit(1)
	}

	serialName := ports[portNum-1]
	fmt.Printf("Selected: %s\n", serialName)

	m := qr.NewMertechQr(&qr.Config{
		Name:        ports[portNum-1],
		Baud:        qr.SpeedBaud, // 115200
		ReadTimeout: time.Second,
		Size:        qr.DataBits, // 8 byte
	})

	if err = m.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func(m *qr.MertechQr) { _ = m.Disconnect() }(m)

	bytes, err := m.ShowQr("hello world")
	fmt.Printf("%s writed: %d", serialName, bytes)
}
```