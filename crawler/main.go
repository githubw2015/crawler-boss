package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	seleniumPath    = `C:\chromedriver_win32\selenium-server-standalone-3.9.1.jar`
	geckoDriverPath = `C:\chromedriver_win32\chromedriver.exe`
	port            = 9515
)

var (
	DB       *gorm.DB
	username = "root"
	password = "root"
	dbName   = "test"
)

var searKeywords = []string{
	"golang",
	"php",
	"Python",
	"Java",
}

var cityMap = map[int]string{
	101020100: "上海",
	//101010100: "北京",
	//101280100: "广州",
	//101280600: "深圳",
	//101210100: "杭州",
}

var proxyIps = []string{
	"http://120.38.241.162:4510",
	"http://58.241.203.160:4545",
	"http://180.125.107.166:4536",
	"http://180.125.33.225:4557",
	"http://124.94.250.26:4560",
	"http://42.54.90.13:4550",
	"http://27.44.216.205:4545",
	"http://180.125.2.213:4567",
	"http://117.60.242.32:4547",
}

func init() {
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbName))
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}

	DB.SingularTable(true)
}

var wg sync.WaitGroup

func main() {
	//初始化基本参数
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(ioutil.Discard),        // Output debug information to STDERR.
	}
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	defer service.Stop()

	for index, val := range cityMap {
		for _, item := range searKeywords {
			wg.Add(1)
			go func(item string, index int, val string) {
				if err != nil {
					panic(err) // panic is used only as an example and is not otherwise recommended.
				}
				//打开 chrome 浏览器
				caps := selenium.Capabilities{"browserName": "chrome"}
				//禁止图片加载，加快渲染速度
				imagCaps := map[string]interface{}{
					"profile.managed_default_content_settings.images": 2,
				}
				rand.Seed(time.Now().Unix())
				proxyIndex := rand.Intn(len(proxyIps))
				chromeCaps := chrome.Capabilities{
					Prefs: imagCaps,
					Path:  "",
					Args: []string{
						"--headless",
						"--start-maximized",
						//"--window-size=1200x600",
						"--no-sandbox",
						"--user-agent=" + GetRandomUserAgent(),
						"--disable-gpu",
						"--disable-impl-side-painting",
						"--disable-gpu-sandbox",
						"--disable-accelerated-2d-canvas",
						"--disable-accelerated-jpeg-decoding",
						"--test-type=ui",
						"--proxy-server=" + proxyIps[proxyIndex],
					},
				}

				//以上是设置浏览器参数
				caps.AddChrome(chromeCaps)
				//打开 chrome 浏览器
				wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
				if err != nil {
					panic(err)
				}
				//wd.AddCookie(&selenium.Cookie{
				//	Name:  "__zp_stoken__",
				//	Value: "__fid=c2b051dc22170700021a31d7606054c0; wt2=DLzLXzYb7kUaCJMioMFtrwMhk6eQlRn81wUGE0NkP6lHW1BcFpJSAPKv89ZXnViD933HW6_mmU-_734s4nYYbMg~~; _bl_uid=C4kRkrpp7k7s80ge9i5wrwst599d; acw_tc=0bdd34b616265981976475631e01e066c1d589888f9a880aaa421e2f892e21; lastCity=101020100; __zp_seo_uuid__=c18ed6d2-ab6f-4214-896a-c55a0d9bc586; __c=1626599142; __g=-; __l=r=https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3D_ycLarYk8_yn0W3nbwH-I2939KNJrnyYRn7Ahn43fZp1bMhDMqRI1cFTkozRfT9F%26wd%3D%26eqid%3Db762e6bb0008e8df0000000560f3eedf&l=%2Fwww.zhipin.com%2Fshanghai%2F&s=1&g=&s=3&friend_source=0; __a=11523211.1626525304.1626597684.1626599142.97.7.1.97; Hm_lvt_194df3105ad7148dcf2b98a91b5e727a=1626594121,1626594699,1626597681,1626599142; Hm_lpvt_194df3105ad7148dcf2b98a91b5e727a=1626599142; __zp_stoken__=83bfcZ1IkQXl8FiE1czd0FAMGQnRYLjhrJWQoUQVFMUZMYWt2Y1RxQGEYKGYFR1gFPEx0BnQhKQIuAygUBl9TeUptCy9YbCw4VR1VYlYZOUdRbzpNUDRMPFZxAV0zMih4DG9kO30kVnYNQTo0",
				//})
				var count = 0
				for i := 1; ; i++ {
					urls := `https://www.zhipin.com/c` + strconv.Itoa(index) + `/?query=` + item + `&page=` + strconv.Itoa(i)
					fmt.Println(urls)
					//加载网页
					if err := wd.Get(urls); err != nil {
						panic(err)
					}
					time.Sleep(time.Second * 10)
					jsRt, err := wd.ExecuteScript("return document.readyState", nil)
					if err != nil {
						log.Println("exe js err", err)
					}
					fmt.Println("jsRt", jsRt)
					if jsRt != "complete" {
						log.Println(item + "网页加载未完成" + strconv.Itoa(i))
						time.Sleep(time.Second * 5)
					}
					// next disabled
					// 获取网站内容
					var frameHtml string
					frameHtml, err = wd.PageSource()
					if err != nil {
						log.Println(err)
						return
					}
					//解析 html 文件
					var doc *goquery.Document
					doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHtml)))
					if err != nil {
						log.Println(err)
						return
					}
					var Workexperience, Education, rongzi, staffNumber string
					doc.Find("#main ul li").Each(func(i int, context *goquery.Selection) {
						jobName := trimSpase(context.Find("span[class=\"job-name\"]").Text())
						salary := trimSpase(context.Find("span[class=\"red\"]").Text())
						href, _ := context.Find("div[class=\"info-primary\"] a").Attr("href")
						company := trimSpase(context.Find("div[class=\"info-company\"] h3").Text())
						address := trimSpase(context.Find("span[class=\"job-area\"]").Text())
						worklimit, _ := context.Find("div[class=\"job-limit clearfix\"] p").Html()
						industry := trimSpase(context.Find("div[class=\"info-company\"] a[class=\"false-link\"]").Text())

						data1 := strings.Split(worklimit, "<em class=\"vline\"></em>")
						for index, val := range data1 {
							if index == 0 {
								Workexperience = trimSpase(val)
							} else if index == 1 {
								Education = trimSpase(val)
							}
						}

						href = "https://www.zhipin.com" + href
						rognstuff, _ := context.Find("div[class=\"info-company\"] p").Html()
						data2 := strings.Split(rognstuff, "<em class=\"vline\"></em>")
						for index, val := range data2 {
							if index == 1 {
								rongzi = trimSpase(val)
							} else if index == 2 {
								staffNumber = trimSpase(val)
							}
						}
						if jobName != "" {
							sp := SpBossJobs{
								JobName:        jobName,
								Salary:         salary,
								Href:           href,
								JobType:        item,
								City:           val,
								CompanyName:    company,
								CompanyAddress: address,
								WorkYears:      Workexperience,
								Education:      Education,
								CompanyLabel:   industry,
								FinancingStage: rongzi,
								StaffNumber:    staffNumber,
							}
							sp.Add()
							count++
						}
					})
					_, errs := wd.FindElement(selenium.ByCSSSelector, "a[class='next disabled']")
					if errs == nil {
						fmt.Println(item + "找到隐藏;抓取总数:" + strconv.Itoa(count))
						break
					}
				}
				wg.Done()
				wd.Quit() // 关闭浏览器
			}(item, index, val)
		}
	}

	wg.Wait()
	fmt.Println("结束")
}

