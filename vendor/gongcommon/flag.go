package gongcommon

import (
	"flag"
)

type SubCmd struct {
	fs         *flag.FlagSet
	Name       string
	HelperText string
	Parse      func(*flag.FlagSet) error
}

func (c *SubCmd) Init() {
	c.fs = flag.NewFlagSet(c.Name, flag.ExitOnError)
}

func (c *SubCmd) RunParse() error {
	return c.Parse(c.fs)
}