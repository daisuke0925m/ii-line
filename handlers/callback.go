package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

type tickers struct {
	Daily []struct {
		ID        int       `json:"id"`
		Symbol    string    `json:"symbol"`
		Date      time.Time `json:"date"`
		Open      float64   `json:"open"`
		High      float64   `json:"high"`
		Low       float64   `json:"low"`
		Close     float64   `json:"close"`
		Volume    int       `json:"volume"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"daily"`
}

func fetchAPI(message string) (tickers, error) {

	url := "https://api.index-indicators.com/ticker?symbol=" + message

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return tickers{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return tickers{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tickers{}, err
	}

	var tickers tickers
	if err := json.Unmarshal(body, &tickers); err != nil {
		log.Fatal(err)
	}

	return tickers, nil
}

func parseReplyMsg(symbol string, date time.Time, open float64, high float64, low float64, close float64, volume int) (repMsg string, err error) {
	parsedDate := strings.Split(date.String(), " ")[0]
	repMsg = symbol + "\n" +
		parsedDate + "\n" +
		"open " + strconv.FormatFloat(open, 'f', 0, 64) + "\n" +
		"high " + strconv.FormatFloat(high, 'f', 0, 64) + "\n" +
		"low " + strconv.FormatFloat(low, 'f', 0, 64) + "\n" +
		"close " + strconv.FormatFloat(close, 'f', 0, 64) + "\n" +
		"volume " + strconv.Itoa(volume)
	return repMsg, nil
}

func LineHandler(w http.ResponseWriter, r *http.Request) {
	// BOTを初期化
	bot, err := linebot.New(
		os.Getenv("LINE_SECRET"),
		os.Getenv("LINE_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// リクエストからBOTのイベントを取得
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				tickers, err := fetchAPI(message.Text)
				if err != nil {
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
					if err != nil {
						log.Print(err)
					}
				}

				latestData := tickers.Daily[0]
				replyMessage, err := parseReplyMsg(latestData.Symbol, latestData.Date, latestData.Open, latestData.High, latestData.Low, latestData.Close, latestData.Volume)
				if err != nil {
					log.Print(err)
				}

				if err != nil {
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
					if err != nil {
						log.Print(err)
					}
				}
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}
