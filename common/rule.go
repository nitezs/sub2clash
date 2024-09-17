package common

import (
	"fmt"
	"strings"

	"github.com/nitezs/sub2clash/model"
)

func PrependRuleProvider(
	sub *model.Subscription, providerName string, group string, provider model.RuleProvider,
) {
	if sub.RuleProviders == nil {
		sub.RuleProviders = make(map[string]model.RuleProvider)
	}
	sub.RuleProviders[providerName] = provider
	PrependRules(
		sub,
		fmt.Sprintf("RULE-SET,%s,%s", providerName, group),
	)
}

func AppenddRuleProvider(
	sub *model.Subscription, providerName string, group string, provider model.RuleProvider,
) {
	if sub.RuleProviders == nil {
		sub.RuleProviders = make(map[string]model.RuleProvider)
	}
	sub.RuleProviders[providerName] = provider
	AppendRules(sub, fmt.Sprintf("RULE-SET,%s,%s", providerName, group))
}

func PrependRules(sub *model.Subscription, rules ...string) {
	if sub.Rules == nil {
		sub.Rules = make([]string, 0)
	}
	sub.Rules = append(rules, sub.Rules...)
}

func AppendRules(sub *model.Subscription, rules ...string) {
	if sub.Rules == nil {
		sub.Rules = make([]string, 0)
	}
	matchRule := sub.Rules[len(sub.Rules)-1]
	if strings.Contains(matchRule, "MATCH") {
		sub.Rules = append(sub.Rules[:len(sub.Rules)-1], rules...)
		sub.Rules = append(sub.Rules, matchRule)
		return
	}
	sub.Rules = append(sub.Rules, rules...)
}
