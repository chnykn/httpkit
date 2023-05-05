// Copyright 2019-2023 chnykn@gmail.com All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpkit

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
)

//=========================================================

//The following types are used as v ...any parameter members
//  passed to the NetKit.Get/Delete/Put/Post method calls.

// ReqAccept Used to set the Accept value in the http.Request header.
type ReqAccept struct {
	string
}

// ReqQuery Used to set the URL.RawQuery value in the http.Request object.
type ReqQuery map[string]string

// ReqBody Used to set the body and contentType values in the http.Request object.
type ReqBody struct {
	Type   string
	Body   io.Reader
	Length int
}

// RespBody Used to parse the JSON response in the response.
type RespBody struct {
	buff *bytes.Buffer
}

//===================================

func NewReqAccept(value string) *ReqAccept {
	return &ReqAccept{value}
}

func NewReqQuery(value map[string]string) ReqQuery {
	return ReqQuery(value)
}

//----------------------------------------

func NewReqBody(body io.Reader, conentLength int, contentType string) *ReqBody {
	res := &ReqBody{Body: body, Type: contentType, Length: conentLength}
	return res
}

func XmlReqBody(data any) *ReqBody {
	bts, _ := xml.Marshal(data)
	return NewReqBody(bytes.NewBuffer(bts), len(bts), XmlContentType)
}

func JsonReqBody(data any) *ReqBody {
	bts, _ := json.Marshal(data)
	return NewReqBody(bytes.NewBuffer(bts), len(bts), JsonContentType)
}

func BytesReqBody(data []byte) *ReqBody {
	rd := bytes.NewReader(data)
	return NewReqBody(io.NopCloser(rd), len(data), JsonContentType) //+"; charset=UTF-8"
}

func BufferReqBody(data bytes.Buffer) *ReqBody {
	return BytesReqBody(data.Bytes()) //+"; charset=UTF-8"
}

//----------------------------------------

func NewRespBody(buff *bytes.Buffer) *RespBody {
	return &RespBody{buff: buff}
}