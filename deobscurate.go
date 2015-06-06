package obscurate

import (
	"fmt"
)

// Deobscurates a chunk of data using the key
func (key *Key) Deobscurate(data []byte) []byte {
	ret := make([]byte, len(data)+256)
	revops := make([]Operation, len(key.Ops))
	for k, op := range key.Ops {
		revops[len(key.Ops)-(k+1)] = op
	}

	for k := 0; k < len(data)+256; k++ {
		if k < 256 {
			ret[k] = byte(k)
			continue
		}
		v := data[k-256]
		for _, op := range revops {
			var value byte
			value = 0
			if op.Constant {
				value = op.Value
			} else {
				pos := k - int(op.Value+1)
				value = ret[pos]
			}
			switch op.Type {
			case OpXOR:
				v ^= value
			case OpAdd:
				v -= value
			case OpSub:
				v += value
			case OpSHL:
				value %= 8
				v = v>>value | v<<(8-value)
			case OpSHR:
				value %= 8
				v = v<<value | v>>(8-value)
			}
		}
		ret[k] = v
	}
	return ret[256:]
}

// Generates code for a function that will deobscurate a []byte using the key
func (key *Key) DeobscurateFunc(name string) string {
	code := "func " + name + `(data []byte) []byte {
	ret := make([]byte, len(data)+256)
	for k := 0; k < len(data)+256; k++ {
		if k < 256 {
			ret[k] = byte(k)
			continue
		}
		v := data[k-256]
`
	opscode := ""
	for _, op := range key.Ops {
		var value string
		if op.Constant {
			value = fmt.Sprintf("%d", op.Value)
		} else {
			value = fmt.Sprintf(`ret[k - %d]`, op.Value+1)
		}
		switch op.Type {
		case OpXOR:
			opscode = fmt.Sprintf(
				`		v ^= %s
`,
				value) + opscode
		case OpAdd:
			opscode = fmt.Sprintf(
				`		v -= %s
`,
				value) + opscode
		case OpSub:
			opscode = fmt.Sprintf(
				`		v += %s
`,
				value) + opscode
		case OpSHL:
			opscode = fmt.Sprintf(
				`		v = (v >> (%s %% 8)) | (v << (8- (%s %% 8 )))
`,
				value, value) + opscode
		case OpSHR:
			opscode = fmt.Sprintf(
				`		v = (v << (%s %% 8)) | (v >> (8- (%s %% 8 )))
`,
				value, value) + opscode
		}

	}
	code += opscode
	code += `		ret[k] = v
	}
	return ret[256:]
}`
	return code
}
