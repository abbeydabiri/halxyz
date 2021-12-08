package main

import "fmt"

var ChanTriggers = make(chan *Trigger)

type Consumer struct {
}

func (consumer *Consumer) HandleBatchResult() {

	db := Database{&SQLite{}}
	db.Connect("halxyz.sqlite")

	for trigger := range ChanTriggers {
		db.Save(trigger)
		fmt.Printf("Trigger: %s Saved to DB!! \n", trigger.TriggerName)
		fmt.Printf("%+v \n\n", trigger)
	}
}
