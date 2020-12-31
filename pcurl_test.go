package pcurl

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

type H map[string]interface{}

func createGeneral(data string) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.Any("/*test", func(c *gin.Context) {
			if len(data) > 0 {
				c.String(200, data)
			}
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func createGeneral2() *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.Any("/*test", func(c *gin.Context) {
			io.Copy(c.Writer, c.Request.Body)
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_Method(t *testing.T) {
	methodServer := func() *httptest.Server {
		router := func() *gin.Engine {
			router := gin.New()

			router.DELETE("/", func(c *gin.Context) {
				c.String(200, "DELETE")
			})

			router.GET("/", func(c *gin.Context) {
				c.String(200, "GET")
			})
			return router
		}()

		return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	}

	need := []string{"DELETE", "DELETE", "GET"}

	for index, curlStr := range []string{
		`curl -X DELETE -G `,
		`curl -G -XDELETE `,
		`curl -G `,
	} {
		ts := methodServer()
		req, err := ParseAndRequest(curlStr + ts.URL)

		assert.NoError(t, err)

		got := ""
		err = gout.New().SetRequest(req).BindBody(&got).Do()

		assert.Equal(t, got, need[index])
		assert.NoError(t, err)
	}
}

func Test_URL(t *testing.T) {
	type testURL struct {
		curl []string
		need string
	}

	for _, urlData := range []testURL{
		testURL{
			curl: []string{"curl", "-X", "POST"},
			need: "1",
		},
		testURL{
			curl: []string{"curl", "-X", "POST"},
			need: "2",
		},
	} {

		code := 0
		// 创建测试服务端
		ts := createGeneral("1")
		ts2 := createGeneral("2")

		// 解析curl表达式
		var curl []string
		if urlData.need == "1" {
			curl = append(urlData.curl, "--url", ts2.URL, ts.URL)
		} else {
			curl = append(urlData.curl, ts.URL, "--url", ts2.URL)

		}

		req, err := ParseSlice(curl).Request()
		assert.NoError(t, err)

		s := ""
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindBody(&s).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
		assert.Equal(t, urlData.need, s)
	}
}

func TestParseString(t *testing.T) {
	data := `curl 'http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login' \
  -H 'Connection: keep-alive' \
  -H 'Cache-Control: max-age=0' \
  -H 'Upgrade-Insecure-Requests: 1' \
  -H 'Origin: http://www.ewebeditor.net' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36' \
  -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' \
  -H 'Referer: http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login' \
  -H 'Accept-Language: zh-CN,zh;q=0.9' \
  -H 'Cookie: ASPSESSIONIDAATBTDQT=CLIDJDJCEDECGOHNFHBGBCCA; IPCity=%E5%B9%BF%E5%B7%9E' \
  --data-raw 'h=www.ewebeditor.net&usr=admin&pwd=admin' \
  --insecure`
	req, _ := ParseAndRequest(data)

	resp := ""
	err := gout.New().Debug(false).SetRequest(req).BindBody(&resp).Do()

	fmt.Println(err, "resp.size = ", len(resp))
}
