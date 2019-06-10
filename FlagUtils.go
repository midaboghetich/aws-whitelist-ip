package main

import "fmt"

type flagArray []string

func (i *flagArray) String() string {
	return fmt.Sprintf("%d", i)
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}
