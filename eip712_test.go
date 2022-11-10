package eip712_go_example

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildDomainHash(t *testing.T) {
	expectedHash := "0xf2cee375fa42b42143804025fc449deafd50cc031ca257e0b194a650a912090f"

	hash, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash.String())
}

func TestBuildMailHash(t *testing.T) {
	expectedHash := "0x83fb919a6723739a9187fa6b145d321d7a747703fb65ff02e5adca18a3537c2a"

	hash, err := data.HashStruct(data.PrimaryType, data.Message)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash.String())
}

func TestBuildDataHash(t *testing.T) {
	expectedHash := "a4e05c27e7d8447be622e7d88a3097fbb2a31d0cff63064caca48ed57211bd49"

	hash, _, err := apitypes.TypedDataAndHash(data)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hex.EncodeToString(hash))
}

func TestVerifySignature(t *testing.T) {
	//[["Cow","0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"],["Bob","0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"],"Hello, Bob!",10000000000,1667659989]
	privateKey, err := crypto.HexToECDSA("your private key")
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	from := crypto.PubkeyToAddress(*publicKeyECDSA)

	data = TypedData{
		Types: Types{
			"EIP712Domain": EIP712DomainType,
			"Person":       PersonType,
			"Mail":         MailType,
		},
		PrimaryType: "Mail",
		Domain: TypedDataDomain{
			Name:              "Ether Mail",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(1),
			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		},

		Message: TypedDataMessage{
			"from":       map[string]interface{}{"name": "Cow", "wallet": from.String()},
			"to":         map[string]interface{}{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
			"contents":   "Hello, Bob!",
			"amount":     "10000000000",
			"expiration": "1667659989",
		},
	}

	digestHash, _, err := apitypes.TypedDataAndHash(data)
	sig, err := crypto.Sign(digestHash, privateKey)
	require.NoError(t, err)
	fmt.Printf("+%v\n", data)
	fmt.Printf("sig:0x%v", hex.EncodeToString(sig))

	verified := crypto.VerifySignature(publicKeyBytes, digestHash, sig[:64])
	require.True(t, verified)

}

func TestVerifySignature1(t *testing.T) {
	//[["Cow","0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"],["Bob","0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"],"Hello, Bob!",2000000000,1667659999]
	privateKey, err := crypto.HexToECDSA("your private key")
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	from := crypto.PubkeyToAddress(*publicKeyECDSA)

	data = TypedData{
		Types: Types{
			"EIP712Domain": EIP712DomainType,
			"Person":       PersonType,
			"Mail":         MailType,
		},
		PrimaryType: "Mail",
		Domain: TypedDataDomain{
			Name:              "Ether Mail",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(1),
			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		},

		Message: TypedDataMessage{
			"from":       map[string]interface{}{"name": "Cow", "wallet": from.String()},
			"to":         map[string]interface{}{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
			"contents":   "Hello, Bob!",
			"amount":     "2000000000",
			"expiration": "1667659999",
		},
	}

	digestHash, _, err := apitypes.TypedDataAndHash(data)
	sig, err := crypto.Sign(digestHash, privateKey)
	require.NoError(t, err)
	fmt.Printf("+%v\n", data)
	fmt.Printf("sig:%v", hex.EncodeToString(sig))

	verified := crypto.VerifySignature(publicKeyBytes, digestHash, sig[:64])
	require.True(t, verified)

}
