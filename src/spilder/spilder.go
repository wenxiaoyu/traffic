package spilder

import (
	"github.com/axgle/mahonia"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type KD int // kind of item li

const (
	SH = KD(iota)
	JK
	RQ
	YS
	BJ
	MY
	MS
	CS
)

type ItemLi struct {
	Title     string
	Date      time.Time
	ImgUrl    string
	ShortDesc string
	DetailUrl string
	Kind      KD
	Dt        *D
}

type D struct {
	Intro  string
	ImgUrl string
	P      []string
}

func Detail(url string) *D {
	detail := new(D)

	resp, err := http.Get(url)
	if err == nil {
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			return nil
		}
		html := string(body)
		enc := mahonia.NewDecoder("UTF-8")
		encHtml := enc.ConvertString(html)
		doc, _ := gokogiri.ParseHtml([]byte(encHtml))

		img, _ := doc.Search("//section[@class='show_box']/div[2]/img")
		intro, _ := doc.Search("//div[@class='daoyu']")
		content, _ := doc.Search("//section[@class='show_box']/article/p")

		if len(img) > 0 {
			detail.ImgUrl = img[0].Attr("src")
		}
		if len(intro) > 0 {
			detail.Intro = intro[0].Content()
		}

		for _, v := range content {

			if strings.Contains(v.Content(), "浏览大图") {
				continue
			}
			if strings.Contains(v.InnerHtml(), "text-align: center;") {
				continue
			}
			//log.Print(detail.P)
			detail.P = append(detail.P, v.InnerHtml())
		}
		doc.Free()
	} else {
		log.Println(err)
	}
	return detail
}

func Li(d string, url string) *ItemLi {
	li := new(ItemLi)
	if strings.Contains(url, "sh") {
		li.Kind = SH
	} else if strings.Contains(url, "jk") {
		li.Kind = JK
	} else if strings.Contains(url, "rq") {
		li.Kind = RQ
	} else if strings.Contains(url, "ys") {
		li.Kind = YS
	} else if strings.Contains(url, "bj") {
		li.Kind = BJ
	} else if strings.Contains(url, "my") {
		li.Kind = MY
	} else if strings.Contains(url, "ms") {
		li.Kind = MS
	} else if strings.Contains(url, "cs") {
		li.Kind = CS
	}
	resp, err := http.Get(url)
	if err == nil {
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			return nil
		}

		html := string(body)
		enc := mahonia.NewDecoder("UTF-8")
		encHtml := enc.ConvertString(html)
		doc, _ := gokogiri.ParseHtml([]byte(encHtml))

		title, _ := doc.Search("//ul/li[1]/a/h3")
		t, _ := doc.Search("//ul/li[1]/a/div[@class='time']")
		a, _ := doc.Search("//ul/li[1]/a")
		icon, _ := doc.Search("//ul/li[1]/a/img")
		intro, _ := doc.Search("//ul/li[1]/a/p")

		if strings.Compare(d, "") == 0 || strings.Contains(t[0].Content(), d) {

			if len(t) > 0 {
				li.Date = time.Now()
			}
			if len(a) > 0 {
				li.DetailUrl = a[0].Attr("href")
			}
			if len(icon) > 0 {
				li.ImgUrl = icon[0].Attr("src")
			}
			if len(intro) > 0 {
				li.ShortDesc = intro[0].Content()
			}
			if len(title) > 0 {
				li.Title = title[0].Content()
			}
		}

		doc.Free()
	} else {
		log.Println(err)
	}
	return li
}
