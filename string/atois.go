package string

import "strconv"

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ParseUint(s string) (uint, error) {
	n, err := strconv.Atoi(s)
	return uint(n), err
}

func ParseInt8(s string) (int8, error) {
	n, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(n), nil
}

func ParseUint8(s string) (uint8, error) {
	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}

func ParseInt16(s string) (int16, error) {
	n, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}

func ParseUint16(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}

func ParseInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

func ParseUint32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}

func ParseInt64(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

func ParseUint64(s string) (uint64, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(n), nil
}

func ParseIntDefault(s string, def int) int {
	n, err := ParseInt(s)
	if err != nil {
		return def
	}
	return n
}

func ParseUintDefault(s string, def uint) uint {
	n, err := ParseUint(s)
	if err != nil {
		return def
	}
	return n
}

func ParseInt8Default(s string, def int8) int8 {
	n, err := ParseInt8(s)
	if err != nil {
		return def
	}
	return n
}

func ParseUint8Default(s string, def uint8) uint8 {
	n, err := ParseUint8(s)
	if err != nil {
		return def
	}
	return n
}

func ParseInt16Default(s string, def int16) int16 {
	n, err := ParseInt16(s)
	if err != nil {
		return def
	}
	return n
}

func ParseUint16Default(s string, def uint16) uint16 {
	n, err := ParseUint16(s)
	if err != nil {
		return def
	}
	return n
}

func ParseInt32Default(s string, def int32) int32 {
	n, err := ParseInt32(s)
	if err != nil {
		return def
	}
	return n
}

func ParseUint32Default(s string, def uint32) uint32 {
	n, err := ParseUint32(s)
	if err != nil {
		return def
	}
	return n
}

func ParseInt64Default(s string, def int64) int64 {
	n, err := ParseInt64(s)
	if err != nil {
		return def
	}
	return n
}

func ParseUint64Default(s string, def uint64) uint64 {
	n, err := ParseUint64(s)
	if err != nil {
		return def
	}
	return n
}
