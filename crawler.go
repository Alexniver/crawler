package crawler

import (
	"github.com/Alexniver/logger4go"
	"sync"
)

type SpiderFunc func(spider *Spider, spiderChannel chan *Spider) error

//url:tocrawl url, spiderCrawlFunc: how spiderCrawl, spiderAnalystFunc: how spider deal with response func
func DoCrawl(seedSpider *Spider, spiderCrawlFunc SpiderFunc, spiderAnlystFunc SpiderFunc, maxConcurrencyNum int) {
	logger := logger4go.GetDefaultLogger()

	visited := map[string]bool{}

	// crawler has two channel, one url channel, one data channel
	var spiderChannel = make(chan *Spider) // spider to crawl channel
	var dataChannel = make(chan *Spider)   // spider to analyst channel

	var wg sync.WaitGroup
	logger.Info("Crawl start!")
	//seedSpider := NewDefaultGetSpider(url)

	go func() {
		spiderChannel <- seedSpider
	}()

	// max crawl concurrency
	//start maxConcurrencyNum goroutine to crawl
	for i := 0; i < maxConcurrencyNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case spider := <-spiderChannel:
					if !visited[spider.Url] {
						//fmt.Println(len(visited))
						visited[spider.Url] = true
						//logger.Info(len(visited))
						err := spiderCrawlFunc(spider, dataChannel)
						if err != nil {
							logger.Error(err)
						}
						//fmt.Println("crawl end")
					}
				}
			}
		}()
	}

	//analys function will start a goroutine to analys when a response throw to dataChannel
	for {
		select {
		case data := <-dataChannel:
			wg.Add(1)
			go func() {
				defer wg.Done()
				//fmt.Println("analysising")
				err := spiderAnlystFunc(data, spiderChannel)
				if err != nil {
					logger.Error(err)
				}
				//fmt.Println("analys end")
			}()
		}
	}

	wg.Wait()

}
