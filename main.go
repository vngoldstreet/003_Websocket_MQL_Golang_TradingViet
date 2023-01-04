package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	// origin := r.Header.Get("Origin")
	// if strings.Contains(origin, "vang247") || strings.Contains(origin, "localhost") {
	// 	return true
	// } else {
	// 	return false
	// }
	return true
}}

func main() {
	router := gin.Default()
	router.POST("/api/orders", apiOrderRequest)
	router.GET("/", func(c *gin.Context) {
		DataStreamFromMQL5(c.Writer, c.Request)
	})
	router.GET("/ws", func(c *gin.Context) {
		StreamDatas(c.Writer, c.Request)
	})
	router.Run(":8080")
}

type Data struct {
	XMLName   xml.Name `xml:"data"`
	Text      string   `xml:",chardata"`
	Name      string   `xml:"name"`
	ID        int64    `xml:"id"`
	Balance   float64  `xml:"balance"`
	Equity    float64  `xml:"equity"`
	PosProfit float64  `xml:"pos_profit"`
	HisProfit float64  `xml:"his_profit"`
	PerProfit float64  `xml:"per_profit"`
}

func DataStreamFromMQL5(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	defer conn.Close()
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		data := Data{}
		xml.Unmarshal(msg, &data)
		message = []DataStreams{
			{
				// Name:      data.Name,
				ID:        data.ID,
				Balance:   data.Balance,
				Equity:    data.Equity,
				PosProfit: data.PosProfit,
				HisProfit: data.HisProfit,
				PerProfit: data.PerProfit,
			},
		}
		tick = time.NewTicker(time.Millisecond * 100)
		conn.WriteMessage(t, msg)
	}
}

type DataStreams struct {
	// Name      string  `json:"name"`
	ID        int64   `json:"id"`
	Balance   float64 `json:"balance"`
	Equity    float64 `json:"equity"`
	PosProfit float64 `json:"pos_profit"`
	HisProfit float64 `json:"his_profit"`
	PerProfit float64 `json:"per_profit"`
}

var message []DataStreams

var tick *time.Ticker

func StreamDatas(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	defer conn.Close()
	quit := make(chan struct{})
	for {
		select {
		case <-tick.C:
			{
				conn.WriteJSON(message)
			}
		case <-quit:
			tick.Stop()
			return
		}
	}
}

type Orders struct {
	ID         string  `json:"id"`
	Side       string  `json:"side"`
	Price      float64 `json:"price"`
	Stoploss   float64 `json:"stoploss"`
	TakeProfit float64 `json:"takeprofit"`
	Volume     float64 `json:"volume"`
	CPrice     float64 `json:"cprice"`
	CProfit    float64 `json:"cprofit"`
}

var res Orders

func apiOrderRequest(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "7884d9ca4d43403acc22f51b597ad15fbc0cac34d2625fef9b3e8be8a25faffa" {
		id := c.PostForm("id")
		side := c.PostForm("side")
		price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
		stoploss, _ := strconv.ParseFloat(c.PostForm("stoploss"), 64)
		takeprofit, _ := strconv.ParseFloat(c.PostForm("takeprofit"), 64)
		volume, _ := strconv.ParseFloat(c.PostForm("volume"), 64)
		cprice, _ := strconv.ParseFloat(c.PostForm("cprice"), 64)
		cprofit, _ := strconv.ParseFloat(c.PostForm("cprofit"), 64)
		res = Orders{
			ID:         id,
			Side:       side,
			Price:      price,
			Stoploss:   stoploss,
			TakeProfit: takeprofit,
			Volume:     volume,
			CPrice:     cprice,
			CProfit:    cprofit,
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Wrong token!",
		})
	}
}
