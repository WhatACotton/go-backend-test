package database

import (
	"log"
)

type CartContent struct {

	//Auto Incrment
	Order int `json:"Order"`
	//From CartRequestPayload
	ItemID   string `json:"ItemID"`
	Quantity int    `json:"Quantity"`
	//From Item
	Status string `json:"Status"`
	//From Item
	ItemName string `json:"ItemName"`
	Price    int    `json:"Price"`
	Stock    int    `json:"Stock"`
}
type CartContents []CartContent

func (c *CartContentRequestPayload) Cart(CartID string) {
	log.Println("CartID : " + CartID)
	log.Print("ItemID : "+c.ItemID, " Quantity : ", c.Quantity)
	Item := new(Item)
	Item.ItemGet(c.ItemID)
	//リクエストされた商品が登録可能か判定
	log.Println("ItemStatus : " + Item.Status)
	if Item.Status == "Available" {
		Carts := GetCartContents(CartID)
		//リクエストされた商品がカートに存在するか確認
		//存在する場合
		if searchCartContents(Carts, c.ItemID) {
			if c.Quantity == 0 {
				DeleteCartContent(CartID, c.ItemID)
			} else {
				if Item.Stock >= c.Quantity {
					c.UpdateCartContent(CartID)
				} else {
					log.Print("stock is not enough")
				}
			}
		} else {
			//存在しない場合
			if c.Quantity != 0 {
				if Item.Stock >= c.Quantity {
					c.PostCartContent(CartID)
				} else {
					log.Print("stock is not enough")
				}

			} else {
				log.Print("CartReq Quantity is 0")
			}
		}
	}
}
func GetCartContents(CartID string) (CartContents CartContents) {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		Item.Status ,
		Item.Price , 
		Item.Name , 
		Item.Stock 
		CartContents.Order , 
		CartContents.ItemID , 
		CartContents.Quantity , 
	FROM 
		Item
	JOIN 
		CartContents 
	ON 
	CartContents.ItemID = Item.ItemID 
	WHERE 
		CartID = ?`, CartID)
	defer rows.Close()
	for rows.Next() {
		CartContent := new(CartContent)
		rows.Scan(&CartContent.Status, &CartContent.Price, &CartContent.ItemName, &CartContent.Stock, &CartContent.Order, &CartContent.ItemID, &CartContent.Quantity)
		if CartContent.Quantity <= CartContent.Stock {
			CartContent.Status = "OutOfStock"
		}
		CartContents = append(CartContents, *CartContent)
	}
	return CartContents
}

type CartContentRequestPayload struct {
	ItemID   string `json:"ItemId"`
	Quantity int    `json:"Quantity"`
}

func (c CartContentRequestPayload) PostCartContent(CartID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,ItemID,Quantity
	ins, _ := db.Prepare(`
	INSERT 
	INTO 
		CartContent 
		(CartID , ItemID , Quantity) 
		VALUES 
		(? , ? , ?)`)
	defer ins.Close()
	// SQLの実行
	ins.Exec(CartID, c.ItemID, c.Quantity)
}
func searchCartContents(CartContents CartContents, ItemID string) bool {
	for _, CartContent := range CartContents {
		if CartContent.ItemID == ItemID {
			return true
		}
	}
	return false
}
func (c CartContentRequestPayload) UpdateCartContent(CartID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, _ := db.Prepare(`
	UPDATE 
		CartContent 
	SET 
		Quantity = ? 
	WHERE 
		CartId = ? 
	AND 
		ItemID = ?`)
	// SQLの実行
	ins.Exec(c.Quantity, CartID, c.ItemID)
	defer ins.Close()
}
func DeleteCartContent(CartID string, ItemID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, _ := db.Prepare(`
	DELETE 
	FROM
		CartContents

	WHERE 
		CartID = ? 

	AND 
		ItemID = ?`)
	// SQLの実行
	ins.Exec(CartID, ItemID)
	defer ins.Close()
}
func DeleteCartContentforTransaction(CartID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, _ := db.Prepare(`
	DELETE 
	FROM
		CartContents

	WHERE 
		CartID = ? `)
	// SQLの実行
	ins.Exec(CartID)
	defer ins.Close()
}
func DeleteItemFromCartContent(ItemID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, _ := db.Prepare(`
	DELETE 
	FROM 
		CartContents 

	WHERE 
		ItemID = ?`)
	defer ins.Close()

	// SQLの実行
	ins.Exec(ItemID)
}
func DeleteCart(CartID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	del, _ := db.Prepare(`
	DELETE 
	FROM 
		CartContents 
	WHERE 
		CartID = ?
	`)
	defer del.Close()

	// SQLの実行
	del.Exec(CartID)
}
