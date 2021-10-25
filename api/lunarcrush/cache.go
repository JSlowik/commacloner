package lunarcrush

import (
	"github.com/jslowik/commacloner/api/lunarcrush/rest"
	"github.com/jslowik/commacloner/api/lunarcrush/rest/dobjs"
	"github.com/jslowik/commacloner/config"
	"go.uber.org/zap"
	"sort"
	"sync"
	"time"
)

type LunarCache struct {
	Logger *zap.SugaredLogger
	Config config.LunarCrushAPI

	mu            sync.Mutex
	pairs [] dobjs.PairData
	Blacklist map[string]bool
}

func (cache *LunarCache) UpdateCache(){
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for {
		p, err  := rest.GetPairs(cache.Config.Key,"")
		if err == nil {
			cache.pairs = p
			break
		} else {
			cache.Logger.Warnf("could not update cache: %v - Will wait 5 seconds then try again", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (cache *LunarCache)GetByGalaxyScore(validPairs map[string][]string, max int) []string {
	return cache.getMatchingPairs(sortByGalaxyScore,validPairs,max)
}

func (cache *LunarCache)GetByAltRank(validPairs map[string][]string, max int) []string {
	return cache.getMatchingPairs(sortByAltRank,validPairs,max)
}

type sorter func([] dobjs.PairData) [] dobjs.PairData

func (cache *LunarCache)getMatchingPairs (sortFunc sorter, validPairs map[string][]string, max int) []string {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	sorted := sortFunc(cache.pairs)

	pairs := make([]string,0)
	for _, pair := range sorted {
		if validPairs[pair.S] != nil {
			if !cache.Blacklist[pair.S] {
				pairs = append(pairs,pair.S)
			} else {
				cache.Logger.Infof("Pair %s in CommaCloner blacklist, ignoring", pair.S)
			}
		}
		if len(pairs) == max {
			break
		}
	}
	return pairs
}



func sortByAltRank(pairs [] dobjs.PairData) [] dobjs.PairData {
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Acr < pairs[j].Acr
	})
	return pairs
}



func sortByGalaxyScore(pairs [] dobjs.PairData) [] dobjs.PairData {
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Gs > pairs[j].Gs
	})
	return pairs
}