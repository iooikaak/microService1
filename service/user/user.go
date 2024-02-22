package user

import (
	"context"
	frame "github.com/iooikaak/frame/core"
	conf "github.com/iooikaak/microService1/config"
	"github.com/iooikaak/microService1/dao/rocketmq/test"

	model "github.com/iooikaak/microService1/database/mysql/user"

	"github.com/iooikaak/frame/json"
	pbms1 "github.com/iooikaak/pb/microService1/http"
	pbms2 "github.com/iooikaak/pb/microService2/http"
	"github.com/opentracing/opentracing-go"
	openLog "github.com/opentracing/opentracing-go/log"
)

func (s *UserService) GetUserInfo(ctx context.Context, req *pbms1.GetUserInfoReq) (*model.UserInfo, error) {
	res, err := s.db.GetUserInfo(ctx, req.UserId)
	if err := test.New(conf.Conf).SendToTest(ctx, "===== musheng =====", "", "testing"); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) RpcUserInfo(ctx context.Context, span opentracing.Span, userID int32) *pbms2.GetUserInfoRep_Data {
	resp, err := pbms2.GetUserInfo(span, frame.TODO(), &pbms2.GetUserInfoReq{
		UserId: userID,
	})
	if err != nil || resp.GetErrcode() != 0 {
		return nil
	}
	go func() {
		b, _ := json.Marshal(resp)
		span.LogFields(
			openLog.String("event", "microService1 rpc request microservice2 method getuserinfo() success and get data"),
			openLog.String("value", string(b)),
		)
	}()
	return resp.GetData()
}
