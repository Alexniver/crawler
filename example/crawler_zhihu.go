package main

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strconv"
	"strings"
	"time"

	"github.com/Alexniver/crawler"
	"github.com/PuerkitoBio/goquery"
	"runtime"
)

func main() {
	logger := crawler.GetDefaultLogger()

	crawlFunc := func(spider *crawler.Spider, dataChannel chan *crawler.Spider) error {

		time.Sleep(spider.CrawlDelay)
		var client *http.Client

		// is use proxy
		if len(spider.ProxyUrl) > 0 {
			proxyUrl, err := url.Parse(spider.ProxyUrl)

			if nil != err {
				return err
			}

			client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: spider.TTL}
		} else {
			client = &http.Client{Timeout: spider.TTL}
		}

		logger.Info(spider.Url)
		logger.Info(spider.Generation)
		logger.Info(runtime.NumGoroutine())
		req, err := http.NewRequest(spider.Method, spider.Url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", spider.ContentType)
		//req.Header.Set("Cookie", "name=anny")

		resp, err := client.Do(req)

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		spider.RequestURL = resp.Request.URL
		spider.ResponseData = string(body)
		//logger.Info(spider.ResponseData)
		go func() {
			dataChannel <- spider
		}()

		return nil
	}

	analystFunc := func(spider *crawler.Spider, spiderChannel chan *crawler.Spider) error {

		if len(spider.ResponseData) <= 0 {
			return nil
		}
		reader := strings.NewReader(spider.ResponseData)

		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			return err
		}
		//fmt.Println(strconv.FormatInt(spider.Generation, 10) + "---" + spider.Url)

		doc.Url = spider.RequestURL

		// find link in response, and gen child spider, set child spider to spider channel
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			link, exist := s.Attr("href")
			if exist {
				if len(link) > 0 && strings.Index(link, "#") != 0 {
					if parsed, e := url.Parse(link); e == nil {
						parsed = doc.Url.ResolveReference(parsed)
						var parsedLink = parsed.String()
						//fmt.Println("parsing...")

						if strings.Index(parsedLink, "http") == 0 && strings.Index(parsedLink, spider.RequestURL.Host) > 0 {
							childSpider := spider.NewChildSpider(parsedLink)
							go func() {
								spiderChannel <- childSpider
							}()
						}
					}
				}
			}

		})

		//查找赞过万的页面
		/*if strings.Index(spider.Url, "question") > 0 && strings.Index(spider.Url, "answer") > 0 {
			doc.Find("span.count").Each(func(i int, s *goquery.Selection) {
				upCountStr, _ := s.Html()
				upCountStr = strings.Replace(upCountStr, "K", "000", -1)
				upCount, err := strconv.ParseInt(upCountStr, 10, 32)
				if err != nil {
					logger.Error(err)
				}
				fmt.Println(upCount)
				if upCount > 10000 {
					logger.Info(spider.Url)
				}
			})
		}*/

		//fmt.Println("analys end")
		//zhihu 相关, 如果当前页面是有问有答的, 则尝试拼接评论链接, 加到搜索中
		/*if strings.Index(spider.Url, "question") > 0 && strings.Index(spider.Url, "#answer") > 0 {
			//l := "http://www.zhihu.com/node/AnswerCommentBoxV2?params=%7B%22answer_id%22%3A%22"
			//l += spider.Url[strings.LastIndex(spider.Url, "answer\\"):]
			re, _ := regexp.Compile("[0-9]+")
			submatch := re.FindAllString(spider.Url, -1)
			//allindex := re.FindAllIndex([]byte(spider.Url), -1)
			l := "http://www.zhihu.com/node/AnswerCommentBoxV2?params=%7B%22answer_id%22%3A%22"
			l += submatch[1]
			l += "%22%2C%22load_all%22%3Atrue%7D"
			childSpider := spider.NewChildSpider(l)
			go func() {
				spiderChannel <- childSpider
			}()
			fmt.Println(spider.Url)
			fmt.Println(l)
		}

		//zhihu 相关, 如果访问的是评论的链接, 则在缓存中查找之前存储的答案和问题链接的map
		if strings.Index(spider.Url, "AnswerCommentBoxV2") > 0 {
			fmt.Println("comment")
		}*/

		//do other analyst
		/*doc.Find("zm-comment-content").Each(func(i int, s *goquery.Selection) {
			html, _ := s.Html()
			fmt.Println(html)
		})*/

		return nil

	}

	//crawler.DoCrawl("http://www.mi.com", crawlFunc, analystFunc, 1000)
	seedSpider := crawler.NewDefaultGetSpider("http://10.236.121.56:8080/admin/page!main.action")
	crawler.DoCrawl(seedSpider, crawlFunc, analystFunc, 1000)
}
