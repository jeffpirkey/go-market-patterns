package main

import (
	"errors"
	"fmt"
)

func predict(tsym, pattern string) error {

	p := Periods[tsym]
	if p == nil {
		return errors.New("ticker sym not found")
	}

	fmt.Println(p.Last())

	return nil
}
