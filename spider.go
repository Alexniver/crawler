package crawler

import (
	"errors"
	"net/url"
	"time"
)

var (
	// ErrEmptyHost is returned if a command to be enqueued has an URL with an empty host.
	ErrEmptyHost = errors.New("crawler fail : invalid empty host")

	// ErrDisallowed is returned when the requested URL is disallowed by the robots.txt
	// policy.
	ErrDisallowed = errors.New("crawler fail : disallowed by robots.txt")

	// ErrQueueClosed is returned when a Send call is made on a closed Queue.
	ErrQueueClosed = errors.New("crawler fail : send on a closed queue")

	//
	ErrTTL = errors.New("crawler fail : search time too long")
)

//spider implements crawler
type Spider struct {
	SeedUrl      string
	ParentUrl    string //this spider's parent spider's url
	Url          string // url to search
	Method       string
	ProxyUrl     string // proxyUrl to use
	CrawlDelay   time.Duration
	TTL          time.Duration //spider time-to-live
	ContentType  string
	RobotsTxtUrl string
	Generation   int64 // which generation this spider is
	RequestURL   *url.URL
	ResponseData string
}

func NewSpider(_seedUrl string, _parentUrl string, _url string, _method string, _proxyUrl string, _crawlDelay time.Duration, _ttl time.Duration, _contentType string, _robotsTxtUrl string, _generation int64) *Spider {
	return &Spider{Url: _url, Method: _method, ProxyUrl: _proxyUrl, CrawlDelay: _crawlDelay, TTL: _ttl, ContentType: _contentType, RobotsTxtUrl: _robotsTxtUrl, Generation: _generation}
}

// new default spider , method get crawlDelay 2 second, ttl 10 second
func NewDefaultGetSpider(_url string) *Spider {
	return NewSpider(_url, "", _url, "GET", "", 1*time.Second, 3*time.Second, "application/x-www-form-urlencoded", "/robots.txt", 0)
}

//return a child spider by this spider
func (spider *Spider) NewChildSpider(_url string) *Spider {
	return NewSpider(spider.SeedUrl, spider.Url, _url, spider.Method, spider.ProxyUrl, spider.CrawlDelay, spider.TTL, spider.ContentType, spider.RobotsTxtUrl, spider.Generation+1)
}
