package show

import (
	"context"
	"github.com/iooikaak/frame/json"
	"github.com/iooikaak/microService1/config"
	"github.com/iooikaak/microService1/model/enum"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	//_ = flag.Set("conf", "../../build/config.yaml")
	//err := conf.ApolloInit()
	//if err != nil {
	//	panic(err)
	//}
	Init()
	d = New(config.Conf)
	os.Exit(m.Run())
}

func Init() (err error) {
	var (
		jsonFile, confPath string
	)
	if confPath != "" {
		jsonFile, err = filepath.Abs(confPath)
	} else {
		jsonFile, err = filepath.Abs(enum.MicroService1JsonDaoPath.String())
	}
	if err != nil {
		return
	}
	jsonRead, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonRead, config.Conf)
	if err != nil {
		return
	}
	return
}
