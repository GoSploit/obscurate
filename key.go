//go:generate jsonenums -type OpType
//go:generate stringer -type OpType

package obscurate

import (
	"crypto/rand"
)

type Key struct {
	Ops []Operation
}

type Operation struct {
	Type     OpType
	Constant bool
	Value    byte
}

type OpType byte

const (
	OpXOR OpType = iota
	OpAdd
	OpSub
	OpSHL
	OpSHR
)

// Generates a "key" of a specified number of operations
func GenerateKey(Ops int) (*Key, error) {
	rawkey := make([]byte, Ops*2)
	_, err := rand.Read(rawkey)
	if err != nil {
		return nil, err
	}
	ops := make([]Operation, 0, Ops)
	for l1 := 0; l1 < Ops; l1++ {
		constant := (rawkey[l1*2] & 1) == 1
		opType := rawkey[l1*2] >> 1 & 0x07
		opType = opType % 5
		val := rawkey[l1*2+1]
		ops = append(ops, Operation{
			Type:     OpType(opType),
			Constant: constant,
			Value:    val,
		})
	}
	key := &Key{
		Ops: ops,
	}
	return key, nil
}
