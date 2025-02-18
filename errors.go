package main

import "errors"

var (
	ErrNaN = errors.New("not a number")
	Blown  = errors.New("Sonic Dead")
	ErrLol = errors.New("for strings numbers, use NewParsedTokenFromBytes")
)
