package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var pbUrl, dbPath string
	flag.StringVar(&pbUrl, "pastebin", "https://pastebin.com/raw/PjFkFsAk", "pastebin url json file containing triggers")
	flag.StringVar(&dbPath, "sqlite", "halxyz.sqlite3", "sqlite3 file path to store triggers")
	flag.Parse()

	//1.A - Pastebin Get
	pastebin := Pastebin{pbUrl, &http.Client{}}
	body, err := pastebin.Get()
	if err != nil {
		log.Panic(err)
	}

	//1.B - Pastebin Parse
	if err := pastebin.Parse(body); err != nil {
		log.Panic(err)
	}

	//2.A - Producer Handle Batches GoRoutine
	producer := Producer{}
	if err := producer.BatchCalls(); err != nil {
		log.Panic(err)
	}

	//2.B - Consumer Hnadle Batch Result GoRoutine
	consumer := Consumer{}
	go consumer.HandleBatchResult()

	// for trigger, key := range Triggers {
	// 	log.Printf("%v: %v", trigger, key)
	// }

	//2.C - Producer Handle MultiBatch Calls
	producer.HandleBatches()
}
