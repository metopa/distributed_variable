package variable

import (
	"errors"
	"fmt"

	"github.com/metopa/distributed_variable/net"
)

type VariableType int32

type Variable interface {
	Get() (VariableType, error)
	Set(VariableType) error
}

type LocalVariable struct {
	value VariableType
}

type RemoteVariable struct {
	ring net.Ring
}

func (v *LocalVariable) Get() (VariableType, error) {
	return v.value, nil
}

func (v *LocalVariable) Set(value VariableType) error {
	v.value = value
	return nil
}

func (v *RemoteVariable) Get() (VariableType, error) {
	return 0, errors.New("not implemented")
}

func (v *RemoteVariable) Set(value VariableType) error {
	return errors.New("not implemented")
}
