package model

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type ProxyGroup struct {
	Name          string   `yaml:"name,omitempty"`
	Type          string   `yaml:"type,omitempty"`
	Proxies       []string `yaml:"proxies,omitempty"`
	IsCountryGrop bool     `yaml:"-"`
	Url           string   `yaml:"url,omitempty"`
	Interval      int      `yaml:"interval,omitempty"`
	Tolerance     int      `yaml:"tolerance,omitempty"`
	Lazy          bool     `yaml:"lazy"`
	Size          int      `yaml:"-"`
}

type ProxyGroupsSortByName []ProxyGroup
type ProxyGroupsSortBySize []ProxyGroup

func (p ProxyGroupsSortByName) Len() int {
	return len(p)
}
func (p ProxyGroupsSortBySize) Len() int {
	return len(p)
}

func (p ProxyGroupsSortByName) Less(i, j int) bool {
	// 定义一组备选语言：首选英语，其次中文
	tags := []language.Tag{
		language.English,
		language.Chinese,
	}
	matcher := language.NewMatcher(tags)

	// 假设我们的请求语言是 "zh"（中文），则使用匹配器找到最佳匹配的语言
	bestMatch, _, _ := matcher.Match(language.Make("zh"))
	// 使用最佳匹配的语言进行排序
	c := collate.New(bestMatch)

	return c.CompareString(p[i].Name, p[j].Name) < 0
}

func (p ProxyGroupsSortBySize) Less(i, j int) bool {
	if p[i].Size == p[j].Size {
		return p[i].Name < p[j].Name
	}
	return p[i].Size < p[j].Size
}

func (p ProxyGroupsSortByName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p ProxyGroupsSortBySize) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
