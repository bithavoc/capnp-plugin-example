package common

import (
	"io"
)

type StdStreamJoint struct {
	in     io.Reader
	out    io.Writer
	closed bool
}

func NewStdStreamJoint(in io.Reader, out io.Writer) *StdStreamJoint {
	return &StdStreamJoint{
		in:  in,
		out: out,
	}
}

func (s *StdStreamJoint) Read(b []byte) (n int, err error) {
	return s.in.Read(b)
}

func (s *StdStreamJoint) Write(b []byte) (n int, err error) {
	return s.out.Write(b)
}

func (s *StdStreamJoint) Close() error {
	// panic(fmt.Errorf("StdStreamJoint closed"))
	s.closed = true
	return nil
}
