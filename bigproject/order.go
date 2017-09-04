package bigproject

import (
	"database/sql"

	"github.com/tokopedia/gosample/database"
)

type Order struct {
	OrderId       int64       `db:"order_id" json:"order_id"`
	InvoiceNumber string      `db:"invoice_ref_num" json:"invoice_ref_num"`
	OrderStatus   int         `db:"order_status" json:"order_status"`
	ShopId        int64       `db:"shop_id" json:"-"`
	ShopDetail    Shop        `json:"shop"`
	OrderDet      OrderDetail `json:"order_detail"`
}

type OrderAPI struct {
	Data struct {
		List []Order `json:"list"`
	} `json:"data"`
}

func GetOrder(orderid int64) (Order, error) {

	query := `SELECT 
				order_id, 
				invoice_ref_num, 
				order_status,
				shop_id
			FROM ws_order 
			WHERE order_id = $1
	`

	var result Order
	err := database.DBPool.MainDB.Get(&result, query, orderid)
	if err != nil && err != sql.ErrNoRows {
		return Order{}, err
	}

	return result, nil
}
