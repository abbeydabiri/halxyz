package main

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/HAL-xyz/ethrpc"
	"github.com/HAL-xyz/web3-multicall-go/multicall"
	"github.com/ethereum/go-ethereum/common"
)

var lastBlockNo = 0
var lastBlockHex = ""

const batchSize = 3
const goroutineSize = 3

const ethURL = "https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d"

// const ethURL = "https://mainnet.infura.io/v3/349d1750fa35425d9625f0fa1e03895e"

var wg sync.WaitGroup
var Batches chan multicall.ViewCalls

type CustomEthrpc struct {
	Eth *ethrpc.EthRPC
}

func (c CustomEthrpc) MakeEthRpcCall(cntAddress, data string, blockNumber int) (string, error) {
	params := ethrpc.T{
		To:   cntAddress,
		From: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		Data: data,
	}

	hexBlockNo := fmt.Sprintf("0x%x", blockNumber)
	return c.Eth.EthCall(params, hexBlockNo)
}

type Producer struct {
	MC multicall.Multicall
}

func (producer *Producer) HandleBatches() {
	wgLen := 0
	for batchedViewCalls := range Batches {
		wgLen++
		wg.Add(1)

		go func() {
			defer wg.Done()
			res, err := producer.MC.Call(batchedViewCalls, lastBlockHex)
			if err != nil {
				panic(err)
			}

			for key, call := range res.Calls {
				bigNumber := call.Decoded[0].(*big.Int)
				number := bigNumber.Int64()

				if number%2 == 0 {
					newTrigger := Triggers[key]
					newTrigger.Number = number
					ChanTriggers <- &newTrigger
				}
			}
			wgLen--
		}()
		if wgLen >= 3 {
			wg.Wait()
		}
	}
}

func (producer *Producer) BatchCalls() (err error) {

	Batches = make(chan multicall.ViewCalls, 3)
	var multicallList []multicall.ViewCall = nil
	ethRPC := CustomEthrpc{ethrpc.New(ethURL)}

	lastBlockNo, _ = ethRPC.Eth.EthBlockNumber()
	lastBlockHex = fmt.Sprintf("%x", lastBlockNo)

	producer.MC, err = multicall.New(ethRPC, multicall.ContractAddress(multicall.MainnetAddress))
	if err != nil {
		return
	}

	var sendBatch = func() {
		switch len(multicallList) {
		case 1:
			Batches <- multicall.ViewCalls{multicallList[0]}
		case 2:
			Batches <- multicall.ViewCalls{multicallList[0], multicallList[1]}
		case 3:
			Batches <- multicall.ViewCalls{multicallList[0], multicallList[1], multicallList[2]}
		case 4:
			Batches <- multicall.ViewCalls{multicallList[0], multicallList[1], multicallList[2], multicallList[3]}
		case 5:
			Batches <- multicall.ViewCalls{multicallList[0], multicallList[1], multicallList[2], multicallList[3], multicallList[4]}
		}
	}

	for key, trigger := range Triggers {

		if len(multicallList) > 0 && len(multicallList)%batchSize == 0 {
			multicallList = nil
		}
		newVC := multicall.NewViewCall(key, trigger.ContractAddress, trigger.Method, []interface{}{common.HexToAddress(trigger.UserAddress)})
		if errV := newVC.Validate(); errV != nil {
			println(errV.Error())
		}
		multicallList = append(multicallList, newVC)

	}
	sendBatch()
	multicallList = nil
	return
}
