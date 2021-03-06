package contract

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase"
	"github.com/seeleteam/go-seele/common"
	"github.com/stretchr/testify/assert"
)

// const
const (
	CmdClient = "../../../bin/client"
	CmdLight  = "../../../bin/light"

	KeyFileShard11 = "../../../config/keyfile/shard1-0x4fb7c8b0287378f0cf8b5a9262bf3ef7e101f8d1"
	KeyFileShard12 = "../../../config/keyfile/shard1-0xec759db47a65f6537d630517f6cd3ca39c6f93d1"
)

// HandleTx handle tx and return the receipt
func HandleTx(t *testing.T, amount int, command, from, contract, payload string) (receipt *testcase.ReceiptInfo) {
	txHash, _, err1 := testcase.SendTx(t, command, amount, 0, 0, from, contract, payload, testcase.ServerAddr)
	if err1 != nil {
		t.Fatal(err1)
	}

	timeoutC := time.After(150 * time.Second)
	for {
		var err2 error
		// fmt.Println("txHash:", txHash)
		receipt, err2 = testcase.GetReceipt(t, command, txHash, testcase.ServerAddr)
		if err2 != nil && !strings.Contains(err2.Error(), "leveldb: not found") {
			t.Fatal(err2)
		}

		// fmt.Println("receipt:", receipt)
		if receipt != nil && receipt.Failed == false {
			break
		}

		select {
		case <-timeoutC:
			t.Fatalf("over time. err: %s", err2.Error())
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	return receipt
}

// GeneratePayload generate a payload
func GeneratePayload(t *testing.T, command, abi, method string, args ...string) (payload string) {
	cmd := exec.Command(command, "payload", "--abi", abi, "--method", method)
	for _, arg := range args {
		cmd.Args = append(cmd.Args, "--args", arg)
	}

	bytes, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	return string(bytes)
}

// ParseBinFile parse bin
func ParseBinFile(t *testing.T, filePath string) string {
	if !common.FileOrFolderExists(filePath) {
		t.Fatal("bin file not found")
	}
	bytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	return string(bytes)
}
