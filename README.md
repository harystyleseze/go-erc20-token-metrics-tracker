# Go ERC-20 Token Metrics Tracker

This project provides a backend API to interact with an ERC-20 token on the Ethereum blockchain. It allows users to retrieve basic token details, total supply, and the balance of a specific address. The API connects to an Ethereum node, fetches data from a specified ERC-20 token contract, and returns it in a JSON format.

## Features

- **Get Token Details**: Fetches the name and symbol of the ERC-20 token.
- **Get Total Supply**: Retrieves the total supply of the token.
- **Get Token Balance**: Returns the balance of the specified address for the ERC-20 token.
- **Built with Go**: Leverages Go's `net/http` and the `go-ethereum` library for blockchain interaction.

## Prerequisites

Before running the application, ensure you have the following installed:

- **Go 1.18+**: This project uses Go for the backend.
- **Ethereum Node Access**: The project connects to an Ethereum node, which is available via an RPC URL. The testnet is used for interactions in this example.

## Setup Instructions

Follow these steps to set up the project on your local machine.

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/harystyleseze/go-erc20-token-metrics-tracker.git
cd go-erc20-token-metrics-tracker
```

### 2. Install Dependencies

Install the necessary Go dependencies using:

```bash
go mod tidy
```

This will download and install the required dependencies, including the `go-ethereum` library used to interact with the Ethereum blockchain.

### 3. Configure the Ethereum RPC URL

The project uses an Ethereum node (Ephemery testnet) as an RPC endpoint to interact with the blockchain. By default, it connects to the following URL:

```go
const EPHEMERY_RPC_URL = "https://otter.bordel.wtf/erigon"
```

If you want to use a different Ethereum node or testnet, replace the URL with your own node's RPC endpoint.

### 4. Update the Token Contract Address

The default contract address is set to an example ERC-20 token contract:

```go
var contractAddr = common.HexToAddress("0x68E1Acf6b9f56267adDf65e1249B6aE321c0560E")
```

If you want to interact with a different ERC-20 token, replace this with the desired token's contract address.

### 5. Run the Application

Start the Go backend API with the following command:

```bash
go run main.go
```

This will start an HTTP server on port `8080` that exposes the API endpoints.

## API Endpoints

The application exposes the following API routes:

### 1. **Get Token Details**

- **Endpoint**: `GET /api/token/details`
- **Description**: Fetches the name and symbol of the ERC-20 token.
- **Response**:

```json
{
  "name": "HarystylesMainToken",
  "symbol": "HMT"
}
```

### 2. **Get Total Supply**

- **Endpoint**: `GET /api/token/totalSupply`
- **Description**: Fetches the total supply of the ERC-20 token.
- **Response**:

```json
{
  "total_supply": "1000000000000000000000"
}
```

### 3. **Get Token Balance**

- **Endpoint**: `GET /api/token/balance?address=<address>`
- **Description**: Fetches the token balance of a specific Ethereum address.
- **Parameters**: The `address` query parameter must be a valid Ethereum address.
- **Example Request**:

```
GET /api/token/balance?address=0xYourEthereumAddress
```

- **Response**:

```json
{
  "balance": "900000000000000000000"
}
```

## Error Handling

The API returns HTTP error codes with a message in case of failure:

- **400 Bad Request**: Missing required parameters (e.g., missing address for balance query).
- **500 Internal Server Error**: If an error occurs while interacting with the Ethereum contract.

## Code Overview

### Main Components

- **`client`**: The `ethclient.Client` object used to interact with the Ethereum blockchain.
- **`contractAddr`**: The address of the ERC-20 token contract.
- **`contractABI`**: The ABI (Application Binary Interface) of the ERC-20 contract used to decode contract data.

### Functions

- **`init()`**: Initializes the Ethereum client and parses the ERC-20 contract ABI.
- **`GetTokenDetails()`**: Retrieves the token's name and symbol from the contract.
- **`GetTotalSupply()`**: Retrieves the total supply of the token.
- **`GetTokenBalance()`**: Retrieves the balance of a specified Ethereum address.
- **`getStringFromContract()`**: Helper function to fetch a string value from the contract (e.g., token name, symbol).
- **`getUint256FromContract()`**: Helper function to fetch a `uint256` value from the contract (e.g., total supply, balance).
- **`callContract()`**: Sends a call to the Ethereum contract and retrieves the response.

### Dependencies

- **`github.com/ethereum/go-ethereum`**: The official Go library for Ethereum.
- **`github.com/gorilla/mux`**: A lightweight router and URL matcher for building Go web applications.

## Testing the API

You can test the API using any HTTP client like:

- **Postman**: Send GET requests to the endpoints and view responses.
- **cURL**: Use the command line to send requests.

Example `cURL` request to get token details:

```bash
curl http://localhost:8080/api/token/details
```

### Automated Tests

Automated tests for the Go application can be added by writing test cases using Go's testing framework (`testing` package). You can create unit tests for the functions like `callContract()` and others that interact with the blockchain.

## License

This project is licensed under the MIT License.

## Contributing

Feel free to open issues or submit pull requests. Contributions are welcome to improve the functionality or fix bugs.
