package main

import (
	"github.com/icrowley/fake"
	"github.com/joho/godotenv"
	"github.com/tarm/serial"
	"log"
	"testing"
	"time"
)

const (
	sleepTime = time.Second * 2
)

func connect() *MertechQr {
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serialName := envs["SERIAL"]

	conf := &serial.Config{
		Name:        serialName,
		Baud:        SpeedBaud,
		ReadTimeout: time.Second,
		Size:        DataBits,
		StopBits:    serial.Stop1,
		Parity:      serial.ParityNone,
	}

	mertech := NewMertechQr(conf)

	err = mertech.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return mertech
}

func disconnect(mertech *MertechQr) {
	if err := mertech.Disconnect(); err != nil {
		log.Fatal(err)
	}
}

func TestCheckFirmware(t *testing.T) {
	mertech := connect()
	defer disconnect(mertech)

	ver, err := mertech.CheckFirmware()
	if err != nil {
		t.Fatal(err)
	}

	if len(ver) == 0 {
		t.Fatal("error 'CheckFirmware' empty resp")
	}
}

func TestShowPicOk(t *testing.T) {
	mertech := connect()
	defer disconnect(mertech)

	ShowPicOkCnt, err := mertech.ShowPicOk()
	if err != nil {
		t.Fatal(err)
	}
	if ShowPicOkCnt == 0 {
		t.Fatal("error 'CheckFirmware' empty resp")
	}

	time.Sleep(sleepTime)

	screenClear, err := mertech.ScreenClear()
	if err != nil {
		t.Fatal(err)
	}

	if screenClear == 0 {
		t.Fatal("error 'screenClear'. wrote zero bytes")
	}
}

func TestShowPicFalse(t *testing.T) {
	mertech := connect()
	defer disconnect(mertech)

	ShowPicFalseCnt, err := mertech.ShowPicFalse()
	if err != nil {
		t.Fatal(err)
	}

	if ShowPicFalseCnt == 0 {
		t.Fatal("error 'CheckFirmware' empty resp")
	}

	time.Sleep(sleepTime)

	screenClear, err := mertech.ScreenClear()
	if err != nil {
		t.Fatal(err)
	}

	if screenClear == 0 {
		t.Fatal("error 'screenClear'. wrote zero bytes")
	}
}

func TestEnableBluetooth(t *testing.T) {
	mertech := connect()
	defer disconnect(mertech)

	EnableBluetoothCnt, err := mertech.EnableBluetooth()
	if err != nil {
		t.Fatal(err)
	}
	if EnableBluetoothCnt == 0 {
		t.Fatal("error 'CheckFirmware' empty resp")
	}

	time.Sleep(sleepTime)
}

func TestDisableBluetooth(t *testing.T) {
	mertech := connect()
	defer disconnect(mertech)

	DisableBluetoothCnt, err := mertech.DisableBluetooth()
	if err != nil {
		t.Fatal(err)
	}
	if DisableBluetoothCnt == 0 {
		t.Fatal("error 'CheckFirmware' empty resp")
	}

	time.Sleep(sleepTime)
}

func TestShowQr(t *testing.T) {
	mertech := connect()
	defer disconnect(mertech)

	ShowQr, err := mertech.ShowQr(fake.City())
	if err != nil {
		t.Fatal(err)
	}

	if ShowQr == 0 {
		t.Fatal("error 'ShowQr'. wrote zero bytes")
	}

	time.Sleep(sleepTime)

	screenClear, err := mertech.ScreenClear()
	if err != nil {
		t.Fatal(err)
	}

	if screenClear == 0 {
		t.Fatal("error 'screenClear'. wrote zero bytes")
	}
}
