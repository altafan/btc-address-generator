package main

import (
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var tests = []struct {
	xpub              string
	network           string
	path              string
	expectedAddresses []string
}{
	{
		xpub:    "upub5FxWB33rLzN7YHBFqbu9tB9Xx26T2QZPHExHE9zqVYnoBQpryFqFMMZPMh7hezCMahwDfHcowFhvNMXab494VqgUiby5Z6xQ3c8ZpxgSbwt",
		network: "testnet",
		path:    "m/0",
		expectedAddresses: []string{
			"2MxzYwk78jNYPS6XGECqUSB82NvMhTMyzkA",
			"2NGSnhJWbzryQ8WoyJLQ4rU4ckoSTcVNxuT",
			"2N5Xpav1Qhwa2jy9cdZNuHYE65BspaACrhk",
			"2NGVSSaRFFEajBEQ1kDCFVZ16vqXeAeEDks",
			"2N8pxxS8uPjjTZi5BFT74de5pLqBicmWexP",
		},
	},
	{
		xpub:    "upub5F3UWX5R8VifmenDhnqHeXZjuyRtS5p6rs3tHNbUggRCTuCz6QSAanbhMsrvg7W4KH3WSGT1EeKPbjr7wiUPgcuq8Yspgp5EWWZbppZ6ncu",
		network: "testnet",
		path:    "m/10",
		expectedAddresses: []string{
			"2Mw5kmk5H7uoSfPxtqRhZLuLzT59Xme4KaX",
			"2N3gfpD9jc3cHefUZigeheY7Y8C3o7gXWPE",
			"2N5QoFdbv6jegFvSWy9TWAT2FAEq1GFQdbv",
			"2N4n9mpD5krQK8ragwjoZbKEu4bEMjGyk9H",
			"2NFBD1uk4fsTXiKYws3oWBHevPEgpqRRpMo",
		},
	},
	{
		xpub:    "ypub6YfHQjUErEQHjnc8S1CpEUVf62RiC87d4gVvDeHGf9Q4QEa8j68YAHQfs9q8YridTecRseVBQszioTDNDsTsRfzgDHUW7H5ZKgBXcsD8eeT",
		network: "mainnet",
		path:    "m/0",
		expectedAddresses: []string{
			"3QauCCcxUEat4gKtpXbQW67JevCm41wNVw",
			"349i6N4NytcdQ4gupxYFc17xyHPLcoqZgs",
			"399B3ZimhFcHkaNasaUHJWUrJUFvmszmE3",
			"32Rf8qXdewtR8AiQSSWmKZ1dYTDAtkwXZy",
			"3G4d5BbJFC7Hs7WaYccUmkkrbfS4EwkCPf",
		},
	},
	{
		xpub:    "upub5EaNW66c8VqYiZiVmQgdcSHLD59Z6msrsqDrBGsFE4rGJmfWyQUe8MppgF5iqCRYGNj6ARs6Cg7P8NAokSevJNRuJu17VBhgpiazMVDsbEM",
		network: "testnet",
		path:    "m/0/0",
		expectedAddresses: []string{
			"2MvjawTaVVc2o7nR9PRXyo8vF7GazVteRie",
			"2MtL8LnLDGmNMstiWhjYGkYLF53a1rkC73F",
			"2N3P7HSiH5WTgQ4jLgsRcBaVN5CRtm1ANur",
			"2N5A72K8etk8GHFErgVYX9ymGFdw5zUJYDg",
			"2My3nfFhKxYSZzVpeeYBh1Tz2wnodNhrMFk",
		},
	},
}

func TestHandler(t *testing.T) {
	for _, test := range tests {
		l, err := NewLambda(test.network, test.xpub, test.path)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(3 * time.Second)

		for _, expectedAddress := range test.expectedAddresses {
			resp, err := l.HandleRequest(events.APIGatewayProxyRequest{})
			if err != nil {
				t.Fatal(err)
			}
			address := resp.Body
			if address != expectedAddress {
				t.Fatalf("Got %s, expected: %s", address, expectedAddress)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
