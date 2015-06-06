package obscurate

import (
	"fmt"
)

// Obscurates a chunk of data using the key
func (key *Key) Obscurate(data []byte) []byte {
	ret := make([]byte, len(data)+256)
	pre := make([]byte, 256)
	for k := range pre {
		pre[k] = byte(k)
	}
	data = append(pre, data...)
	for k, v := range data {
		for _, op := range key.Ops {
			var value byte
			value = 0
			if op.Constant {
				value = op.Value
			} else {
				pos := k - int(op.Value+1)
				if pos >= 0 {
					value = data[pos]
				}
			}
			switch op.Type {
			case OpXOR:
				v ^= value
			case OpAdd:
				v += value
			case OpSub:
				v -= value
			case OpSHL:
				value %= 8
				v = v<<value | v>>(8-value)
			case OpSHR:
				value %= 8
				v = v>>value | v<<(8-value)
			}
		}
		ret[k] = v
	}
	return ret[256:]
}

// Generates code for a function that will obscurate a []byte using the key
func (key *Key) ObscurateFunc(name string) string {
	code := "func " + name + `(data []byte) []byte {
	ret := make([]byte, len(data)+256)
	pre := make([]byte, 256)
	for k := range pre {
		pre[k] = byte(k)
	}
	data = append(pre, data...)
	for k, v := range data {
`
	for _, op := range key.Ops {
		var value string
		if op.Constant {
			value = fmt.Sprintf("%d", op.Value)
		} else {
			value = fmt.Sprintf(`data[k - %d]`, op.Value+1)
			code += fmt.Sprintf(`		if k >= %d {
`, op.Value+1)
		}
		switch op.Type {
		case OpXOR:
			code = code + fmt.Sprintf(
				`		v ^= %s
`,
				value)
		case OpAdd:
			code = code + fmt.Sprintf(
				`		v += %s
`,
				value)
		case OpSub:
			code = code + fmt.Sprintf(
				`		v -= %s
`,
				value)
		case OpSHL:
			code = code + fmt.Sprintf(
				`		v = (v << (%s %% 8)) | (v >> (8- (%s %% 8 )))
`,
				value, value)
		case OpSHR:
			code = code + fmt.Sprintf(
				`		v = (v >> (%s %% 8)) | (v << (8- (%s %% 8 )))
`,
				value, value)
		}
		if !op.Constant {
			code = code + `		}
`
		}

	}
	code += `		ret[k] = v
	}
	return ret[256:]
}`
	return code
}
