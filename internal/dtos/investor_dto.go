package dtos

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"scaffold-api-server/internal/entities"
)

type InvestorDTO struct {
	entities.Investor
}

func (d *InvestorDTO) ToEntity() (interface{}, error) {
	return &d.Investor, nil
}

func (d *InvestorDTO) IsValid() bool {
	client, err := ethclient.Dial(viper.GetString("network.url"))
	if err != nil {
		logrus.Error(fmt.Sprintf("connect network failed, details: %v", err))
		return false
	}
	defer client.Close()
	if !common.IsHexAddress(d.Address) {
		return false
	}
	b, err := client.BalanceAt(context.Background(), common.HexToAddress(d.Address), nil)
	if err != nil {
		logrus.Error(fmt.Sprintf("get balance failed, details: %v", err))
		return false
	}
	logrus.Info("balance(wei): ", b.String())
	return true
}
