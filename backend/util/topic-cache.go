package util

import (
	"sync"
)

type TopicCache struct {
	mu   sync.RWMutex
	data map[string]map[string][]int64
}

var StarGazerTopicCache = NewTopicCache()

type TopicInfo struct {
	Name  string  `json:"name"`
	Repos []int64 `json:"repos"`
}

func NewTopicCache() *TopicCache {
	return &TopicCache{
		data: make(map[string]map[string][]int64),
	}
}

func (c *TopicCache) SetTopics(email string, data map[string][]int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.data[email]; !exists {
		c.data[email] = make(map[string][]int64)
	}
	c.data[email] = data
}

func (c *TopicCache) GetTopics(email string) ([]*TopicInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	topicsMap, exists := c.data[email]
	if !exists {
		return []*TopicInfo{}, false
	}

	topics := make([]*TopicInfo, 0, len(topicsMap))
	for topicName, repos := range topicsMap {
		topics = append(topics, &TopicInfo{
			Name:  topicName,
			Repos: repos,
		})
	}
	return topics, true
}
