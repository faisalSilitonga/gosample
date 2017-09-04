package bigproject

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Product struct {
	ProductID          int64  `json:"product_id"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
	ProductStatus      string `json:"product_status"`
}

type ProductAPI struct {
	Data struct {
		Info Product `json:"info"`
	} `json:"data"`
}

func GetProductDetail(productId int64) (Product, error) {

	apiURL := fmt.Sprintf("http://devel-go.tkpd:3002/v4/product/get_detail.pl?product_id=%d", productId)
	var productAPI ProductAPI
	// create request
	req, err := http.NewRequest("GET", apiURL, nil)

	// send http get request to pulsa
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	dataJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Product{}, err
	}

	err = json.Unmarshal([]byte(dataJson), &productAPI)
	if err != nil {
		return Product{}, err
	}

	return productAPI.Data.Info, nil
}
