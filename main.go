package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
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
	router.GET("/", func(c *gin.Context) {
		DataStreamFromMQL5(c.Writer, c.Request)
	})
	router.GET("/ws", func(c *gin.Context) {
		StreamDatas(c.Writer, c.Request)
	})
	router.Run(":8802")
}

type AllDatas struct {
	AccountDatas []AccountDatas
}

type AccountDatas struct {
	XMLName xml.Name `xml:"data"`
	Text    string   `xml:",chardata"`
	Info    struct {
		Text      string `xml:",chardata"`
		Name      string `xml:"name,attr"`
		ID        string `xml:"id,attr"`
		Balance   string `xml:"balance,attr"`
		Equity    string `xml:"equity,attr"`
		PosProfit string `xml:"pos_profit,attr"`
		HisProfit string `xml:"his_profit,attr"`
		PerProfit string `xml:"per_profit,attr"`
	} `xml:"info"`
	Order []struct {
		Text      string `xml:",chardata"`
		Symbol    string `xml:"symbol,attr"`
		Type      string `xml:"type,attr"`
		Volume    string `xml:"volume,attr"`
		SlPrice   string `xml:"sl_price,attr"`
		TpPrice   string `xml:"tp_price,attr"`
		Profit    string `xml:"profit,attr"`
		CurrentPr string `xml:"current_pr,attr"`
	} `xml:"order"`
}

type DataIndi struct {
	XMLName xml.Name `xml:"data_indicator"`
	Text    string   `xml:",chardata"`
	Item    []struct {
		Text     string  `xml:",chardata"`
		IndiName string  `xml:"indi_name,attr"`
		Data     float64 `xml:"data,attr"`
	} `xml:"item"`
}

var account_datas []AllDatas

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
		fmt.Printf("msg: %s\n", msg)
		var data_resp []AllDatas
		xml.Unmarshal(msg, &data_resp)
		fmt.Printf("data_resp: %v\n", data_resp)
		// json, err := json.Unmarshal()

		// data_resp := AccountDatas{}
		// if err_mars := xml.Unmarshal(msg, &data_resp); err_mars == nil {
		// 	data := []DataStreamEAs{}
		// 	for _, v := range data_resp.Item {
		// 		cur := DataStreamEAs{
		// 			ID:        v.ID,
		// 			Name:      v.Name,
		// 			Balance:   v.Balance,
		// 			Equity:    v.Equity,
		// 			HisProfit: v.HisProfit,
		// 			PosProfit: v.PosProfit,
		// 			PerProfit: v.PerProfit,
		// 		}
		// 		data = append(data, cur)
		// 	}
		// 	ea_data = data
		// }
		account_datas = data_resp
		conn.WriteMessage(t, msg)
		tick = time.NewTicker(time.Second)
	}
}

type DataStreamIndis struct {
	Name string  `json:"name"`
	Data float64 `json:"data"`
}

type DataStreamEAs struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Equity    float64 `json:"equity"`
	PosProfit float64 `json:"pos_profit"`
	HisProfit float64 `json:"his_profit"`
	PerProfit float64 `json:"per_profit"`
}

type DataStreams struct {
	AccountDatas []AllDatas
}

var tick *time.Ticker

func StreamDatas(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	defer conn.Close()
	tick = time.NewTicker(time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-tick.C:
			{
				datastream := DataStreams{
					AccountDatas: account_datas,
				}
				conn.WriteJSON(datastream)
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
