package faceScore

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Scorer struct {
	appCode string
}

var (
	baseUrl = "http://faceapi.remarkdip.com/user/faceScore"
)

func NewScorer(appCode string) *Scorer {
	return &Scorer{
		appCode: appCode,
	}
}

func (s Scorer) LocalScore(imgPath string) (*Result, error) {
	//打开文件
	log.Println("本地文件评分")
	if encoded, err := FromLocal(imgPath); err != nil {
		return nil, fmt.Errorf("无法打开文件:%s", err.Error())
	} else {
		log.Println("base64编码")

		//base64 编码

		value := url.Values{
			"img_base64": []string{encoded},
		}
		//if _,err = fmt.Fprintf(buffer,"{\"img_base64\":\"%s\"",encoded); err != nil {
		//	return nil, fmt.Errorf("网络传输错误:%s", err.Error())
		//}

		//构建请求
		log.Println("构建请求")

		if request, err := s.buildRequest(true, strings.NewReader(value.Encode()), ""); err != nil {
			return nil, fmt.Errorf("构建请求错误:%s", err.Error())
		} else {
			log.Println("发出请求")
			//request.Write(os.Stdout)
			if resp, err := http.DefaultClient.Do(request); err != nil {
				return nil, fmt.Errorf("发送请求错误:%s", err.Error())
			} else {
				//处理应答
				return s.handResult(resp)
			}
		}
	}

}

func (s Scorer) WebScore(imgUrl string) (*Result, error) {

	//检查url
	if _, err := url.Parse(imgUrl); err != nil {
		return nil, fmt.Errorf("地址不合法:%s", err.Error())
	}
	//构建请求
	if request, err := s.buildRequest(false, nil, imgUrl); err != nil {
		return nil, fmt.Errorf("构建请求错误:%s", err.Error())
	} else {
		if resp, err := http.DefaultClient.Do(request); err != nil {
			return nil, fmt.Errorf("发送请求错误:%s", err.Error())
		} else {
			//处理应答
			return s.handResult(resp)
		}
	}
}

func (s Scorer) handResult(resp *http.Response) (*Result, error) {
	if resp != nil {
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		result := &Result{}

		if err := decoder.Decode(result); err != nil {
			return nil, fmt.Errorf("结果解析错误:%s", err.Error())
		} else {
			return result, nil
		}
	} else {
		return nil, errors.New("结果不能为空")
	}
}

func (s Scorer) buildRequest(local bool, reader io.Reader, url string) (*http.Request, error) {
	if local {
		if request, err := http.NewRequest("POST", baseUrl, reader); err == nil {
			request.Header.Set("Authorization", "APPCODE "+s.appCode)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			return request, nil
		} else {
			return nil, err

		}
	} else {
		if request, err := http.NewRequest("GET", baseUrl, nil); err == nil {
			q := request.URL.Query()
			q.Set("img_url", url)
			request.URL.RawQuery = q.Encode()
			request.Header.Set("Authorization", "APPCODE "+s.appCode)
			return request, nil
		} else {
			return nil, err

		}
	}
}
func FromLocal(fname string) (string, error) {
	var b bytes.Buffer

	fileExists, _ := exists(fname)
	if !fileExists {
		return "", fmt.Errorf("File does not exist\n")
	}

	file, err := os.Open(fname)
	if err != nil {
		return "", fmt.Errorf("Error opening file\n")
	}

	_, err = b.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file to buffer\n")
	}

	return FromBuffer(b), nil
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func FromBuffer(buf bytes.Buffer) string {
	enc := encode(buf.Bytes())
	mime := http.DetectContentType(buf.Bytes())

	return format(enc, mime)
}
func encode(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
}
func format(enc []byte, mime string) string {
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;base64,%s", mime, enc)
	default:
	}

	return fmt.Sprintf("data:image/png;base64,%s", enc)
}
