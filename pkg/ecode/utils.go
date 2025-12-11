package ecode

import (
	"errors"
	"fmt"
	"runtime"
)

func WithCaller(err error) error {
	return WithCallerAndSkip(err, 2)
}

func WithCallerAndSkip(err error, skip int) error {
	pc, file, no, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		caller := fmt.Errorf("caller: func:%s file:%s line:%d", details.Name(), file, no)
		if eErr := FromError(err); eErr != nil && eErr.Code != UnknownCode {
			return eErr.WithCause(caller)
		}
		return errors.Join(err, caller)
	}
	return err
}

func JoinStr(err error, msg string) error {
	if err == nil {
		return nil
	}
	if eErr := FromError(err); eErr != nil {
		if eErr.cause != nil {
			eErr.cause = errors.Join(eErr.cause, err)
			return eErr
		}
		return eErr.WithCause(errors.New(msg))
	}
	return errors.Join(err, errors.New(msg))
}
