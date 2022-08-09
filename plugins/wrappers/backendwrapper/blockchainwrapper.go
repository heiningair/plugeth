package backendwrapper

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	gcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/openrelayxyz/plugeth-utils/core"
)

type Blockchain struct {
	bc *gcore.BlockChain
}

func NewBlockchain(bc *gcore.BlockChain) core.BlockChain {
	return &Blockchain{bc}
}

func (b *Blockchain) SetHead(head uint64) error {
	return b.bc.SetHead(head)
}

func (b *Blockchain) SnapSyncCommitHead(hash core.Hash) error {
	return b.bc.SnapSyncCommitHead(common.Hash(hash))
}

func (b *Blockchain) Reset() error {
	return b.bc.Reset()
}

// RLP encoded *types.Block
func (b *Blockchain) ResetWithGenesisBlock(genesis []byte) error {
	block := &types.Block{}
	err := rlp.DecodeBytes(genesis, block)
	if err != nil {
		return err
	}
	return b.bc.ResetWithGenesisBlock(block)
}

func (b *Blockchain) Export(w io.Writer) error {
	return b.bc.Export(w)
}

func (b *Blockchain) ExportN(w io.Writer, first uint64, last uint64) error {
	return b.bc.ExportN(w, first, last)
}

func (b *Blockchain) Stop() {
	b.bc.Stop()
}

func (b *Blockchain) StopInsert() {
	b.bc.StopInsert()
}

// List of RLP encoded blocks, list of RLP encoded receipts
func (b *Blockchain) InsertReceiptChain(blockChain [][]byte, receiptChain [][][]byte, ancientLimit uint64) (int, error) {
	blocks := make([]*types.Block, len(blockChain))
	for i, block := range blockChain {
		new := &types.Block{}
		err := rlp.DecodeBytes(block, new)
		if err != nil {
			return 0, err
		}
		blocks[i] = new
	}
	receiptLists := make([]types.Receipts, len(receiptChain))
	for i, receipts := range receiptChain {
		receiptLists[i] = make(types.Receipts, len(receipts))
		for j, receipt := range receipts {
			new := &types.Receipt{}
			err := rlp.DecodeBytes(receipt, new)
			if err != nil {
				return 0, err
			}
			receiptLists[i][j] = new
		}
	}
	return b.bc.InsertReceiptChain(blocks, receiptLists, ancientLimit)
}

// List of RLP encoded blocks
func (b *Blockchain) InsertChain(chain [][]byte) (int, error) {
	blocks := make([]*types.Block, len(chain))
	for i, block := range chain {
		new := &types.Block{}
		err := rlp.DecodeBytes(block, new)
		if err != nil {
			return 0, err
		}
		blocks[i] = new
	}
	return b.bc.InsertChain(blocks)
}

// RLP encoded *types.BLock
func (b *Blockchain) InsertBlockWithoutSetHead(block []byte) error {
	new := &types.Block{}
	err := rlp.DecodeBytes(block, new)
	if err != nil {
		return err
	}
	return b.bc.InsertBlockWithoutSetHead(new)
}

// RLP encoded *types.Block
func (b *Blockchain) SetChainHead(head []byte) error {
	new := &types.Block{}
	err := rlp.DecodeBytes(head, new)
	if err != nil {
		return err
	}
	return b.bc.SetHead(new.NumberU64())
}

// List of rlp encoded *types.Header
func (b *Blockchain) InsertHeaderChain(chain [][]byte, checkFreq int) (int, error) {
	headers := make([]*types.Header, len(chain))
	for i, header := range chain {
		new := &types.Header{}
		err := rlp.DecodeBytes(header, new)
		if err != nil {
			return 0, err
		}
		headers[i] = new
	}
	return b.bc.InsertHeaderChain(headers, checkFreq)
}
