package main

type Msg struct {
	Code string   `json:"code"`
	Args []string `json:"args,omitempty"`
}

func msg(code string, args ...string) Msg { return Msg{Code: code, Args: args} }

type codedError struct{ m Msg }

func (e codedError) Error() string { return e.m.Code }
func (e codedError) Msg() Msg      { return e.m }

func errMsg(code string, args ...string) error { return codedError{msg(code, args...)} }

func toMsg(err error) Msg {
	if ce, ok := err.(codedError); ok {
		return ce.m
	}
	return msg("generic", err.Error())
}
