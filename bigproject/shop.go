package bigproject

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Shop struct {
	ShopId     int64  `json:"shop_id"`
	ShopName   string `json:"shop_name"`
	ShopDomain string `json:"shop_domain"`
	ShopStatus int    `json:"shop_status"`
}

type ShopAPI struct {
	Id     int64  `json:"shop_id"`
	Name   string `json:"shop_name"`
	Domain string `json:"domain"`
	Status int    `json:"status"`
}

type ShopAPIArr struct {
	Data []ShopAPI `json:"data"`
}

func GetShop(shopID int64) (Shop, error) {

	apiURL := fmt.Sprintf("http://devel-go.tkpd:3002/v1/shop/get_summary?shop_id=%d", shopID)
	var arrShop ShopAPIArr
	// create request
	req, err := http.NewRequest("GET", apiURL, nil)

	// send http get request to pulsa
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		return Shop{}, err
	}
	defer resp.Body.Close()

	dataJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Shop{}, err
	}

	err = json.Unmarshal([]byte(dataJson), &arrShop)
	if err != nil {
		return Shop{}, err
	}

	if len(arrShop.Data) < 0 {
		return Shop{}, errors.New("No Shop Data")
	}

	var shop = Shop{}
	shop.ShopId = arrShop.Data[0].Id
	shop.ShopName = arrShop.Data[0].Name
	shop.ShopDomain = arrShop.Data[0].Domain
	shop.ShopStatus = arrShop.Data[0].Status

	return shop, nil
}
