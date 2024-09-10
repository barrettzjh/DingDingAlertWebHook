package loki

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/barrettzjh/DingDingAlertWebHook/utils"
	"golang.org/x/time/rate"
)

const (
	logMaxLength            = 400
	notificationChannelSize = 10000
	defaultNotifyTimeOut    = time.Minute * 10
)

var (
	webHooks     = map[string]string{}
	Channel      = NewNotificationChannel(notificationChannelSize) // 初始化告警通道
	notifiers    = make(map[string]Notifier)                       // 存储已创建的notifier
	notifierLock sync.Mutex                                        // 保证并发安全
)

// 运行程序先自动初始化钉钉通知通道
func init() {
	// 通过环境变量WEB_HOOK_来确定钉钉机器人webhook和通道的关系
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		key := pair[0]
		if strings.HasPrefix(key, "WEB_HOOK_") {
			key := strings.TrimPrefix(key, "WEB_HOOK_")
			if _, ok := webHooks[key]; !ok {
				webHooks[key] = pair[1]
			}
		}
	}

	// 根据获取到的webhook对应关系，进行创建通知通道
	for name := range webHooks {
		notifiers[name] = createNotifier(name)
	}

	go Channel.StartNotifier()
}

// 告警媒介接口，这里只实现了钉钉的
type Notifier interface {
	Notify(ctx context.Context, msg interface{}) error
	Close()
}

type DingDingNotifier struct {
	WebHook string
	stopCh  chan struct{} // 用于停止通知goroutine
	limiter *rate.Limiter // 添加Limiter成员变量
}

func (cn *DingDingNotifier) Notify(ctx context.Context, msg interface{}) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-cn.stopCh:
			return nil
		default:
			if err := cn.limiter.Wait(context.Background()); err != nil {
				fmt.Printf("Error waiting for limiter: %v\n", err)
			}
			err := utils.SendDingTalkResolvedWithWebHook(cn.WebHook, "告警", msg)
			if err != nil {
				fmt.Printf("Error sending notification: %v\n", err)
			} else {
				return nil
			}
		}
	}
}

func (cn *DingDingNotifier) Close() {
	close(cn.stopCh)
}

type NotifierFactory struct{}

func (nf NotifierFactory) GetNotifier(belong string) Notifier {
	if notifier, ok := notifiers[belong]; ok {
		return notifier
	}
	return nil
}

func createNotifier(key string) Notifier {
	notifierLock.Lock()
	defer notifierLock.Unlock()

	n := &DingDingNotifier{
		WebHook: webHooks[key],
		stopCh:  make(chan struct{}),
		// 每秒钟生成0.3个令牌，最多能存一个令牌，约等于 4秒钟一个令牌，钉钉api限流为1分钟20次
		limiter: rate.NewLimiter(0.3, 1),
	}
	notifiers[key] = n
	return n

}

type NotificationChannel struct {
	c         chan Alert
	wg        sync.WaitGroup
	processed *FixedSizeQueue
	closed    bool
	mutex     sync.Mutex
}

func NewNotificationChannel(fixedSize int) *NotificationChannel {
	return &NotificationChannel{
		c:         make(chan Alert),
		closed:    false,
		processed: NewFixedSizeQueue(fixedSize),
		mutex:     sync.Mutex{},
	}
}

func (nc *NotificationChannel) StartNotifier() {
	cache := NewCache(time.Hour)

	for {
		select {
		case notification, ok := <-nc.c:
			if !ok {
				return
			}

			// 避免接受大量无效重复告警
			if ok, _ := cache.getValue(notification); ok {
				// notification.DropAlertMsg("body内容相同，一小时仅告警一次")
				continue
			}
			cache.setValue(notification)

			// 如果限速后的告警仍然很多，协程太多堆积，这里如果10分钟仍未发出通知，则放弃该条通知
			ctx, cancel := context.WithTimeout(context.Background(), defaultNotifyTimeOut)
			defer cancel()

			go func() {
				nf := NotifierFactory{}
				notifier := nf.GetNotifier(notification.Labels.Belong)
				// 如果在Loki的告警规则中没有配置belong这个label，或webHooks里没有定义的，均不进行告警
				if notifier == nil {
					return
				}
				// 防止钉钉接口异常报错
				if len(notification.Labels.AttributesExceptionStacktrace) > logMaxLength {
					notification.Labels.AttributesExceptionStacktrace = notification.Labels.AttributesExceptionStacktrace[:logMaxLength] + "..."
				}
				if len(notification.Labels.Body) > logMaxLength {
					notification.Labels.Body = notification.Labels.Body[:logMaxLength] + "..."
				}

				err := notifier.Notify(ctx, map[string]interface{}{
					"msgtype": "text",
					"text": map[string]string{
						"content": notification.GetAlertMsg(),
					},
					"at": map[string]interface{}{
						"atMobiles": []string{notification.Labels.At},
						"isAtAll":   false,
					},
				},
				)
				if err != nil {
					fmt.Printf("Error sending notification: %v\n", err)
				}
			}()
		}
	}
}

func (nc *NotificationChannel) Close() {
	if !nc.closed {
		close(nc.c)
		nc.closed = true
		nc.wg.Wait()
	}
}
