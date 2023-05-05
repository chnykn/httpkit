// Copyright 2019-2023 chnykn@gmail.com All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpkit

import (
	"io"
)

type reqArgs struct {
	URL    string
	Method string
	Accept string
	Query  map[string]string

	ReqBody       io.Reader
	ContentType   string
	ContentLength int

	RespPtr any
}

func newReqArgs(url string, method string, v ...any) *reqArgs {

	res := &reqArgs{
		URL:           url,
		Method:        method,
		Accept:        JsonContentType,
		Query:         nil,
		ContentType:   "",
		ContentLength: 0,
		ReqBody:       nil,
		RespPtr:       nil,
	}

	for _, vi := range v {
		switch vv := vi.(type) {

		case *ReqAccept:
			res.Accept = vv.string

		case *ReqQuery:
			res.Query = make(map[string]string)
			for key, val := range *vv {
				res.Query[key] = val
			}

		case *ReqBody:
			res.ReqBody = vv.Body
			res.ContentType = vv.Type
			res.ContentLength = vv.Length

		case *RespBody:
			res.RespPtr = vv.Ptr

			//-----------------

		case ReqAccept:
			res.Accept = vv.string

		case ReqQuery:
			res.Query = make(map[string]string)
			for key, val := range vv {
				res.Query[key] = val
			}

		case ReqBody:
			res.ReqBody = vv.Body
			res.ContentType = vv.Type
			res.ContentLength = vv.Length

		case RespBody:
			res.RespPtr = vv.Ptr
		}

	}

	return res
}

//func getError(status int, respText []byte) error {
//	errName, ok := statusErrorName[status]
//	if !ok {
//		errName = fmt.Sprintf("unexpected error, status: %d", status)
//	}
//
//	detail := string(respText)
//	return fmt.Errorf("%s\ndetail:%s\n", errName, detail)
//}
