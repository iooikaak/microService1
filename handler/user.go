package handler

import (
	frame "github.com/iooikaak/frame/core"
	"github.com/iooikaak/microService1/config"
	"github.com/iooikaak/microService1/service/user"

	pbms1 "github.com/iooikaak/pb/microService1/http"

	"github.com/iooikaak/frame/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openLog "github.com/opentracing/opentracing-go/log"
)

var (
	userHandler *UserHandler
)

type UserHandler struct {
	*BaseHandler
	*user.UserService
	User
}

type User interface {
	GetUserInfo(c frame.Context) error
}

func NewUserHandler(b *BaseHandler) *UserHandler {
	userHandler = &UserHandler{
		BaseHandler: b,
		UserService: user.New(config.Conf),
	}
	return userHandler
}

func (h *UserHandler) GetUserInfo(c frame.Context) error {
	spanCtx, _ := h.Tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(c.Request().Header))
	span := h.Tracer.StartSpan("microService1 GetUserInfo() start tracing", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	var req pbms1.GetUserInfoReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	userInfo, err := userHandler.UserService.GetUserInfo(h.Context, &req)
	if err != nil {
		return c.JSON2(1, err.Error(), nil)
	}

	go func() {
		str, _ := json.Marshal(userInfo)
		span.LogFields(
			openLog.String("event", "microService1 GetUserInfo() get data from tidb"),
			openLog.String("value", string(str)),
		)
	}()
	obj := &pbms1.GetUserInfoResp_Data{
		Name: userInfo.Name,
		Age:  userInfo.Age,
		Job:  userInfo.Job,
	}
	if userInfo.Gender == 2 {
		obj.Gender = "女"
	} else {
		//默认男
		obj.Gender = "男"
	}
	go func() {
		b, _ := json.Marshal(obj)
		span.LogFields(
			openLog.String("event", "microService1 GetUserInfo() return data"),
			openLog.String("value", string(b)),
		)
	}()
	return c.JSON2(0, "success", obj)
}
