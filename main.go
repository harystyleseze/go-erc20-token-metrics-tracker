package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

const (
	// The RPC URL of Ephemery testnet
	EPHEMERY_RPC_URL = "https://otter.bordel.wtf/erigon"
	// ERC-20 Token ABI
	ERC20_ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"}]`
)

var (
	client       *ethclient.Client
	contractAddr = common.HexToAddress("0x68E1Acf6b9f56267adDf65e1249B6aE321c0560E") // ERC20 token address
	contractABI  abi.ABI
)

func init() {
	// Connect to the Ethereum client
	var err error
	client, err = ethclient.Dial(EPHEMERY_RPC_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Parse ERC-20 Token ABI
	contractABI, err = abi.JSON(bytes.NewReader([]byte(ERC20_ABI)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}
}

// GetTokenDetails retrieves the basic details of the ERC-20 token
func GetTokenDetails(w http.ResponseWriter, r *http.Request) {
	// Get token name and symbol from the contract
	tokenName, err := getStringFromContract("name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching token name: %v", err), http.StatusInternalServerError)
		return
	}

	tokenSymbol, err := getStringFromContract("symbol")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching token symbol: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the response as JSON
	response := map[string]string{
		"name":   tokenName,
		"symbol": tokenSymbol,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTotalSupply retrieves the total supply of the ERC-20 token
func GetTotalSupply(w http.ResponseWriter, r *http.Request) {
	totalSupply, err := getUint256FromContract("totalSupply")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching total supply: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert total supply to a readable string
	totalSupplyStr := totalSupply.String()

	// Return the response as JSON
	response := map[string]string{
		"total_supply": totalSupplyStr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTokenBalance retrieves the balance of tokens for a specific address
func GetTokenBalance(w http.ResponseWriter, r *http.Request) {
	// Get address from query parameter
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	addressHex := common.HexToAddress(address)
	balance, err := getUint256FromContract("balanceOf", addressHex)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching token balance: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert balance to a readable string
	balanceStr := balance.String()

	// Return the response as JSON
	response := map[string]string{
		"balance": balanceStr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Utility function to fetch a string from the contract (e.g., token name, symbol)
func getStringFromContract(functionName string) (string, error) {
	data, err := callContract(functionName)
	if err != nil {
		return "", err
	}

	var result string
	err = contractABI.UnpackIntoInterface(&result, functionName, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Utility function to fetch a uint256 from the contract (e.g., totalSupply, balance)
func getUint256FromContract(functionName string, params ...interface{}) (*big.Int, error) {
	data, err := callContract(functionName, params...)
	if err != nil {
		return nil, err
	}

	var result *big.Int
	err = contractABI.UnpackIntoInterface(&result, functionName, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Generic function to call the contract and fetch data
func callContract(functionName string, params ...interface{}) ([]byte, error) {
	// Prepare the data to call the contract
	data, err := contractABI.Pack(functionName, params...)
	if err != nil {
		return nil, err
	}

	// Create the message for the transaction call
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	// Call the contract
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	// Set up the router
	r := mux.NewRouter()

	// Define the API routes
	r.HandleFunc("/api/token/details", GetTokenDetails).Methods("GET")
	r.HandleFunc("/api/token/totalSupply", GetTotalSupply).Methods("GET")
	r.HandleFunc("/api/token/balance", GetTokenBalance).Methods("GET")

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
