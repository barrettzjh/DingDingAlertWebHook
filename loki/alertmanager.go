package loki

import (
	"fmt"
	"time"
)

type AlertManagerAlert struct {
	Receiver string `json:"receiver"`
	Status string `json:"status"`
	Alerts []Alert `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		AttributesExceptionStacktrace string `json:"attributes_exception_stacktrace"`
		Belong string `json:"belong"`
		Body string `json:"body"`
		Job string `json:"job"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
	Version string `json:"version"`
	GroupKey string `json:"groupKey"`
	TruncatedAlerts int `json:"truncatedAlerts"`
}

type Alert struct {
	Status string `json:"status"`
	Labels struct {
		Alertname string `json:"alertname"`
		AttributesExceptionStacktrace string `json:"attributes_exception_stacktrace"`
		Belong string `json:"belong"`
		Body string `json:"body"`
		Job string `json:"job"`
		At string `json:"at"`
		Traceid string `json:"traceid"`
		Metrics string `json:"metrics"`
	} `json:"labels"`
	Annotations struct {
		Summary string `json:"summary"`
	} `json:"annotations"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt time.Time `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
	Fingerprint string `json:"fingerprint"`
}

func (alert Alert)GetAlertMsg()string{
	if alert.Labels.Metrics == "prometheus"{
		return fmt.Sprintf("监控告警:\n应用: %s\n告警内容: %s\n", alert.Labels.Job, alert.Annotations.Summary)
	}else {
		return fmt.Sprintf("日志告警:\n应用: %s\ntraceid: %s\n日志内容: %s\n日志描述: %s\n堆栈信息: %s\n", alert.Labels.Job, alert.Labels.Traceid, alert.Labels.Body, alert.Annotations.Summary, alert.Labels.AttributesExceptionStacktrace)
	}
}

func (alert Alert)DropAlertMsg(s string){
	fmt.Printf("丢弃告警：%s\n原因：%s\ntraceid：%s\nbody:%s\n", alert.Labels.Alertname,s ,alert.Labels.Traceid, alert.Labels.Body)
}