package parser

import (
	"strings"

	"github.com/h2so5/goback/regexp"

	"github.com/samber/lo"
)

// TODO: grouping book by similarity
func adjustName(book string) string {
	book = strings.TrimLeft(book, "《 ")
	book = strings.TrimRight(book, "》 ")
	if strings.HasSuffix(book, "）") {
		if re2, err := regexp.Compile("（第.版）$"); err == nil {
			book = re2.ReplaceAllString(book, "")
		}
	}
	if book == "富爸爸穷爸爸" || book == "穷爸爸富爸爸" || book == "穷爸爸，富爸爸" || book == "富爸爸" || strings.Contains(book, "富爸爸穷爸爸") {
		book = "富爸爸，穷爸爸"
	}
	if book == "聪明者的投资者" {
		book = "聪明的投资者"
	}
	if book == "解读基金" {
		book = "解读基金：我的投资观和实践"
	}
	if book == "有钱人想的和你不一样" {
		book = "有钱人和你想的不一样"
	}
	if book == "约翰聂夫的成功投资" || book == "约翰涅夫的成功投资" || book == "约翰 聂夫的成功投资" {
		book = "约翰·聂夫的成功投资"
	}
	if strings.HasPrefix(book, "小狗钱钱") {
		book = "小狗钱钱"
	}
	if strings.HasPrefix(book, "穷查理宝典") {
		book = "穷查理宝典"
	}
	if strings.HasPrefix(book, "漫步华尔街") {
		book = "漫步华尔街"
	}
	if strings.HasPrefix(book, "巴菲特的估值逻辑") {
		book = "巴菲特的估值逻辑"
	}
	if strings.HasPrefix(book, "投资最重要的事") {
		book = "投资最重要的事"
	}
	if strings.Contains(book, "巴菲特") && strings.Contains(book, "股东的信") {
		book = "巴菲特致股东的信"
	}
	if strings.Contains(book, "30年") && strings.Contains(book, "拿什么养活自己") {
		book = "30年后，你拿什么养活自己"
	}
	if strings.Contains(book, "手把手教你读财报") {
		book = "手把手教你读财报"
	}
	book = strings.Replace(book, "查理.芒格", "查理·芒格", 1)
	book = strings.Replace(book, "彼得.林奇", "彼得·林奇", 1)
	book = strings.Replace(book, "彼得林奇", "彼得·林奇", 1)

	book = strings.Replace(book, ":", "：", 1)

	return book
}

func ParseBooks(content string) map[string]int {

	if len(content) == 0 {
		return make(map[string]int)
	}

	reg, err := regexp.Compile(`(?<=《)(.*?)(?=》)`)
	if err != nil {
		return make(map[string]int)
	}

	res := make(map[string]int, 0)
	found := reg.FindAllString(content, -1)
	lo.ForEach(found, func(book string, _ int) {
		book = adjustName(book)
		res[book] = 1
	})

	return res
}
