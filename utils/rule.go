package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"sub2clash/model"
)

func AddRulesByUrl(sub *model.Subscription, url string, proxy string) {
	get, err := Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(get.Body)
	bytes, err := io.ReadAll(get.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var payload model.Payload
	err = yaml.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range payload.Rules {
		payload.Rules[i] = payload.Rules[i] + "," + proxy
	}
	AddRules(sub, payload.Rules...)
}

func AddRuleProvider(
	sub *model.Subscription, providerName string, proxy string, provider model.RuleProvider,
) {
	if sub.RuleProviders == nil {
		sub.RuleProviders = make(map[string]model.RuleProvider)
	}
	sub.RuleProviders[providerName] = provider
	AddRules(
		sub,
		fmt.Sprintf("RULE-SET,%s,%s", providerName, proxy),
	)
}

func AddRules(sub *model.Subscription, rules ...string) {
	sub.Rules = append(rules, sub.Rules...)
}
