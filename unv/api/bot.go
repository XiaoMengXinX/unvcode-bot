package api

import (
	"encoding/json"
	"fmt"
	"github.com/XiaoMengXinX/go-unvcode"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unv"
)

type response struct {
	QueryID   string              `json:"inline_query_id"`
	Method    string              `json:"method"`
	Results   []inlineQueryResult `json:"results"`
	CacheTime int64               `json:"cache_time"`
}

type inlineQueryResult struct {
	Type                string              `json:"type"`
	Id                  int64               `json:"id"`
	Title               string              `json:"title"`
	InputMessageContent inputMessageContent `json:"input_message_content"`
	Description         string              `json:"description"`
}

type inputMessageContent struct {
	MessageText string `json:"message_text"`
}

var unv *unvcode.Unv

func init() {
	unv, _ = unvcode.New(font.Font)
}

func UnvBot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update

	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.InlineQuery != nil {
		if update.InlineQuery.Query == "" {
			return
		}

		text, mse := unv.Unvcode(update.InlineQuery.Query)
		timeNow := time.Now().UnixNano()

		data := response{
			Method:    "answerInlineQuery",
			QueryID:   update.InlineQuery.ID,
			CacheTime: 3600,
			Results: []inlineQueryResult{
				{
					Type:                "article",
					Id:                  timeNow,
					Title:               "反和谐结果",
					InputMessageContent: inputMessageContent{MessageText: text},
					Description:         text,
				},
				{
					Type:                "article",
					Id:                  timeNow + 1,
					Title:               "字符相似度",
					InputMessageContent: inputMessageContent{MessageText: fmt.Sprintln(mse)},
					Description:         fmt.Sprintln(mse),
				},
			},
		}
		msg, _ := json.Marshal(data)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, string(msg))
	}
}
