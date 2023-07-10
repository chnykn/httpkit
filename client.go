// Copyright 2019-2023 chnykn@gmail.com All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpkit

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BeforeReqFunc func(httpReq *http.Request) error

type Client struct {
	beforeReq  BeforeReqFunc
	httpClient *http.Client
}

func NewClient(beforeReq BeforeReqFunc) *Client {
	res := &Client{
		beforeReq:  beforeReq,
		httpClient: &http.Client{},
	}
	return res
}

//-----------------------------------------------

func (o *Client) Request(url string, method string, v ...any) (int, error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return -1, err
	}

	if o.beforeReq != nil {
		err = o.beforeReq(req)
		if err != nil {
			return -1, err
		}
	}

	var respPtr any = nil

	for _, vi := range v {
		switch vv := vi.(type) {

		case *ReqAccept:
			setReqAccecpt(req, vv.string)

		case *ReqQuery:
			setReqQuery(req, *vv)

		case *ReqBody:
			setReqBody(req, vv)

		case *ReqForm:
			setReqForm(req, vv)

		case *RespBody:
			respPtr = vv.Ptr
		}
	}

	//-----------------

	response, err := o.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	var statusCode = response.StatusCode

	defer response.Body.Close()
	respBody, _ := io.ReadAll(response.Body)

	//statusCode 200 201 204
	if (statusCode != http.StatusOK) && (statusCode != http.StatusCreated) && (statusCode != http.StatusNoContent) {
		err = fmt.Errorf("netkit.doRequest ERR: status=%s, body=%s", response.Status, string(respBody))
		return statusCode, err
	}

	//statusCode 200
	if (statusCode == http.StatusOK) && (respPtr != nil) {

		//If the type is either *[]byte or *bytes.Buffer, return directly.
		switch resp := respPtr.(type) {
		case *[]byte:
			*resp = respBody
			return statusCode, nil
		case *bytes.Buffer:
			resp.Write(respBody)
			return statusCode, nil
		}

		contentType := req.Header.Get(kHeaderContentType)

		//Unmarshal to args.result based on the Accept.
		if strings.HasPrefix(contentType, JsonContentType) {
			return statusCode, json.Unmarshal(respBody, respPtr)
		} else if strings.HasPrefix(contentType, XmlContentType) {
			return statusCode, xml.Unmarshal(respBody, respPtr)
		} else {
			// If it is neither JSON nor XML,
			// there is no need to parse it into args.result.
		}
	}

	return statusCode, nil
}

/*
Status code    Description
-------------------------------------------------
200			OK
201			Created
204			No Content
401			Unauthorized
403			Forbidden
404			Not Found
405			Method Not Allowed
500			Internal Server Error
*/

//-----------------------------------------------

// Get return (requst response'StatusCode, error)
func (o *Client) Get(url string, v ...any) (int, error) {
	return o.Request(url, MethodGet, v...)
}

// Delete return (requst response'StatusCode, error)
func (o *Client) Delete(url string, v ...any) (int, error) {
	return o.Request(url, MethodDelete, v...)
}

// Put return (requst response'StatusCode, error)
func (o *Client) Put(url string, v ...any) (int, error) {
	return o.Request(url, MethodPut, v...)
}

// Post return (requst response'StatusCode, error)
func (o *Client) Post(url string, v ...any) (int, error) {
	return o.Request(url, MethodPost, v...)
}
