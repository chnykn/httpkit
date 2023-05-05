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
/*
	case consts.MethodGet, consts.MethodDelete:
		httpReq = o.newRequest(args.URL, args.Method, args.Accept, "", nil)
	case consts.MethodPost, consts.MethodPut:
		httpReq = o.newRequest(args.URL, args.Method, args.Accept, args.ContentType, args.ReqBody)
	default:
		err := fmt.Errorf("netkit.doRequest ERR: unrecognized http request method '%s'", args.Method)
		return -1, err
	}
*/
func (o *Client) newRequest(args *reqArgs) (*http.Request, error) {

	result, err := http.NewRequest(args.Method, args.URL, args.ReqBody)
	if err != nil {
		return nil, err
	}

	if args.ReqBody != nil {
		result.Header.Set(headerContentType, args.ContentType)
		if args.ContentLength > 0 {
			result.ContentLength = int64(args.ContentLength)
		}
	}

	if args.Accept != "" {
		result.Header.Set(headerAccept, args.Accept)
	}

	return result, nil
}

// return (requst response'StatusCode, error)
func (o *Client) doRequest(args *reqArgs) (int, error) {

	httpReq, err := o.newRequest(args)
	if err != nil {
		return -1, err
	}

	if len(args.Query) > 0 {
		q := httpReq.URL.Query()
		for k, v := range args.Query {
			q.Add(k, v)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	if o.beforeReq != nil {
		err = o.beforeReq(httpReq)
		if err != nil {
			return -1, err
		}
	}

	//-----------------

	response, err := o.httpClient.Do(httpReq)
	if err != nil {
		return 0, err
	}

	var statusCode = response.StatusCode

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	//statusCode 200 201 204
	if (statusCode != http.StatusOK) && (statusCode != http.StatusCreated) && (statusCode != http.StatusNoContent) {
		err = fmt.Errorf("netkit.doRequest ERR: status=%s, body=%s", response.Status, string(body))
		return statusCode, err
	}

	//statusCode 200
	if (statusCode == http.StatusOK) && (args.RespPtr != nil) {

		//If the type is either *[]byte or *bytes.Buffer, return directly.
		switch resp := args.RespPtr.(type) {
		case *[]byte:
			*resp = body
			return statusCode, nil
		case *bytes.Buffer:
			resp.Write(body)
			return statusCode, nil
		}

		//Unmarshal to args.result based on the Accept.
		if strings.HasPrefix(args.Accept, JsonContentType) {
			return statusCode, json.Unmarshal(body, args.RespPtr)
		} else if strings.HasPrefix(args.Accept, XmlContentType) {
			return statusCode, xml.Unmarshal(body, args.RespPtr)
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

// Request Get return (requst response'StatusCode, error)
func (o *Client) Request(url string, method string, v ...any) (int, error) {
	args := newReqArgs(url, method, v...)
	return o.doRequest(args)
}

// Get return (requst response'StatusCode, error)
func (o *Client) Get(url string, v ...any) (int, error) {
	args := newReqArgs(url, MethodGet, v...)
	return o.doRequest(args)
}

// Delete return (requst response'StatusCode, error)
func (o *Client) Delete(url string, v ...any) (int, error) {
	args := newReqArgs(url, MethodDelete, v...)
	return o.doRequest(args)
}

// Put return (requst response'StatusCode, error)
func (o *Client) Put(url string, v ...any) (int, error) {
	args := newReqArgs(url, MethodPut, v...)
	return o.doRequest(args)
}

// Post return (requst response'StatusCode, error)
func (o *Client) Post(url string, v ...any) (int, error) {
	args := newReqArgs(url, MethodPost, v...)
	return o.doRequest(args)
}
