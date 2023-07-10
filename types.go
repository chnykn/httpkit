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
	Body         io.Reader
	ConentType   string
	ConentLength int64
}

//-----------------------

type UploadFile struct {
	FieldName string
	FilePath  string
}

type ReqForm struct {
	Fields      map[string]string
	UploadFiles []*UploadFile
}

//-----------------------

// RespBody Used to parse the JSON response in the response.Body into a corresponding
// object pointer (which must be a pointer type) after executing the request.
type RespBody struct {
	Ptr any //  *[]byte , *bytes.Buffer,  or pointer of other object
}

//=========================================================

func NewReqAccept(value string) *ReqAccept {
	return &ReqAccept{value}
}

func NewReqQuery(value map[string]string) ReqQuery {
	return ReqQuery(value)
}

func NewReqForm(fields map[string]string, filePath string) *ReqForm {
	uploadFiles := []*UploadFile{
		&UploadFile{
			FieldName: "file",
			FilePath:  filePath,
		},
	}

	return &ReqForm{
		Fields:      fields,
		UploadFiles: uploadFiles,
	}
}

func NewReqFormX(fields map[string]string, uploadFiles []*UploadFile) *ReqForm {
	return &ReqForm{
		Fields:      fields,
		UploadFiles: uploadFiles,
	}
}

//----------------------------------------

func NewReqBody(body io.Reader, conentLength int64, contentType string) *ReqBody {
	res := &ReqBody{Body: body, ConentType: contentType, ConentLength: conentLength}
	return res
}

func XmlReqBody(data any) *ReqBody {
	bts, _ := xml.Marshal(data)
	return NewReqBody(bytes.NewBuffer(bts), int64(len(bts)), XmlContentType)
}

func JsonReqBody(data any) *ReqBody {
	bts, _ := json.Marshal(data)
	return NewReqBody(bytes.NewBuffer(bts), int64(len(bts)), JsonContentType)
}

func BytesReqBody(data []byte) *ReqBody {
	rd := bytes.NewReader(data)
	return NewReqBody(io.NopCloser(rd), int64(len(data)), StreamContentType) //+"; charset=UTF-8"
}

func BufferReqBody(data bytes.Buffer) *ReqBody {
	return BytesReqBody(data.Bytes()) //+"; charset=UTF-8"
}

//----------------------------------------

func NewRespBody(ptr any) *RespBody {
	return &RespBody{Ptr: ptr}
}
