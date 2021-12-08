package main

import (
	"bytes"
	"io/ioutil"

	"net/http"

	"testing"
)

const baseUrl = "https://pastebin.com/raw/PjFkFsAk"
const sampleTriggers = `[
  {
    "TriggerName": "BAT",
    "UserAddress": "0x0000000000000000000000000000000000000000",
    "ContractAddress": "0x0d8775f648430679a709e98d2b0cb6250d2887ef",
    "Method": "balanceOf(address)(uint256)"
  },
  {
    "TriggerName": "DGD",
    "UserAddress": "0x0000000000000000000000000000000000000000",
    "ContractAddress": "0xe0b7927c4af23765cb51314a0e0521a9645f0e2a",
    "Method": "balanceOf(address)(uint256)"
  },
  {
    "TriggerName": "Bytom",
    "UserAddress": "0x0000000000000000000000000000000000000000",
    "ContractAddress": "0xcb97e65f07da24d46bcdd078ebebd7c6e6e3d750",
    "Method": "balanceOf(address)(uint256)"
  },
  {
    "TriggerName": "BAT",
    "UserAddress": "0x0000000000000000000000000000000000000000",
    "ContractAddress": "0x0d8775f648430679a709e98d2b0cb6250d2887ef",
    "Method": "balanceOf(address)(uint256)"
  }
]`

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGet(t *testing.T) {

	client := NewTestClient(func(req *http.Request) *http.Response {
		equals(t, req.URL.String(), baseUrl)
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(sampleTriggers)),
			Header:     make(http.Header),
		}
	})

	pb := Pastebin{baseUrl, client}
	body, err := pb.Get()
	ok(t, err)
	equals(t, []byte(sampleTriggers), body)
}

func TestParse(t *testing.T) {
	pb := Pastebin{}
	err := pb.Parse([]byte(sampleTriggers))
	ok(t, err)
	equals(t, 3, len(Triggers))
}