func trimSpase(str string) string {
	strs := strings.Replace(str, " ", "", -1)
	strs = strings.Replace(strs, "\n", "", -1)
	return strs
}

// boss招聘信息表
type SpBossJobs struct {
	Id             uint   `db:"id"`
	JobName        string `db:"job_name"`        //工作名称
	Salary         string `db:"salary"`          //薪资
	City           string `db:"city"`            //城市
	JobType        string `db:"job_type"`        //薪资
	Href           string `db:"href"`            //详情连接
	CompanyName    string `db:"company_name"`    //公司名称
	CompanyAddress string `db:"company_address"` //公司地址
	WorkYears      string `db:"work_years"`      //工作年限
	Education      string `db:"education"`       //学历要求
	CompanyLabel   string `db:"company_label"`   //公司所属行业
	FinancingStage string `db:"financing_stage"` //融资阶段
	StaffNumber    string `db:"staff_number"`    //公司规模-员工人数
}

//添加数据
func (sp *SpBossJobs) Add() {
	err := DB.Create(sp).Error
	if err != nil {
		fmt.Println("创建失败")
	}
}

var userAgentList = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.109 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/603.2.5 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.5",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
	"Mozilla/5.0 (Windows NT 6.1; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.109 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.109 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0",
	"Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.1 Safari/603.1.30",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 5.1; rv:52.0) Gecko/20100101 Firefox/52.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.109 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:52.0) Gecko/20100101 Firefox/52.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/58.0.3029.110 Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/603.2.5 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.5",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:45.0) Gecko/20100101 Firefox/45.0",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36 OPR/46.0.2597.32",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/59.0.3071.109 Chrome/59.0.3071.109 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:53.0) Gecko/20100101 Firefox/53.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 OPR/45.0.2552.898",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36 OPR/46.0.2597.39",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/601.7.7 (KHTML, like Gecko) Version/9.1.2 Safari/601.7.7",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; rv:52.0) Gecko/20100101 Firefox/52.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36",
}

func GetRandomUserAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return userAgentList[r.Intn(len(userAgentList))]
}
