package handler

import (
	"github.com/iooikaak/microService1/config"
	"github.com/iooikaak/microService1/service/user"

	pbms1 "github.com/iooikaak/pb/microService1/http"

	frame "github.com/iooikaak/frame/core"
	"github.com/iooikaak/frame/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openLog "github.com/opentracing/opentracing-go/log"
)

var (
	userInfoHandler *UserInfoHandler
)

type UserInfoHandler struct {
	*BaseHandler
	*user.UserService
	UserInfo
}

type UserInfo interface {
	RpcGetUserInfo(c frame.Context) error
	AddUserInfo(c frame.Context) error
}

func NewUserInfoHandler(b *BaseHandler) *UserInfoHandler {
	userInfoHandler = &UserInfoHandler{
		BaseHandler: b,
		UserService: user.New(config.Conf),
	}
	return userInfoHandler
}

func (h *UserInfoHandler) AddUserInfo(c frame.Context) error {
	return nil
}

func (h *UserInfoHandler) RpcGetUserInfo(c frame.Context) error {
	spanCtx, _ := h.Tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(c.Request().Header))
	span := h.Tracer.StartSpan("microService1 GetUserInfo() start tracing", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	var req pbms1.RpcGetUserInfoReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	userInfo := userHandler.UserService.RpcUserInfo(h.Context, span, req.UserId)
	go func() {
		b, _ := json.Marshal(userInfo)
		span.LogFields(
			openLog.String("event", "microService1 RpcGetUserInfo() return data"),
			openLog.String("value", string(b)),
		)
	}()
	return c.JSON2(0, "success", userInfo)
}
