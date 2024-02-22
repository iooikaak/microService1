package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/iooikaak/microService1/model/enum"

	"github.com/iooikaak/frame/stat/metric/metrics"
	"github.com/iooikaak/frame/xlog"
	"github.com/iooikaak/microService1/config"
	"github.com/iooikaak/microService1/handler"
	pb "github.com/iooikaak/pb/microService1/http"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//远程读取apollo配置文件
	//初始化apollo
	//if err := config.ApolloInit(); err != nil {
	//	panic(err)
	//}
	config.Init()
	tracer, closer, err := config.CreateTracer(enum.ServiceName.String())
	defer closer.Close()
	if err != nil {
		panic(err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan(enum.JaegerStartSpan.String())
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	//metrics
	go func() {
		http.Handle("/microService1/metrics", promhttp.Handler())
		var m *metrics.Metrics
		var h http.Handler
		var once = sync.Once{}
		once.Do(
			func() {
				m = metrics.New(enum.ServiceName.String())
			},
		)
		once.Do(
			func() {
				h = promhttp.InstrumentMetricHandler(m.RegInstance(), promhttp.HandlerFor(m.Gather(), promhttp.HandlerOpts{}))
			},
		)
		//h := promhttp.InstrumentMetricHandler(m.RegInstance(), promhttp.HandlerFor(m.Gather(), promhttp.HandlerOpts{}))
		err = http.ListenAndServe(":3202", h)
		if err != nil {
			xlog.Fatal("ListenAndServe: ", err.Error())
			return
		}
	}()
	//注册到frame
	server := handler.NewBaseHandler(ctx, tracer)
	pb.RegisterMicroService1Server(server, config.Conf.BaseCfg)
}
