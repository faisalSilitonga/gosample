package bigproject

import (
	"database/sql"

	"github.com/tokopedia/gosample/database"
)

type OrderDetail struct {
	Quantity      int64   `db:"quantity_deliver" json:"quantity_delivered"`
	ProductID     int64   `db:"product_id" json:"-"`
	ProductDetail Product `json:"product"`
}

func GetOrderDetail(orderid int64) ([]OrderDetail, error) {

	query := `SELECT 
				product_id, 
				quantity_deliver
			FROM ws_order_dtl
			WHERE order_id = $1
	`

	var result []OrderDetail
	err := database.DBPool.MainDB.Select(&result, query, orderid)
	if err != nil && err != sql.ErrNoRows {
		return []OrderDetail{}, err
	}

	return result, nil
}
