package main

type Trigger struct {
	TriggerName, UserAddress,
	ContractAddress, Method,
	Status string
	Number int64
}

var Triggers map[string]Trigger
