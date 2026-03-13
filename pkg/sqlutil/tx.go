package sqlutil

import (
	"context"
	"database/sql"

	pkgerr "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Bnei-Baruch/gxydb-api/pkg/errs"
)

type TxError struct {
	errs.WithMessage
}

func WrappingTxError(err error, msg string) *TxError {
	return &TxError{errs.WithMessage{
		Msg: msg,
		Err: err,
	}}
}

func InTx(ctx context.Context, beginner boil.Beginner, f func(*sql.Tx) error) error {
	var tx *sql.Tx
	var err error
	if ctxBeginner, ok := beginner.(boil.ContextBeginner); ok {
		tx, err = ctxBeginner.BeginTx(ctx, nil)
	} else {
		tx, err = beginner.Begin()
	}
	if err != nil {
		return pkgerr.WithStack(WrappingTxError(err, "begin tx"))
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		_ = tx.Rollback()
	}()

	if err := f(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return pkgerr.WithStack(WrappingTxError(err, "tx.Commit"))
	}

	return nil
}
