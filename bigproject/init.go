package bigproject

import (
	"encoding/json"
	"net/http"
	"sync"

	nsq "github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"

	"strconv"

	"fmt"

	"log"

	"github.com/tokopedia/gosample/database"
)

type NsqRedis struct {
	Key   string
	Order OrderAPI
}

func CheckRedis(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "ERROR - FORBIDDEN")
			return
		}

		redisConn := database.RedisPool.RedisDev.Get()
		defer redisConn.Close()

		key := fmt.Sprint("faisal:", r.FormValue("order_id"))

		val, err := redis.String(redisConn.Do("GET", key))
		if err != nil {
			log.Printf("ERROR GET from REDIS")
			next.ServeHTTP(w, r)
		} else {
			log.Printf("GET from REDIS")
			fmt.Fprintf(w, "got %s", val)
		}
	})
}

func GetOrderDetailHandler(w http.ResponseWriter, r *http.Request) {
	var orderAPI = OrderAPI{}

	orderID := r.FormValue("order_id")
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "error orderid: %s", err.Error())
	}

	orderDataTemp, err := GetOrder(orderIDInt)
	if err != nil {
		fmt.Fprintf(w, "error get data order from db : %s", err.Error())
		return
	}

	var shopID = orderDataTemp.ShopId
	orderDataTemp.ShopDetail, err = GetShop(shopID)
	if err != nil {
		fmt.Fprintf(w, "error get data shop from api : %s", err.Error())
		return
	}

	var orderDetails []OrderDetail

	orderDetails, err = GetOrderDetail(orderIDInt)
	if err != nil {
		fmt.Fprintf(w, "error get order detail from db : %s", err.Error())
		return
	}

	for _, orderDetail := range orderDetails {
		var productId = orderDetail.ProductID
		orderDetail.ProductDetail, err = GetProductDetail(productId)
		if err != nil {
			fmt.Fprintf(w, "error get product from api: %s", err.Error())
			return
		}

		orderDataTemp.OrderDet = orderDetail
		orderAPI.Data.List = append(orderAPI.Data.List, orderDataTemp)
	}

	hasilJson, err := json.Marshal(orderAPI)
	if err != nil {
		fmt.Fprintf(w, "ada error ketika parsing json: %s", err.Error())
	}

	key := fmt.Sprint("faisal:", r.FormValue("order_id"))

	var result = string(hasilJson)

	var nsqData NsqRedis
	nsqData.Key = key
	nsqData.Order = orderAPI

	hasilJson, err = json.Marshal(nsqData)
	if err != nil {
		fmt.Fprintf(w, "ada error ketika parsing json: %s", err.Error())
	}

	err = CreateProducer(hasilJson)
	if err != nil {
		fmt.Fprintf(w, "ada error ketika nulis ke redis: %s", err.Error())
	}

	fmt.Fprint(w, result)
	return
}

func CreateProducer(value []byte) error {
	log.Printf("create Producer")
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("10.164.4.112:4150", config)

	err := w.Publish("faisal_order", value)
	if err != nil {
		return err
	}

	w.Stop()
	return nil
}

func CreateConsumer() error {
	log.Printf("create Consumer")
	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("faisal_order", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Handle Consumer")
		var nsqRedis NsqRedis
		err := json.Unmarshal(message.Body, &nsqRedis)
		log.Printf(string(message.Body))
		log.Printf(nsqRedis.Key)

		redisConn := database.RedisPool.RedisDev.Get()
		defer redisConn.Close()

		hasilJson, err := json.Marshal(nsqRedis.Order)
		if err != nil {
			return err
		}

		redisExpire := 60 * 24

		_, err = redisConn.Do("SETEX", nsqRedis.Key, redisExpire, string(hasilJson))
		if err != nil {
			return err
		}

		log.Printf("SETEX to REDIS")

		return nil
	}))
	err := q.ConnectToNSQD("10.164.4.112:4150")
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}
