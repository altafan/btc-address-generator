package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// Lambda type
type Lambda struct {
	XPub           *hdkeychain.ExtendedKey
	DerivationPath []uint32
	NetParams      *chaincfg.Params
	db             *dynamoDB
}

// NewLambda returns a new Lambda instance
func NewLambda(network, xpub, derviationPath string) (*Lambda, error) {
	masterPubKey, err := hdkeychain.NewKeyFromString(xpub)
	if err != nil {
		return nil, err
	}

	netParams, err := parseNetworkParams(network)
	if err != nil {
		return nil, err
	}

	path, err := parseDerivationPath(derviationPath)
	if err != nil {
		return nil, err
	}

	db := NewDynamoDB(xpub)
	if !db.ItemExists() {
		err := db.AddItem("0")
		if err != nil {
			return nil, err
		}
	}

	l := new(Lambda)
	l.XPub = masterPubKey
	l.DerivationPath = path
	l.NetParams = netParams
	l.db = db

	return l, nil
}

// HandleRequest implements lambda function
func (l *Lambda) HandleRequest(req events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, err error) {
	counter, err := l.db.GetCounter()
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()
		return
	}

	accountIndex := uint32(counter) + l.DerivationPath[len(l.DerivationPath)-1]
	if accountIndex >= hdkeychain.HardenedKeyStart {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = "Reached max number of derived receiving addresses"
		return
	}

	derivedKey := l.XPub
	for i, p := range l.DerivationPath {
		if i == len(l.DerivationPath)-1 {
			p = p + uint32(counter)
		}
		k, e := derivedKey.Child(p)
		if e != nil {
			resp.StatusCode = http.StatusInternalServerError
			resp.Body = e.Error()
			return
		}
		derivedKey = k
	}

	err = l.db.IncrementCounter()
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()
	}

	pubkey, _ := derivedKey.ECPubKey()
	pkhash := btcutil.Hash160(pubkey.SerializeCompressed())
	witAddr, _ := btcutil.NewAddressWitnessPubKeyHash(pkhash, l.NetParams)
	witnessProgram, _ := txscript.PayToAddrScript(witAddr)
	address, _ := btcutil.NewAddressScriptHash(witnessProgram, l.NetParams)

	resp.StatusCode = http.StatusOK
	resp.Body = address.String()

	return
}

// Start the lambda function
func (l *Lambda) Start() {
	lambda.Start(l.HandleRequest)
}

func parseNetworkParams(network string) (*chaincfg.Params, error) {
	if network == "mainnet" {
		return &chaincfg.MainNetParams, nil
	}

	if network == "testnet" {
		return &chaincfg.TestNet3Params, nil
	}

	if network == "regtest" {
		return &chaincfg.RegressionNetParams, nil
	}

	return nil, fmt.Errorf("invalid network: must be wither mainnet, testnet or regtest")
}

func parseDerivationPath(path string) ([]uint32, error) {
	splitpath := strings.Split(path, "/")

	if len(splitpath) < 2 {
		return nil, fmt.Errorf("derivation path must be in the form m/0 or m/0/0")
	}

	parsedPath := []uint32{}
	for _, p := range splitpath[1:] {
		parsedPath = append(parsedPath, stringPathToUint(p))
	}

	return parsedPath, nil
}

func isHardened(s string) bool {
	return strings.HasSuffix(s, "'")
}

func stringPathToUint(s string) uint32 {
	offset := uint32(0)
	stringPath := s
	if isHardened(s) {
		stringPath = strings.Replace(s, "'", "", 1)
		offset = hdkeychain.HardenedKeyStart
	}

	uintPath, _ := strconv.Atoi(stringPath)

	return offset + uint32(uintPath)
}
