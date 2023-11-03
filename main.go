package main

import "sync"

// Outlet：代表TechSolutions的一个出口。
type Outlet struct {
	ID      int     `json:"id"`      // 唯一标识符
	Name    string  `json:"name"`    // 出口名称
	Address string  `json:"address"` // 出口地址
	Offers  []Offer `json:"offers"`  // 当前出口的特别优惠
}

// Product：代表一个产品或配件。
type Product struct {
	ID          int     `json:"id"`           // 唯一标识符
	Name        string  `json:"name"`         // 产品名称
	Description string  `json:"description"`  // 产品描述
	Price       float64 `json:"price"`        // 产品价格
	IsAccessory bool    `json:"is_accessory"` // 该产品是否是配件
}

// Offer：代表特殊的优惠或促销。
type Offer struct {
	ProductID int     `json:"product_id"` // 参与优惠的产品ID
	Discount  float64 `json:"discount"`   // 折扣金额
}

// Purchase：代表客户的一个购买。
type Purchase struct {
	ID         int    `json:"id"`          // 唯一标识符
	CustomerID int    `json:"customer_id"` // 购买该产品的顾客ID
	ProductID  int    `json:"product_id"`  // 购买的产品ID
	OutletID   int    `json:"outlet_id"`   // 购买地点的出口ID
	Date       string `json:"date"`        // 购买日期
}

// LoyaltyCard：代表GadgetPoints的忠诚卡。
type LoyaltyCard struct {
	ID          int `json:"id"`          // 唯一标识符
	Points      int `json:"points"`      // 当前的购买次数
	Redemptions int `json:"redemptions"` // 可兑换次数
}

// SubscriptionService：代表通过订阅获得的服务。
type SubscriptionService struct {
	ID          int    `json:"id"`          // 唯一标识符
	Name        string `json:"name"`        // 服务名称，如"无限的设备服务"
	Description string `json:"description"` // 服务描述
}

// DeviceDiscount：代表订阅者的设备折扣。
type DeviceDiscount struct {
	ID             int     `json:"id"`              // 唯一标识符
	ProductID      int     `json:"product_id"`      // 受折扣影响的产品ID
	DiscountAmount float64 `json:"discount_amount"` // 折扣金额或百分比
	IsPercentage   bool    `json:"is_percentage"`   // 折扣是否是百分比形式
}

// Subscription：代表在线商店的订阅服务。
type Subscription struct {
	ID              int                   `json:"id"`               // 唯一标识符
	CustomerID      int                   `json:"customer_id"`      // 订阅服务的顾客ID
	StartDate       string                `json:"start_date"`       // 订阅开始日期
	EndDate         string                `json:"end_date"`         // 订阅结束日期
	IsActive        bool                  `json:"is_active"`        // 是否当前有效
	ServiceUsed     int                   `json:"service_used"`     // 已经使用的服务数量
	Services        []SubscriptionService `json:"services"`         // 通过此订阅获得的服务列表
	DeviceDiscounts []DeviceDiscount      `json:"device_discounts"` // 为订阅者提供的设备折扣
}

// Customer：代表一个顾客。
type Customer struct {
	ID            int            `json:"id"`            // 唯一标识符
	Name          string         `json:"name"`          // 顾客名称
	Email         string         `json:"email"`         // 顾客电子邮件
	LoyaltyCard   LoyaltyCard    `json:"loyalty_card"`  // 顾客的GadgetPoints忠诚卡
	Subscriptions []Subscription `json:"subscriptions"` // 顾客的在线商店订阅
	mu            sync.Mutex     // 互斥锁
}

// AccumulatePoints：累积购买次数和可兑换次数
func (c *Customer) AccumulatePoints(wg *sync.WaitGroup) {
	defer wg.Done() // 告诉WaitGroup一个goroutine已经完成

	c.mu.Lock()         // 获取锁
	defer c.mu.Unlock() // 释放锁

	c.LoyaltyCard.Points++
	if c.LoyaltyCard.Points%10 == 0 {
		c.LoyaltyCard.Redemptions++
	}
}

// UseRedemption：使用红利点数
func (c *Customer) UseRedemption(wg *sync.WaitGroup) {
	defer wg.Done() // 告诉WaitGroup一个goroutine已经完成

	c.mu.Lock()         // 获取锁
	defer c.mu.Unlock() // 释放锁

	if c.LoyaltyCard.Redemptions > 0 {
		c.LoyaltyCard.Redemptions--
	}
}
func main() {
	customer := Customer{}

	var wg sync.WaitGroup

	// 使用WaitGroup和goroutines进行点数累积和赎回操作示例
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go customer.AccumulatePoints(&wg)
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go customer.UseRedemption(&wg)
	}

	wg.Wait() // 等待所有goroutines完成

	// 打印最后的结果
	println("Points:", customer.LoyaltyCard.Points)
	println("Redemptions:", customer.LoyaltyCard.Redemptions)
}
