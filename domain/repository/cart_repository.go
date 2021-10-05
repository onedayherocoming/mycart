package repository

import (
	"errors"
	gorm "github.com/jinzhu/gorm"
	"github.com/onedayherocoming/mycart/domain/model"
)


type ICartRepository interface {
	InitTable()error
	FindCartByID(int64)(*model.Cart,error)
	CreateCart(Cart *model.Cart)(int64,error)
	DeleteCartByID(int64)error
	UpdateCart(Cart *model.Cart)error
	FindAll(int64 )([]model.Cart,error)

	CleanCart(int64) error
	IncrNum(int64,int64)error
	DecrNum(int64, int64)error
}
//创建CartRepository
func NewCartRepository(db *gorm.DB) ICartRepository{
	return &CartRepository{mysqlDb: db}
}


type CartRepository struct {
	mysqlDb *gorm.DB
}
func (u *CartRepository)InitTable()error{
	//创建多张表
	return u.mysqlDb.CreateTable(&model.Cart{}).Error
}

//根据CartID查找
func (u *CartRepository) FindCartByID(ID int64)(*model.Cart,error){
	Cart := &model.Cart{}
	return Cart,u.mysqlDb.First(Cart,ID).Error
}

//创建Cart
func (u *CartRepository)CreateCart(Cart *model.Cart)(int64,error){
	//有则不创建，没有就创建
	db := u.mysqlDb.FirstOrCreate(Cart,model.Cart{ProductID: Cart.ProductID,
		SizeID: Cart.SizeID,UserID: Cart.UserID})
	if db.Error!=nil{
		return 0,db.Error
	}
	if db.RowsAffected==0{
		return 0,errors.New("购物车插入失败")
	}
	return Cart.ID,nil
}

//根据CartID删除Cart
func (u *CartRepository)DeleteCartByID(ID int64)error{
	return u.mysqlDb.Where("id=?",ID).Delete(&model.Cart{}).Error
}

//更新Cart信息
func (u *CartRepository)UpdateCart(Cart *model.Cart)error{
	return u.mysqlDb.Model(Cart).Update(Cart).Error
}
//查找所有
func (u *CartRepository)FindAll(userID int64)(CartAll []model.Cart,err error){
	return CartAll,u.mysqlDb.Where("user_id=?",userID).Find(&CartAll).Error
}


func (u *CartRepository)CleanCart(userID int64) error{
	return u.mysqlDb.Where("user_id=?",userID).Delete(model.Cart{}).Error

}
//添加商品数量
func (u *CartRepository)IncrNum(cartID int64,num int64)error{
	cart := &model.Cart{ID:cartID}
	return u.mysqlDb.Model(cart).UpdateColumn("num",gorm.Expr("num+?",num)).Error
}
//减少商品
func (u *CartRepository)DecrNum(cartID int64, num int64)error{
	cart := &model.Cart{ID: cartID}
	db := u.mysqlDb.Model(cart).Where("num>=?",num).UpdateColumn("num",gorm.Expr("num-?",num))
	if db.Error!=nil{
		return db.Error
	}
	if db.RowsAffected==0{
		return errors.New("减少失败")
	}
	return nil
}