package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const (
	ethToUsdPriceFeedAddress = "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419"
	ethereumNodeURL          = "wss://mainnet.gateway.tenderly.co"
)

const priceFeedABI = `[{
	"inputs": [],
	"name": "latestRoundData",
	"outputs": [
		{"name": "roundId", "type": "uint80"},
		{"name": "answer", "type": "int256"},
		{"name": "startedAt", "type": "uint256"},
		{"name": "updatedAt", "type": "uint256"},
		{"name": "answeredInRound", "type": "uint80"}
	],
	"stateMutability": "view",
	"type": "function"
}]`

const (
	GREEN = "\033[32m"
	RESET = "\033[0m"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <Ethereum Address>")
	}
	ethAddress := os.Args[1]

	if !common.IsHexAddress(ethAddress) {
		log.Fatal("Invalid Ethereum address. Please try again.")
	}

	client, err := ethclient.Dial(ethereumNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	account := common.HexToAddress(ethAddress)
	weiBalance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal("Failed getting balance : ", err)
	}

	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(weiBalance), big.NewFloat(1e18))

	course, err := getEthCourse(client)
	if err != nil {
		log.Fatal("Failed getting course : ", err)
	}

	usdBalance := new(big.Float).Mul(ethBalance, course)

	fmt.Printf(GREEN+"Balance in Ethereum:%s  %v\n", RESET, ethBalance)
	fmt.Printf(GREEN+"Balance in USD:%s       %.2f\n", RESET, usdBalance)
}

type AbiData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func getEthCourse(client *ethclient.Client) (*big.Float, error) {
	contractAbi, err := abi.JSON(strings.NewReader(priceFeedABI))
	if err != nil {
		return nil, errors.Wrap(err, "abi.JSON")
	}

	contractAddress := common.HexToAddress(ethToUsdPriceFeedAddress)

	data, err := contractAbi.Pack("latestRoundData")
	if err != nil {
		return nil, errors.Wrap(err, "Pack")
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "CallContract")
	}

	var roundData AbiData

	err = contractAbi.UnpackIntoInterface(&roundData, "latestRoundData", result)
	if err != nil {
		return nil, errors.Wrap(err, "UnpackIntoInterface")
	}

	price := new(big.Float).Quo(
		new(big.Float).SetInt(roundData.Answer),
		new(big.Float).SetInt(big.NewInt(1e8)),
	)

	return price, nil
}
