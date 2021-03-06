

package dbft

import (
	"fmt"
	"io"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/log"
	ser "github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/types"
)

type PrepareRequest struct {
	msgData        ConsensusMessageData
	Nonce          uint64
	NextBookkeeper common.Address
	Transactions   []*types.Transaction
	Signature      []byte
}

func (pr *PrepareRequest) Serialize(w io.Writer) error {
	log.Debug()

	pr.msgData.Serialize(w)
	if err := ser.WriteVarUint(w, pr.Nonce); err != nil {
		return fmt.Errorf("[PrepareRequest] nonce serialization failed: %s", err)
	}
	if err := pr.NextBookkeeper.Serialize(w); err != nil {
		return fmt.Errorf("[PrepareRequest] nextbookkeeper serialization failed: %s", err)
	}
	if err := ser.WriteVarUint(w, uint64(len(pr.Transactions))); err != nil {
		return fmt.Errorf("[PrepareRequest] length serialization failed: %s", err)
	}
	for _, t := range pr.Transactions {
		if err := t.Serialize(w); err != nil {
			return fmt.Errorf("[PrepareRequest] transactions serialization failed: %s", err)
		}
	}
	if err := ser.WriteVarBytes(w, pr.Signature); err != nil {
		return fmt.Errorf("[PrepareRequest] signature serialization failed: %s", err)
	}
	return nil
}

func (pr *PrepareRequest) Deserialize(r io.Reader) error {
	pr.msgData = ConsensusMessageData{}
	pr.msgData.Deserialize(r)
	pr.Nonce, _ = ser.ReadVarUint(r, 0)

	if err := pr.NextBookkeeper.Deserialize(r); err != nil {
		return fmt.Errorf("[PrepareRequest] nextbookkeeper deserialization failed: %s", err)
	}

	length, err := ser.ReadVarUint(r, 0)
	if err != nil {
		return fmt.Errorf("[PrepareRequest] length deserialization failed: %s", err)
	}

	pr.Transactions = make([]*types.Transaction, length)
	for i := 0; i < len(pr.Transactions); i++ {
		var t types.Transaction
		if err := t.Deserialize(r); err != nil {
			return fmt.Errorf("[PrepareRequest] transactions deserialization failed: %s", err)
		}
		pr.Transactions[i] = &t
	}

	pr.Signature, err = ser.ReadVarBytes(r)
	if err != nil {
		return fmt.Errorf("[PrepareRequest] signature deserialization failed: %s", err)
	}

	return nil
}

func (pr *PrepareRequest) Type() ConsensusMessageType {
	log.Debug()
	return pr.ConsensusMessageData().Type
}

func (pr *PrepareRequest) ViewNumber() byte {
	log.Debug()
	return pr.msgData.ViewNumber
}

func (pr *PrepareRequest) ConsensusMessageData() *ConsensusMessageData {
	log.Debug()
	return &(pr.msgData)
}
