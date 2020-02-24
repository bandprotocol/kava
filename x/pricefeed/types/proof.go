package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func merkleLeafHash(leaf []byte) []byte {
	return merkle.SimpleHashFromByteSlices([][]byte{
		leaf,
	})
}

func merkleInnerHash(left, right []byte) []byte {
	return merkle.SimpleHashFromByteSlices([][]byte{
		left, right,
	})
}

func uint64ToBytes(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	// Encode uint64 must not have error
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type IAVLMerklePath struct {
	IsDataOnRight  bool         `json:"isDataOnRight"`
	SubtreeHeight  uint8        `json:"subtreeHeight"`
	SubtreeSize    uint64       `json:"subtreeSize"`
	SubtreeVersion uint64       `json:"subtreeVersion"`
	SiblingHash    cmn.HexBytes `json:"siblingHash"`
}

func (path *IAVLMerklePath) GetParentHash(currentMerkleHash []byte) []byte {
	letfSubtree := currentMerkleHash
	rightSubtree := path.SiblingHash
	if path.IsDataOnRight {
		letfSubtree, rightSubtree = rightSubtree, letfSubtree
	}
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(path.SubtreeSize))
	data := append([]byte{byte(path.SubtreeHeight * 2)}, buf[:n]...)
	n = binary.PutVarint(buf, int64(path.SubtreeVersion))
	data = append(data, buf[:n]...)
	data = append(data, byte(32))
	data = append(data, letfSubtree...)
	data = append(data, byte(32))
	data = append(data, rightSubtree...)
	return tmhash.Sum(data)
}

type BlockHeaderMerkleParts struct {
	VersionAndChainIdHash       cmn.HexBytes `json:"versionAndChainIdHash"`
	TimeHash                    cmn.HexBytes `json:"timeHash"`
	TxCountAndLastBlockInfoHash cmn.HexBytes `json:"txCountAndLastBlockInfoHash"`
	ConsensusDataHash           cmn.HexBytes `json:"consensusDataHash"`
	LastResultsHash             cmn.HexBytes `json:"lastResultsHash"`
	EvidenceAndProposerHash     cmn.HexBytes `json:"evidenceAndProposerHash"`
}

func (header *BlockHeaderMerkleParts) GetBlockHeader(
	blockHeight uint64, appHash []byte,
) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, blockHeight)
	return merkleInnerHash( // [BlockHeader]
		merkleInnerHash( // [3A]
			merkleInnerHash( // [2A]
				header.VersionAndChainIdHash, // [1A]
				merkleInnerHash( // [1B]
					merkleLeafHash(buf[:n]), // [2]
					header.TimeHash)),       // [3]
			header.TxCountAndLastBlockInfoHash), // [2B]
		merkleInnerHash( // [3B]
			header.ConsensusDataHash, // [2C]
			merkleInnerHash( // [2D]
				merkleInnerHash( // [1G]
					merkleLeafHash(append([]byte{byte(32)}, appHash...)), // [C]
					header.LastResultsHash),                              // [D]
				header.EvidenceAndProposerHash))) // [1H]
}

type BlockRelayProof struct {
	OracleIAVLStateHash    cmn.HexBytes           `json:"oracleIAVLStateHash"`
	OtherStoresMerkleHash  cmn.HexBytes           `json:"otherStoresMerkleHash"`
	SupplyStoresMerkleHash cmn.HexBytes           `json:"supplyStoresMerkleHash"`
	BlockHeaderMerkleParts BlockHeaderMerkleParts `json:"blockHeaderMerkleParts"`
	SignedDataPrefix       cmn.HexBytes           `json:"signedDataPrefix"`
	Signatures             []TMSignature          `json:"signatures"`
}

func (blockProof *BlockRelayProof) GetAppHash() []byte {
	zoracleStorePrefix, _ := hex.DecodeString("077a6f7261636c6520")
	return merkleInnerHash(
		blockProof.OtherStoresMerkleHash,
		merkleInnerHash(
			blockProof.SupplyStoresMerkleHash,
			merkleLeafHash(
				append(zoracleStorePrefix, tmhash.Sum(tmhash.Sum(blockProof.OracleIAVLStateHash))...),
			),
		),
	)
}

type OracleDataProof struct {
	Version     uint64           `json:"version"`
	RequestID   uint64           `json:"requestID"`
	CodeHash    cmn.HexBytes     `json:"codeHash"`
	Params      cmn.HexBytes     `json:"params"`
	Data        cmn.HexBytes     `json:"data"`
	MerklePaths []IAVLMerklePath `json:"merklePaths"`
}

func (proof *OracleDataProof) GetOracleStateRoot(blockHeight uint64) []byte {
	merkleHash := append([]byte{byte(0)}, byte(2))
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(proof.Version))
	merkleHash = append(merkleHash, buf[:n]...)
	merkleHash = append(merkleHash, byte(41+len(proof.Params)))
	merkleHash = append(merkleHash, byte(255))
	merkleHash = append(merkleHash, uint64ToBytes(proof.RequestID)...)
	merkleHash = append(merkleHash, proof.CodeHash...)
	merkleHash = append(merkleHash, proof.Params...)
	merkleHash = append(merkleHash, byte(32))
	merkleHash = append(merkleHash, tmhash.Sum(proof.Data)...)

	merkleHash = tmhash.Sum(merkleHash)

	for idx := range proof.MerklePaths {
		merkleHash = proof.MerklePaths[idx].GetParentHash(merkleHash)
	}
	return merkleHash
}

type TMSignature struct {
	R                cmn.HexBytes `json:"r"`
	S                cmn.HexBytes `json:"s"`
	V                uint8        `json:"v"`
	SignedDataSuffix cmn.HexBytes `json:"signedDataSuffix"`
}

func (signature *TMSignature) VerifyWithPub(
	pubKey *secp256k1.PubKeySecp256k1, blockHeader, signedDataPrefix []byte,
) bool {
	sig := append(signature.R, signature.S...)
	sig = append(sig, byte(signature.V))
	msg := append(signedDataPrefix, blockHeader...)
	msg = append(msg, signature.SignedDataSuffix...)
	return pubKey.VerifyBytes(msg, sig)
}

type Proof struct {
	BlockHeight     uint64          `json:"blockHeight"`
	OracleDataProof OracleDataProof `json:"oracleDataProof"`
	BlockRelayProof BlockRelayProof `json:"blockRelayProof"`
}
