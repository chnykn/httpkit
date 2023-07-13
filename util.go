// Copyright 2019-2023 chnykn@gmail.com All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpkit

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func setReqAccecpt(req *http.Request, accept string) {
	req.Header.Set(kHeaderAccept, accept)
}

func setReqQuery(req *http.Request, query ReqQuery) {
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func setReqBody(req *http.Request, body *ReqBody) {
	req.Body = io.NopCloser(body.Body)
	req.ContentLength = body.ConentLength
	req.Header.Set(kHeaderContentType, body.ConentType)
}

func setReqForm(req *http.Request, form *ReqForm) {
	pr, pw := io.Pipe()
	bodyWriter := multipart.NewWriter(pw)

	go func() {
		for key, value := range form.Fields {
			bodyWriter.WriteField(key, value)
		}

		n := 0
		for _, ufile := range form.UploadFiles {

			ofile, err := os.Open(ufile.FilePath)
			if err != nil {
				continue
			}
			defer ofile.Close()

			if ufile.FieldName == "" {
				n++
				ufile.FieldName = "file" + strconv.Itoa(n)
			}

			fileName := filepath.Base(ufile.FilePath)
			fileWriter, err := bodyWriter.CreateFormFile(ufile.FieldName, fileName)
			if err != nil {
				continue
			}

			io.Copy(fileWriter, ofile)
		}
		bodyWriter.Close()
		pw.Close()
	}()

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Body = io.NopCloser(pr)
}

/*
func setReqForm(req *http.Request, form *ReqForm) {
	var buffer bytes.Buffer
	bodyWriter := multipart.NewWriter(&buffer)

	for key, value := range form.Fields {
		_ = bodyWriter.WriteField(key, value)
	}

	n := 0
	for i := 0; i < len(form.UploadFiles); i++ {
		uploadFile := form.UploadFiles[i]

		file, err := os.Open(uploadFile.FilePath)
		if err != nil {
			continue
		}
		defer file.Close()

		if uploadFile.FieldName == "" {
			n++
			uploadFile.FieldName = "file" + strconv.Itoa(n)
		}

		fileName := filepath.Base(uploadFile.FilePath)
		fileWriter, err := bodyWriter.CreateFormFile(uploadFile.FieldName, fileName)
		if err != nil {
			continue
		}

		io.Copy(fileWriter, file)
	}
	bodyWriter.Close()

	req.Header.Set(kHeaderContentType, bodyWriter.FormDataContentType())
	req.Body = io.NopCloser(&buffer)
}
*/

//-------------------------------------------------------

//func getError(status int, respText []byte) error {
//	errName, ok := statusErrorName[status]
//	if !ok {
//		errName = fmt.Sprintf("unexpected error, status: %d", status)
//	}
//
//	detail := string(respText)
//	return fmt.Errorf("%s\ndetail:%s\n", errName, detail)
//}
