package service

import (
	"github.com/onedayherocoming/mycart/domain/model"
	"github.com/onedayherocoming/mycart/domain/repository"
)

type IcartDataService interface {
	AddCart(user *model.Cart)(int64,error)
	DeleteCart(int64)error
	UpdateCart(user *model.Cart)(err error)
	FindCartByID(int64)(*model.Cart,error)
	FindAllCart(int64)([]model.Cart, error)

	CleanCart(int64)error
	DecrNum(int64,int64)error
	IncrNum(int64,int64)error
}

//创建
func NewcartDataService(cartRepository repository.ICartRepository)IcartDataService{
	return &cartDataService{cartRepository: cartRepository}
}

type cartDataService struct {
	cartRepository repository.ICartRepository
}

func (u *cartDataService)AddCart(user *model.Cart)(int64,error){
	return u.cartRepository.CreateCart(user)
}
func (u *cartDataService)DeleteCart(ID int64)error{
	return u.cartRepository.DeleteCartByID(ID)
}

func (u *cartDataService)UpdateCart(user *model.Cart)(err error){
	return u.cartRepository.UpdateCart(user)
}

func (u *cartDataService)FindCartByID(ID int64)(*model.Cart,error){
	return u.cartRepository.FindCartByID(ID)
}
func (u *cartDataService) FindAllCart(userID int64)([]model.Cart, error){
	return u.cartRepository.FindAll(userID)
}

func (u *cartDataService)CleanCart(userID int64)error{
	return u.cartRepository.CleanCart(userID)
}

func (u *cartDataService) DecrNum(cartID int64,num int64)error{
	return u.cartRepository.DecrNum(cartID,num)
}

func (u *cartDataService) IncrNum(cartID int64,num int64)error{
	return u.cartRepository.IncrNum(cartID,num)
}