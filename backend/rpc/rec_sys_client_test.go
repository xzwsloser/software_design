package rpc

import (
	"testing"

	"github.com/xzwsloser/software_design/backend/utils"
)

func TestGrpcClinet(t * testing.T) {
	utils.InitLogger()

	addr := "localhost"
	port := 7777

	NewGrpcClient(addr, port)

	resp, err := GetRecSysClient().GetRecResult(
		1,
		13,
		3,
		0,
		[]int{4,6,10},
		[]int{1,2,3},
		[]int{4,5,6},
		false,
		200,
	)

	if err != nil {
		t.Error(err.Error())
		return
	}

	for _, siteIdx := range resp {
		t.Logf("siteIdx = %d", siteIdx)
	}
}