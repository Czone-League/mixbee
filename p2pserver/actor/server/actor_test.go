

package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/p2pserver"
	"github.com/mixbee/mixbee/p2pserver/common"
)

func TestP2PActorServer(t *testing.T) {
	log.Init(log.Stdout)
	fmt.Println("Start test the p2pserver by actor...")

	acct := account.NewAccount("SHA256withECDSA")
	p2p, err := p2pserver.NewServer(acct)
	if err != nil {
		t.Fatalf("TestP2PActorServer: p2pserver NewServer error %s", err)
	}

	p2pActor := NewP2PActor(p2p)
	p2pPID, err := p2pActor.Start()
	if err != nil {
		t.Fatalf("p2pActor init error %s", err)
	}

	//test server api

	//false: disable sync,running without ledger
	future := p2pPID.RequestFuture(&StartServerReq{StartSync: false}, common.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		t.Fatalf("TestP2PActorServer: p2p start error %s", err)
	}

	future = p2pPID.RequestFuture(&GetConnectionCntReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetConnectionCntReq error %s", err)
	}
	_, ok := result.(*GetConnectionCntRsp)
	if !ok {
		t.Error("GetConnectionCntRsp error")
	}

	future = p2pPID.RequestFuture(&GetNeighborAddrsReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetNeighborAddrsReq error %s", err)
	}
	_, ok = result.(*GetNeighborAddrsRsp)
	if !ok {
		t.Error("GetNeighborAddrsRsp error")
	}

	future = p2pPID.RequestFuture(&GetConnectionStateReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetConnectionStateReq error %s", err)
	}
	_, ok = result.(*GetConnectionStateRsp)
	if !ok {
		t.Error("GetConnectionStateRsp error")
	}

	future = p2pPID.RequestFuture(&GetTimeReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetTimeReq error %s", err)
	}
	_, ok = result.(*GetTimeRsp)
	if !ok {
		t.Error("GetTimeRsp error")
	}

	future = p2pPID.RequestFuture(&GetPortReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetPortReq error %s", err)
	}
	_, ok = result.(*GetPortRsp)
	if !ok {
		t.Error("GetPortRsp error")
	}

	future = p2pPID.RequestFuture(&GetIdReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetIdReq error %s", err)
	}
	_, ok = result.(*GetIdRsp)
	if !ok {
		t.Error("GetIdRsp error")
	}

	future = p2pPID.RequestFuture(&GetRelayStateReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetRelayStateReq error %s", err)
	}
	_, ok = result.(*GetRelayStateRsp)
	if !ok {
		t.Error("GetRelayStateRsp error")
	}

	future = p2pPID.RequestFuture(&GetVersionReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetVersionReq error %s", err)
	}
	_, ok = result.(*GetVersionRsp)
	if !ok {
		t.Error("GetVersionRsp error")
	}

	future = p2pPID.RequestFuture(&GetNodeTypeReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Errorf("GetNodeTypeReq error %s", err)
	}
	_, ok = result.(*GetNodeTypeRsp)
	if !ok {
		t.Error("GetNodeTypeRsp error")
	}

	future = p2pPID.RequestFuture(&StopServerReq{}, common.ACTOR_TIMEOUT*time.Second)
	result, err = future.Result()
	if err != nil {
		t.Fatalf("TestP2PActorServer: p2p halt error %s", err)
	}

}
