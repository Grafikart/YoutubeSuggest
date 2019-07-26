package main

import (
	"fmt"
	"github.com/Grafikart/YoutubeSuggest/api"
	"log"
	"os"
	"sort"
	"sync"
)

type RankingItem struct {
	api.Subscription
	count int
}

type Ranking []RankingItem

func (r Ranking) Len() int {
	return len(r)
}

func (r Ranking) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Ranking) Less(i, j int) bool {
	return r[i].count > r[j].count
}

func (r Ranking) Increment(s api.Subscription) Ranking {
	ok := false
	for k, i := range r {
		if i.ID == s.ID {
			r[k].count++
			ok = true
		}
	}
	if ok == false {
		r = append(r, RankingItem{s, 1})
	}
	return r
}

func main() {
	c := api.NewFileCache("./cache")
	app, err := api.NewAPI(os.Getenv("GOOGLE_API_KEY"), c)
	catch(err)
	subscriptions, err := app.Subscriptions(os.Args[0])
	catch(err)
	var wg sync.WaitGroup
	var mux sync.Mutex
	ranks := Ranking{}
	wg.Add(len(subscriptions))
	for _, s := range subscriptions {
		go func(sub api.Subscription) {
			defer wg.Done()
			subsubscriptions, err := app.Subscriptions(sub.ID)
			if err == nil {
				for _, s := range subsubscriptions {
					if isIn(s, subscriptions) == false {
						mux.Lock()
						ranks = ranks.Increment(s)
						mux.Unlock()
					}
				}
			} else {
				fmt.Println(err)
			}
		}(s)
	}
	wg.Wait()
	sort.Sort(ranks)
	for k, r := range ranks {
		if k < 20 {
			fmt.Printf("#%d %s <https://www.youtube.com/channel/%s> (%d liens)\n", k+1, r.Name, r.ID, r.count)
		}
	}
}

func isIn(sub api.Subscription, subs api.Subscriptions) bool {
	for _, s := range subs {
		if s.ID == sub.ID {
			return true
		}
	}
	return false
}

func catch(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
