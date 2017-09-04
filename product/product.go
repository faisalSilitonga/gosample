package product

import (
	"database/sql"

	"net/http"

	"strconv"

	"fmt"

	"encoding/json"

	"github.com/tokopedia/gosample/database"
)

type Product struct {
	ID    int64   `db:"product_id" json:"product_id"`
	Name  string  `db:"product_name" json:"product_name"`
	Price float64 `db:"normal_price" json:"price"`
}

func GetProduct(productID int64) (Product, error) {
	query := `SELECT 
				product_id,
				product_name,
				normal_price
			FROM
			ws_product
			WHERE 
			product_id = $1
	`

	var result Product
	err := database.DBPool.MainDB.Get(&result, query, productID)
	if err != nil && err != sql.ErrNoRows {
		return Product{}, err
	}

	return result, nil
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	productIDInt, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "ada error: %s", err.Error())
		return
	}

	productData, err := GetProduct(productIDInt)
	if err != nil {
		fmt.Fprintf(w, "ada error db : %s", err.Error())
		return
	}

	hasilJson, err := json.Marshal(productData)
	if err != nil {
		fmt.Fprintf(w, "ada error ketika parsing json: %s", err.Error())
	}

	hasilJson, err = json.Marshal(map[string]interface{}{
		"product_id":   productData.ID,
		"nama_product": productData.Name,
		"harga":        productData.Price,
	})

	if err != nil {
		fmt.Fprintf(w, "ada error ketika parsing json: %s", err.Error())
	}

	fmt.Fprint(w, string(hasilJson))
	return
}
