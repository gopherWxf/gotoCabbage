package core

import (
	"gotoCabbage/interface"
)

type BaseRouter struct{}

func (b *BaseRouter) PreHandle(iRequest _interface.IRequest) {
}

func (b *BaseRouter) Handle(iRequest _interface.IRequest) {
}

func (b *BaseRouter) PostHandle(iRequest _interface.IRequest) {
}
