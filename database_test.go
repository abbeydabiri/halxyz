package main

import (
	"testing"
)

func TestOpen(t *testing.T) {
	database := Database{&SQLiteMock{}}
	err := database.Connect("hal.xyz")
	ok(t, err)
}

func TestFind(t *testing.T) {
	database := Database{&SQLiteMock{}}
	trigger, err := database.Find("BAT", "0x0000000000000000000000000000000000000000")
	ok(t, err)
	equals(t, "BAT", trigger.TriggerName)
}

func TestSave(t *testing.T) {
	database := Database{&SQLiteMock{}}
	trigger, err := database.Find("BAT", "0x0000000000000000000000000000000000000000")
	ok(t, err)
	equals(t, "BAT", trigger.TriggerName)
}
