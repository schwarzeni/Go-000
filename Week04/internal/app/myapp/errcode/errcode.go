package errcode

type ErrCode int

func (e ErrCode) Error() string {
	return codemap[e]
}

const (
	ErrDB ErrCode = iota
)

var codemap = map[ErrCode]string{
	ErrDB: "database error",
}
