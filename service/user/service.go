package user

import (
	conf "github.com/iooikaak/microService1/config"
	"github.com/iooikaak/microService1/dao/mysql/user"
)

type UserService struct {
	//*service.BaseService
	db *user.Dao
}

func New(cfg *conf.Configuration) *UserService {
	srv := &UserService{
		//BaseService: service.New(cfg),
		db: user.New(cfg),
	}
	return srv
}
