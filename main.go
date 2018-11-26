//新增了index輸入關鍵字，還有使用mux成功
package main

import (
	"net/http"
	"fmt"
	"regexp"
	"io/ioutil"
	"strconv"
	"github.com/gorilla/mux"
)

var keyword string
var final [][]string
var max_page int
var max_pic int

func searchGet(keyword string) {
    resp, _ := http.Get("https://www.pixiv.net/search.php?word=" + keyword + "&s_mode=s_tag_full&order=date_d&p=1")
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    
	all_pid := S_pid_img(string(body))
	max_page = S_page(string(body))
	if(max_page > 10){
		max_page = 10
	}
    for i := 2;i <= max_page;i++{
        resp, _ := http.Get("https://www.pixiv.net/search.php?word=" + keyword + "&s_mode=s_tag_full&order=date_d&p=" + strconv.Itoa(i))
        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        temp := S_pid_img(string(body))
		for i, _ := range temp{
			all_pid = append(all_pid, temp[i])
		}
    }
	
	workGet(all_pid)

	fmt.Println("search for: " + keyword + " finish.")
}

func workGet(target []string){
	var to_zero [][]string
	final = to_zero
	for _, j := range target{
		resp, _ := http.Get("https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" + j)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		pname, aid, aname, like := A_pname_aid_aname_like(string(body))
		temp := []string{j, pname, aid, aname, like}
		
		final = append(final, temp)
	}
	final = sort(final)
}

func sort(data [][]string) [][]string{
	for i := 0;i < len(data) - 1;i++{
		for j := i + 1;j < len(data);j++{
			ii, _ := strconv.Atoi(data[i][4])
			ji, _ := strconv.Atoi(data[j][4])
			if(ii < ji){
				data[i], data[j] = data[j], data[i]
			}
		}
	}
	return data
}
const indexPage = `
<h1>歡迎使用未完成的P站人氣排序</h1>
<form method="post" action="/meddle">
    <label for="keyword">關鍵字：</label>
    <input type="text" id="keyword" name="keyword">
    <button type="submit">開始搜尋</button>
</form>
<h3>可能需要一分鐘的時間，請稍候</h3>
`
func search_handler(w http.ResponseWriter, r *http.Request) {
	//搜尋網址
	keyword = r.URL.Path[8:]
	searchGet(keyword)	
	http.Redirect(w, r, "/result/1", 302)
}


const head_1 = `<!DOCTYPE html>
<html>
	<head>
		<meta charset = "utf-8">
		<title>Pixiv search result</title>
		<link rel="stylesheet" href="P_search.css">
		<style type="text/css">
`
const head_2 = `</style>
		<script type="text/javascript">
`
const tail_1 =`		</script>
	</head>
	<body>
		`
const tail_2 =`	</body>
</html>`

const script_def = `var pid = new Array();
var pname = new Array();
var aid = new Array();
var aname = new Array();
var like = new Array();
`

func result_handler(w http.ResponseWriter, r *http.Request){
	script, _ := ioutil.ReadFile("source.js")
	css, _ := ioutil.ReadFile("source.css")
	//全js變數值寫入
	script_pid := "pid = ["
	script_pname := "pname = ["
	script_aid := "aid = ["
	script_aname := "aname = ["
	script_like := "like = ["
	for i, _ := range final{
		script_pid += ("\"" + final[i][0] + "\"")
		script_pname += ("\"" + final[i][1] + "\"")
		script_aid += ("\"" + final[i][2] + "\"")
		script_aname += ("\"" + final[i][3] + "\"")
		script_like += ("\"" + final[i][4] + "\"")
		if(i != len(final)-1){
			script_pid += ","
			script_pname += ","
			script_aid += ","
			script_aname += ","
			script_like += ","
		}
	}
	script_pid += "]\n"
	script_pname += "]\n"
	script_aid += "]\n"
	script_aname += "]\n"
	script_like += "]\n"
	page := r.URL.Path[8:]
	script_page := "var page = " + page + "\n"
	
	fmt.Fprintf(w, head_1 + string(css) + head_2)
	fmt.Fprintf(w, script_def + script_pid + script_pname + script_aid + script_aname + script_like + script_page)
	fmt.Fprintf(w, string(script))
	fmt.Fprintf(w, tail_1 + "Keyword = " + keyword + "<br>總頁數：" + strconv.Itoa(max_page) + "<br>\n" + "		<p id=\"result\"></p>")
	
	page_i, _ := strconv.Atoi(page)
	button := ""
	if page == "1"{
		button += "<a text-align=\"center\" href=\"/result/" + strconv.Itoa(page_i + 1) + "\">" + "<button>下一頁</button></a>\n"
	}else if page == "10"{
		button += "<a text-align=\"center\" href=\"/result/" + strconv.Itoa(page_i - 1) + "\">" + "<button>上一頁</button></a>\n"
	}else{
		button += "<a text-align=\"center\" href=\"/result/" + strconv.Itoa(page_i - 1) + "\">" + "<button>上一頁</button></a>\n"
		button += "<a text-align=\"center\" href=\"/result/" + strconv.Itoa(page_i + 1) + "\">" + "<button>下一頁</button></a>\n"
	}
	
	fmt.Fprintf(w, button + tail_2);
}




func index_handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, indexPage)
}

func meddle_Handler(w http.ResponseWriter, r *http.Request){
	keyword := r.FormValue("keyword")
	redirect_target := "/search/" + keyword
	http.Redirect(w, r, redirect_target, 302)
}

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", index_handler)
	router.HandleFunc("/search/{keyword:.+}", search_handler) // each request calls handler
	router.HandleFunc("/result/{page:[0-9]+}", result_handler)
	
	router.HandleFunc("/meddle", meddle_Handler).Methods("POST")
	
	http.Handle("/", router)
	fmt.Println("listen to http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}



//以下是regex相關

//從作品頁面抓下作品名稱、作者ID、作者名稱、讚數
func A_pname_aid_aname_like(input string)(string, string, string, string){
    re := regexp.MustCompile("<meta property=\"og:title\" content=\"「(.*?)」/「(.*?)」")
	re2 := regexp.MustCompile("<a href=\"member[.]php[?]id=([0-9]+)\">")
    re3 := regexp.MustCompile("<span class=\"views\">([0-9]+)</span></li></ul>")
    result := re.FindAllStringSubmatch(input, -1)
	aid := re2.FindAllStringSubmatch(input, -1)
	like := re3.FindAllStringSubmatch(input, -1)
	
	return result[0][1], aid[0][1], result[0][2], like[0][1]
}
//"
//從搜尋結果頁1回傳最大頁數
func S_page(input string) int{
    re := regexp.MustCompile("<span class=\"count-badge\">([0-9]+)")
    total := re.FindAllStringSubmatch(input, -1)
	
	pictures, _ := strconv.Atoi(total[0][1])
	max_pic = pictures
	max_page := (max_pic + 40)/40
	return max_page
}
//"
//從搜尋頁面回傳至多40個作品ID，與對應的40個預覽圖網址
func S_pid_img(input string) (out_pid []string){
    re := regexp.MustCompile("illustId&quot;:&quot;([0-9]+)&quot;")
    pid := re.FindAllStringSubmatch(input, -1)

	for i := 0;i < len(pid);i++{
        out_pid = append(out_pid, pid[i][1])
    }
	return
}