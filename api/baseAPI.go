package api

//import (
//	"sync"
//
//	"github.com/iooikaak/microService1/config"
//	"github.com/iooikaak/microService1/dao"
//)
//
//type BaseAPI struct {
//	*config.Configuration
//	*dao.BaseDao
//}
//
//var (
//	api  *BaseAPI
//	once sync.Once
//)
//
//func NewBaseAPI(config *config.Configuration) *BaseAPI {
//	once.Do(func() {
//		api = &BaseAPI{
//			Configuration: config,
//			BaseDao:       dao.NewBaseDao(config),
//		}
//	})
//	return api
//}
