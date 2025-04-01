start-node:
	geth --http --http.api web3,eth,net,personal --http.addr 0.0.0.0 --http.port 8545

all: build run
