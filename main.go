package main

import (
	"encoding/json"
	"faceScore/faceScore"
	"flag"
	"fmt"
	"log"
)

var (
	appCode  = ""
	local    = true
	location = ""
)

func init() {
	flag.StringVar(&appCode, "appCode", "", "用户的appCode")
	flag.BoolVar(&local, "local", true, "是否本地图片")
	flag.StringVar(&location, "location", "", "地址")

}

func main() {
	flag.Parse()

	scorer := faceScore.NewScorer(appCode)

	var result *faceScore.Result
	var err error

	if local {
		result, err = scorer.LocalScore(location)
	} else {
		result, err = scorer.WebScore(location)
	}

	view := viewResult{}
	if err == nil {
		view.Ok = result.HasError() == nil
		view.Msg = result.Message
		view.Score = fmt.Sprintf("%.1f", result.GetSocre())

	} else {
		view.Ok = false
		view.Msg = err.Error()
	}
	jsonMsg, _ := json.Marshal(view)
	log.Println(string(jsonMsg))
	fmt.Println(string(jsonMsg))
}

type viewResult struct {
	Ok    bool
	Msg   string
	Score string
}
