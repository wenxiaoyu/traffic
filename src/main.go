package main

import (
	"./spilder"
	"bufio"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var homeDir = "/home/fish"

func main() {

	fileServer := http.FileServer(http.Dir(homeDir + "/changshi"))
	http.HandleFunc("/productLiHtml", productLiHtml)
	http.HandleFunc("/count", count)
	go func() {
		http.ListenAndServe(":8888", nil)
	}()
	go func() {
		err := http.ListenAndServe(":8080", fileServer)
		if err != nil {
			log.Fatal("ListenAndServer ", err)
		}
	}()

	go timer()

	log.Println("Listen http port 8080!")

	select {}
}

//间隔1小时重新拉取一次数据
func timer() {
	timer := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-timer.C:
			go productLiHtml(nil, nil)
		}
	}
}

var l sync.Mutex

func count(w http.ResponseWriter, req *http.Request) {
	l.Lock()
	defer l.Unlock()

	fileName := homeDir + "/logs/changshi.log"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(req.RemoteAddr + " " + req.FormValue("page") + " " + req.UserAgent())
}
func productLiHtml(w http.ResponseWriter, req *http.Request) {

	//log.Println("----->timer test.")
	d := ""

	domain := "http://m.360changshi.com/"
	sh := spilder.Li(d, domain+"sh/")
	jk := spilder.Li(d, domain+"jk/")
	rq := spilder.Li(d, domain+"rq/")
	ys := spilder.Li(d, domain+"ys/")
	bj := spilder.Li(d, domain+"bj/")
	my := spilder.Li(d, domain+"my/")
	ms := spilder.Li(d, domain+"ms/")
	cs := spilder.Li(d, domain+"cs/")

	shd := spilder.Detail(sh.DetailUrl)
	jkd := spilder.Detail(jk.DetailUrl)
	rqd := spilder.Detail(rq.DetailUrl)
	ysd := spilder.Detail(ys.DetailUrl)
	bjd := spilder.Detail(bj.DetailUrl)
	myd := spilder.Detail(my.DetailUrl)
	msd := spilder.Detail(ms.DetailUrl)
	csd := spilder.Detail(cs.DetailUrl)

	sh.Dt = shd
	jk.Dt = jkd
	rq.Dt = rqd
	ys.Dt = ysd
	bj.Dt = bjd
	my.Dt = myd
	ms.Dt = msd
	cs.Dt = csd

	pmap := make(map[string]*spilder.ItemLi)
	pmap["sh"] = sh
	pmap["jk"] = jk
	pmap["rq"] = rq
	pmap["ys"] = ys
	pmap["bj"] = bj
	pmap["my"] = my
	pmap["ms"] = ms
	pmap["cs"] = cs

	t := template.New("Li template")
	t.Funcs(template.FuncMap{"dateconv": Dateconv, "randBg": RandBg, "unescaped": Unescaped})

	t.ParseFiles(homeDir+"/changshi/li.tmpl", "/home/fish/changshi/detail.tmpl")

	fileName := homeDir + "/changshi/today/today.html"
	if _, err := os.Stat(fileName); os.IsExist(err) {
		os.Remove(fileName)
	}
	f, _ := os.Create(fileName)
	bw := bufio.NewWriter(f)

	var hisName = time.Now().Format("2006-01-02")
	bakfile := homeDir + "/changshi/his/li_" + hisName + ".html"
	if _, err := os.Stat(bakfile); os.IsExist(err) {
		os.Remove(bakfile)
	}
	bakf, _ := os.Create(bakfile)
	bakbw := bufio.NewWriter(bakf)
	sh.DetailUrl = "/today/detail_sh.html"
	jk.DetailUrl = "/today/detail_jk.html"
	rq.DetailUrl = "/today/detail_rq.html"
	ys.DetailUrl = "/today/detail_ys.html"
	bj.DetailUrl = "/today/detail_bj.html"
	my.DetailUrl = "/today/detail_my.html"
	ms.DetailUrl = "/today/detail_ms.html"
	cs.DetailUrl = "/today/detail_cs.html"
	t.ExecuteTemplate(bw, "li.tmpl", pmap)

	sh.DetailUrl = "/his/d_sh_" + hisName + ".html"
	jk.DetailUrl = "/his/d_jk_" + hisName + ".html"
	rq.DetailUrl = "/his/d_rq_" + hisName + ".html"
	ys.DetailUrl = "/his/d_ys_" + hisName + ".html"
	bj.DetailUrl = "/his/d_bj_" + hisName + ".html"
	my.DetailUrl = "/his/d_my_" + hisName + ".html"
	ms.DetailUrl = "/his/d_ms_" + hisName + ".html"
	cs.DetailUrl = "/his/d_cs_" + hisName + ".html"
	t.ExecuteTemplate(bakbw, "li.tmpl", pmap)
	bw.Flush()
	f.Close()
	bakbw.Flush()
	bakf.Close()

	for k, vv := range pmap {
		//detail-sh
		shdarray := make([]*DP, 0, 5)
		for _, v := range vv.Dt.P {
			shdp := &DP{Tt: vv.Title, Intro: vv.Dt.Intro, P: v}
			shdarray = append(shdarray, shdp)
		}
		RenderDetail(shdarray, t, k)
	}

}

func RenderDetail(shdarray []*DP, t *template.Template, sh string) error {
	var hisName = time.Now().Format("2006-01-02")
	dshfilename := homeDir + "/changshi/today/detail_" + sh + ".html"
	if _, err := os.Stat(dshfilename); os.IsExist(err) {
		os.Remove(dshfilename)
	}
	dshf, _ := os.Create(dshfilename)
	dshfbw := bufio.NewWriter(dshf)
	err := t.ExecuteTemplate(dshfbw, "detail.tmpl", shdarray)
	if err != nil {
		return err
	}
	dshfbw.Flush()
	dshf.Close()

	dshhfilename := homeDir + "/changshi/his/d_" + sh + "_" + hisName + ".html"
	if _, err := os.Stat(dshhfilename); os.IsExist(err) {
		os.Remove(dshhfilename)
	}
	dshhf, _ := os.Create(dshhfilename)
	dshhfbw := bufio.NewWriter(dshhf)
	err1 := t.ExecuteTemplate(dshhfbw, "detail.tmpl", shdarray)
	if err1 != nil {
		return err
	}
	dshhfbw.Flush()
	dshhf.Close()
	return nil
}

func Dateconv(t time.Time) string {
	return t.Format("2006-01-02")
}

func RandBg() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(5) + 1
}
func Unescaped(x string) interface{} {
	return template.HTML(x)
}

type DP struct {
	Tt    string
	Intro string
	P     string
}
