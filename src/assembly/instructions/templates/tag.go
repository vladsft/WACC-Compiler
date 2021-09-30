package main

import (
	"fmt"
	"regexp"
	"strings"
)

const tagString = "(?U)<%s>%s</%s>"

var (
	oneLine = newTag("1line", "(.*\n)*.*", func(str string) string {
		noNewLines := strings.ReplaceAll(str, "\n", "")
		return strings.Join(strings.Fields(noNewLines), " ")
	})

	operand = newTag("op", ".*", func(str string) string {
		return "<" + str + "." + strings.Title(targetLang) + "Operand>"
	})
)

type tag struct {
	token   string
	matcher *regexp.Regexp
	expand  func(string) string
}

func newTag(token, match string, expand func(string) string) tag {
	return tag{
		token:   token,
		matcher: regexp.MustCompile(fmt.Sprintf(tagString, token, match, token)),
		expand:  expand,
	}
}

func (t tag) apply(str string) string {
	return t.matcher.ReplaceAllStringFunc(str, func(str string) string {
		noStartTag := strings.Replace(str, "<"+t.token+">", "", 1)
		noEndTag := strings.Replace(noStartTag, "</"+t.token+">", "", 1)
		return t.expand(noEndTag)
	})
}
