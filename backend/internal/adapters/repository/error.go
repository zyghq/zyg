package repository

type dbErr string

func (err dbErr) Error() string {
	return string(err)
}

const (
	ErrEmpty   = dbErr("got nothing")
	ErrQuery   = dbErr("db query failed")
	ErrTxQuery = dbErr("db tx query failed")
)
