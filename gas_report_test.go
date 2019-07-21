package main

import (
	"errors"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mariusgiger/iot-abe/cmd"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type GasTestSuite struct {
	suite.Suite
	accManager *acc.Manager
	wm         wallet.Manager
	rpc        rpc.Client
	log        *logrus.Logger
	gasStats   map[string]interface{}
}

var (
	owner              = common.HexToAddress("0x20683Db6E6d7ff53b62BCD6F723f74eC94dC410e")
	user               = common.HexToAddress("0x1e52b030261C4890A6aCe85Ed48CaE5f459525A0")
	device             = common.HexToAddress("0xE1097bAAA914277A8E2AefE464f8E29557e5f046")
	NumberOfAttributes = 6
	AttrLength         = 10
)

func (ts *GasTestSuite) TestDeploy() {
	ts.deploy()
}

func (ts *GasTestSuite) TestSetPolicy() {
	//arrange
	deployInfo := ts.deploy()
	contract := deployInfo.ContractAddress
	name := "Camera A"

	policy := "("
	for k := 1; k <= NumberOfAttributes; k++ {
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " and "
		}

		policy = policy + utils.RandString(AttrLength) + ")"
	}

	//act
	info, err := ts.accManager.SetDevicePolicy(name, device, contract, owner, policy)

	//assert
	ts.Require().NoError(err)
	ts.NotNil(info)
	ts.log.Info("created device policy:")
	utils.PrintFormatted(info)

	receipt, err := ts.WaitForConfirmation(info.Hash())
	ts.Require().NoError(err)

	ts.gasStats["policy"] = map[string]interface{}{
		"tx":        info.Hash().Hex(),
		"gasUsed":   receipt.GasUsed,
		"policy":    policy,
		"policyLen": len(policy),
	}
}

func (ts *GasTestSuite) TestRequestAccess() {
	ts.request()
}

func (ts *GasTestSuite) TestGrantAccess() {
	ts.grant()
}

func (ts *GasTestSuite) deploy() *acc.DeployInfo {
	//act
	info, err := ts.accManager.Deploy(owner)
	ts.Require().NoError(err)

	//assert
	ts.NotNil(info)
	ts.log.Info("created a new contract:")
	utils.PrintFormatted(info)

	receipt, err := ts.WaitForConfirmation(info.ContractTxHash)
	ts.Require().NoError(err)

	pubKeyBytes, err := hexutil.Decode(info.PublicKey)
	ts.Require().NoError(err)

	mkBytes, err := hexutil.Decode(info.MasterKey)
	ts.Require().NoError(err)

	ts.gasStats["deploy"] = map[string]interface{}{
		"tx":           info.ContractTxHash.Hex(),
		"gasUsed":      receipt.GasUsed,
		"pubKeyLen":    len(pubKeyBytes),
		"masterKeyLen": len(mkBytes),
	}

	return info
}

func (ts *GasTestSuite) request() (*acc.DeployInfo, common.Hash) {
	//arrange
	deployInfo := ts.deploy()
	contract := deployInfo.ContractAddress

	//act
	info, err := ts.accManager.RequestAccess(user, contract)

	//assert
	ts.Require().NoError(err)
	ts.NotNil(info)
	ts.log.Info("created access request:")
	utils.PrintFormatted(info)

	receipt, err := ts.WaitForConfirmation(info)
	ts.Require().NoError(err)

	ts.gasStats["request"] = map[string]interface{}{
		"tx":      info.Hex(),
		"gasUsed": receipt.GasUsed,
	}

	return deployInfo, info
}

func (ts *GasTestSuite) grant() (*acc.DeployInfo, *acc.GrantInfo) {
	//arrange
	deployInfo, requestTx := ts.request()
	var attrs []string
	for i := 0; i < NumberOfAttributes; i++ {
		attrs = append(attrs, utils.RandString(AttrLength))
	}

	tx, _, err := ts.rpc.TransactionByHash(requestTx)
	ts.Require().NoError(err)
	subjectPubKey, err := utils.RecoverPubKey(tx, big.NewInt(3)) //TODO move to constant
	ts.Require().NoError(err)

	//act
	info, err := ts.accManager.GrantAccess(deployInfo.ContractAddress, owner, user, attrs, deployInfo.PublicKey, deployInfo.MasterKey, subjectPubKey)

	//assert
	ts.Require().NoError(err)
	ts.NotNil(info)
	ts.log.Info("granted access:")
	utils.PrintFormatted(info)

	receipt, err := ts.WaitForConfirmation(info.TxHash)
	ts.Require().NoError(err)
	ts.gasStats["grant"] = map[string]interface{}{
		"tx":                 info.TxHash.Hex(),
		"gasUsed":            receipt.GasUsed,
		"keyLength":          len(info.SecretKey),
		"encryptedKeyLength": len(info.EncryptedSecretKey),
		"attributes":         len(attrs),
		"attrLength":         AttrLength,
	}

	return deployInfo, info
}

func (ts *GasTestSuite) SetupSuite() {
	log := logrus.New()
	cfg := &core.Config{}

	configFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(configFile), cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	cfg.DataDir = cmd.NormalizePath(cfg.DataDir)
	cfg.EthKeystoreDir = cmd.NormalizePath(cfg.EthKeystoreDir)

	wm := wallet.NewManager(log, cfg.EthKeystoreDir)
	wm.UsePW("1234")
	client, err := rpc.NewRPCClient(log, cfg)
	if err != nil {
		ts.Failf("could not create rpc client, %v", err.Error())
	}
	accManager, err := acc.NewManager(log, wm, client)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ts.accManager = accManager
	ts.wm = wm
	ts.rpc = client
	ts.log = log
	ts.gasStats = make(map[string]interface{})
}

func (ts *GasTestSuite) TearDownSuite() {
	ts.log.Info("gas stats:")
	ts.log.Info(utils.PrintFormatted(ts.gasStats))
}

var ConfTarget = big.NewInt(2)

func (ts *GasTestSuite) WaitForConfirmation(txHash common.Hash) (*types.Receipt, error) {
	var txReceipt *types.Receipt
	check := func() error {
		receipt, err := ts.rpc.TransactionReceipt(txHash)
		if err != nil {
			ts.log.Infof("tx not mined")
			return err
		}

		latestBlock, err := ts.rpc.BlockByNumber(nil)
		if err != nil {
			return err
		}

		if diff := big.NewInt(0).Sub(latestBlock.Number(), receipt.BlockNumber); ConfTarget.Cmp(diff) > 0 {
			ts.log.Infof("conf target not reached: current %v, wanted: %v", diff.Int64(), ConfTarget.Int64())
			return errors.New("conf target not reached")
		}

		txReceipt = receipt
		return nil // ok
	}

	err := backoff.Retry(check, backoff.NewConstantBackOff(time.Second*2))
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}

func TestGasTestSuite(t *testing.T) {
	suite.Run(t, &GasTestSuite{})
}
