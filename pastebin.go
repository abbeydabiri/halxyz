package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Pastebin struct {
	URL    string
	Client *http.Client
}

func (pastebin *Pastebin) Get() ([]byte, error) {

	resp, err := pastebin.Client.Get(pastebin.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (pastebin *Pastebin) Parse(body []byte) error {
	var triggers []Trigger
	Triggers = make(map[string]Trigger)

	if err := json.Unmarshal(body, &triggers); err != nil {
		return err
	}

	for _, trigger := range triggers {

		if trigger.TriggerName == "" || trigger.ContractAddress == "" || trigger.UserAddress == "" {
			continue
		}

		mapKey := fmt.Sprintf("%s-%s", trigger.TriggerName, trigger.UserAddress)
		Triggers[mapKey] = trigger
	}
	return nil
}
