package eip712_go_example

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// type aliases to avoid importing "core" everywhere
type TypedData = apitypes.TypedData
type TypedDataDomain = apitypes.TypedDataDomain
type Types = apitypes.Types
type Type = apitypes.Type
type TypedDataMessage = apitypes.TypedDataMessage

// EIP712DomainType is the type description for the EIP712 Domain
var EIP712DomainType = []Type{
	{Name: "name", Type: "string"},
	{Name: "version", Type: "string"},
	{Name: "chainId", Type: "uint256"},
	{Name: "verifyingContract", Type: "address"},
}

var PersonType = []Type{
	{Name: "name", Type: "string"},
	{Name: "wallet", Type: "address"},
}

var MailType = []Type{
	{Name: "from", Type: "Person"},
	{Name: "to", Type: "Person"},
	{Name: "contents", Type: "string"},
	{Name: "amount", Type: "uint256"},
	{Name: "expiration", Type: "uint256"},
}

type Person struct {
	Name   string `json:"name"`
	Wallet string `json:"wallet"`
}

type Mail struct {
	From       Person `json:"from"`
	To         Person `json:"to"`
	Contents   string `json:"contents"`
	Amount     string `json:"amount"`     //solidity's uint256 map to go's string in EIP712
	Expiration string `json:"expiration"` //solidity's uint256 map to go's string in EIP712
}

var data = TypedData{
	Types: Types{
		"EIP712Domain": EIP712DomainType,
		"Person":       PersonType,
		"Mail":         MailType,
	},
	PrimaryType: "Mail",
	Domain: TypedDataDomain{
		Name:              "Ether Mail",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(1), //solidity's uint256 map to go's math.NewHexOrDecimal256 in EIP712
		VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
	},

	Message: TypedDataMessage{
		"from":       map[string]interface{}{"name": "Cow", "wallet": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		"to":         map[string]interface{}{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
		"contents":   "Hello, Bob!",
		"amount":     "10000000000",
		"expiration": "1667659989",
	},
}
