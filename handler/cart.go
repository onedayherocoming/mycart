package handler

import (
	"context"
	"github.com/onedayherocoming/mycart/domain/model"
	"github.com/onedayherocoming/mycart/domain/service"
	cart "github.com/onedayherocoming/mycart/proto/cart"
	common "github.com/onedayherocoming/mycommon"
)

type Cart struct{
	CartDataService service.IcartDataService
}

func (c *Cart)AddCart(ctx context.Context,request *cart.CartInfo,response *cart.ResponseAdd) error{
	cart := &model.Cart{}
	if err:=common.SwapTo(cart,request);err!=nil{
		return err
	}
	cartID,err:=c.CartDataService.AddCart(cart)
	if err!=nil{
		return err
	}
	response.CartId=cartID
	response.Message="添加成功"
	return nil
}

func (c *Cart)CleanCart(ctx context.Context,request  *cart.Clean,response  *cart.Response) error{
	if err:= c.CartDataService.CleanCart(request.UserId);err!=nil{
		return err
	}
	response.Meg="清空成功"
	return nil
}

func (c *Cart)Incr(ctx context.Context, request *cart.Item, response *cart.Response) error{
	if err:=c.CartDataService.IncrNum(request.Id,request.ChangeNum);err!=nil{
		return err
	}
	response.Meg="数量增加成功"
	return nil
}

func (c *Cart)Decr(ctx context.Context,request  *cart.Item,response  *cart.Response) error{
	if err:=c.CartDataService.DecrNum(request.Id,request.ChangeNum);err!=nil{
		return err
	}
	response.Meg="数量减少成功"
	return nil
}
func (c *Cart)DeleteItemByID(ctx context.Context,request  *cart.CartID, response *cart.Response) error{
	if err:=c.CartDataService.DeleteCart(request.Id);err!=nil{
		return err
	}
	response.Meg="删除成功"
	return nil
}
func (c *Cart)GetAll(ctx context.Context, request *cart.CartFindAll,response  *cart.CartAll) error{
	allCarts,err:=c.CartDataService.FindAllCart(request.UserUd)
	if err!=nil{
		return err
	}
	for i := range allCarts{
		carInfo := &cart.CartInfo{}
		if err=common.SwapTo(i,carInfo);err!=nil{
			return err
		}
		response.CartInfo = append(response.CartInfo,carInfo)
	}
	return nil
}

