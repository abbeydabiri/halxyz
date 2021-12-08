package main

import (
	"testing"
)

func TestBatchCalls(t *testing.T) {

	Triggers = make(map[string]Trigger)
	Triggers["BAT-0x0000000000000000000000000000000000000000"] = Trigger{
		TriggerName:     "BAT",
		UserAddress:     "0x0000000000000000000000000000000000000000",
		ContractAddress: "0x0d8775f648430679a709e98d2b0cb6250d2887ef",
		Method:          "balanceOf(address)(uint256)",
	}

	Triggers["DGD-0x0000000000000000000000000000000000000000"] = Trigger{
		TriggerName:     "DGD",
		UserAddress:     "0x0000000000000000000000000000000000000000",
		ContractAddress: "0xe0b7927c4af23765cb51314a0e0521a9645f0e2a",
		Method:          "balanceOf(address)(uint256)",
	}

	producer := Producer{}
	err := producer.BatchCalls()
	ok(t, err)
}
