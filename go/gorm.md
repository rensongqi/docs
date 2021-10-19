## 1 gorm的Related和Association的区别
二者区别在于谁主谁从
下述代码中user表是源，company关联源中的字段名，简而言之通过先查user-->commpany 最终的主表记录从user出发
```go
    db.Model(&user).Association("company").Find(&company)
```
而Related方法其内的company表是要查主表源，主查company表user实例只是条件填充对象

```go
db.Model(&user).Related(&company)
//// SELECT * FROM company WHERE user_id = 1; // 1 is user's primary key
```


开启sql打印一目了然

```go
// 启用Logger，显示详细日志
db.LogMode(true)
```

## 2 gorm Preload
默认情况下GORM因为性能问题，不会自动加载关联属性的值，gorm通过Preload函数支持预加载（Eager loading）关联数据
```go
// 用户表
type User struct {
  gorm.Model
  Username string
  Orders []Orders // 关联订单，一对多关联关系
}
// 订单表
type Orders struct {
  gorm.Model
  UserID uint // 外键字段 
  Price float64
}

// 预加载Orders字段值，Orders字段是User的关联字段
db.Preload("Orders").Find(&users)
// 下面是自动生成的SQL，自动完成关联查询
//// SELECT * FROM users;
//// SELECT * FROM orders WHERE user_id IN (1,2,3,4);


// Preload第2，3个参数支持设置SQL语句条件和绑定参数
db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
// 自动生成的SQL如下
//// SELECT * FROM users;
//// SELECT * FROM orders WHERE user_id IN (1,2,3,4) AND state NOT IN ('cancelled');

// 通过组合Where函数一起设置SQL条件
db.Where("state = ?", "active").Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
// 自动生成的SQL如下
//// SELECT * FROM users WHERE state = 'active';
//// SELECT * FROM orders WHERE user_id IN (1,2) AND state NOT IN ('cancelled');

// 预加载Orders、Profile、Role多个关联属性
// ps: 预加载字段，必须是User的属性
db.Preload("Orders").Preload("Profile").Preload("Role").Find(&users)
//// SELECT * FROM users;
//// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
//// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
//// SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
```
**自动预加载**

```go
type User struct {
  gorm.Model
  Name       string
  CompanyID  uint
  Company    Company `gorm:"PRELOAD:false"` // 通过标签属性关闭预加载
  Role       Role                           // 默认开启预加载特性
}
// 通过Set设置gorm:auto_preload属性，开启自动预加载，查询的时候才会自动完成关联查询
db.Set("gorm:auto_preload", true).Find(&users)
```
**嵌套预加载**
```go
// 预加载User.Orders.OrderItems属性值，使用点连接嵌套属性即可
db.Preload("Orders.OrderItems").Find(&users)
db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)
```