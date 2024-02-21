[![An old rock in the desert](https://mertech.ru/image/cache/catalog/goods/qr-code-display/1951/1951-1-700x700.jpg)](https://mertech.ru/image/cache/catalog/goods/qr-code-display/1951/1951-1-700x700.jpg)

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
    StopBits:    serial.Stop1,      // 1
    Parity:      serial.ParityNone, // 0
})

err = mertech.Connect()

bytes, err := mertech.ShowQr()

err = mertech.Disconnect()
```
