package main

import (
	"fmt"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3/util/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/onedayherocoming/mycart/domain/repository"
	"github.com/onedayherocoming/mycart/handler"
	common "github.com/onedayherocoming/mycommon"
	"github.com/opentracing/opentracing-go"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	service2 "github.com/onedayherocoming/mycart/domain/service"
	proto "github.com/onedayherocoming/mycart/proto/cart"
)

var QPS=100

func main() {
	consulIp := "192.168.1.100"
	jaegerIp := "192.168.1.100"
	//配置中心
	consulConfig, err := common.GetConsulConfig(consulIp, 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulIp + ":8500",
		}
	})
	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", jaegerIp+":6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	//数据库
	//获取mysql配置信息，路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	fmt.Println(mysqlInfo)
	//初始化数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@tcp("+mysqlInfo.Host+")/"+mysqlInfo.Database+"?charset=utf8&parseTime=True")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止复表
	db.SingularTable(true)
	rp := repository.NewCartRepository(db)
	err=rp.InitTable()
	if err!=nil{
		log.Error(err)
	}

	// Create service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		//设置地址和需要暴露的端口
		micro.Address("127.0.0.1:8087"),
		//添加consul作为注册中心
		micro.Registry(consulRegistry),
		//链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),

	)

	//初始化服务
	service.Init()

	//创建category服务
	cartDataService := service2.NewcartDataService(rp)
	err = proto.RegisterCartHandler(service.Server(), &handler.Cart{cartDataService})
	if err != nil {
		log.Error(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
