package tools

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

func Post(c context.Context, client *resty.Client, url string, header, param map[string]string, body map[string]any, checkResCode bool) (resp *resty.Response, err error) {
	if header == nil {
		header = make(map[string]string)
		header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
	} else {
		if _, ok := header["Content-Type"]; !ok {
			header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
		}
		if _, ok := header["User-Agent"]; !ok {
			header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
		}
	}

	req := client.R().SetHeaders(header)

	if param != nil {
		req.SetQueryParams(param)
	}
	if body != nil {
		req.SetBody(body)
	}

	resp, err = req.Post(url)

	if err != nil {
		return
	}
	if resp.StatusCode() != 200 {
		err = fmt.Errorf("status code is: %d", resp.StatusCode())
		return
	}

	if checkResCode {
		if gjson.Get(resp.String(), "result_code").String() != "0" {
			err = fmt.Errorf("result_code is: %s, result_message is: %s", gjson.Get(resp.String(), "result_code").String(), gjson.Get(resp.String(), "result_message").String())
			return
		}
	}

	return
}
