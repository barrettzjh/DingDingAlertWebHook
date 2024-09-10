package loki

var (
	AlertChan = make(chan AlertManagerAlert)
)

func Server() {
	for {
		select {
		case alerts := <-AlertChan:
			for _, alert := range alerts.Alerts {
				// prometheus 的监控不参与traceid的过滤
				if alert.Labels.Metrics == "prometheus" {
					alert.Labels.Body = alert.Annotations.Summary

					Channel.c <- alert
					continue
				}

				if alert.Labels.Traceid == "" {
					alert.Labels.Traceid = "TraceIDNotFound"
				}
				if Channel.processed.Contains(alert.Labels.Body + alert.EndsAt.String()) {
					// alert.DropAlertMsg("队列中存在该body")
					continue
				}
				Channel.processed.Enqueue(alert.Labels.Body + alert.EndsAt.String())
				Channel.c <- alert
			}
		}
	}
}
