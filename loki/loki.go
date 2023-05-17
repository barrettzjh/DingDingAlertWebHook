package loki

import (
	"time"
)

var (
	AlertChan = make(chan []LokiRuleAlertStruct)
)

func Server() {
	for {
		select {
		case alerts := <-AlertChan:
			for _, alert := range alerts {
				traceid := alert.Labels.Traceid
				if Channel.processed.Contains(traceid + alert.EndsAt.String()) {
					continue
				}
				Channel.processed.Enqueue(traceid + alert.EndsAt.String())
				Channel.c <- alert
			}
		}
	}
}

type LokiRuleAlertStruct struct {
	Annotations struct {
		Summary string `json:"summary"`
	} `json:"annotations"`
	EndsAt       time.Time `json:"endsAt"`
	StartsAt     time.Time `json:"startsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Labels       struct {
		Alertname string `json:"alertname"`
		Belong    string `json:"belong"`
		At		  string `json:"at"`
		Body      string `json:"body"`
		Job       string `json:"job"`
		Traceid   string `json:"traceid"`
		Stack     string `json:"attributes_exception_stacktrace"`
	} `json:"labels"`
}