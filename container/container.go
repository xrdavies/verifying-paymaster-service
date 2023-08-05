package container

import (
	"github.com/ququzone/verifying-paymaster-service/db"
)

type Container interface {
	GetRepository() db.Repository
}

func NewContainer(rep db.Repository) Container {
	return &container{
		rep: rep,
	}
}

type container struct {
	rep db.Repository
}

func (c *container) GetRepository() db.Repository {
	return c.rep
}
