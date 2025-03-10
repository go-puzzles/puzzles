// File:		oss.go
// Created by:	Hoven
// Created on:	2024-11-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package poss

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

type IOSS interface {
	UploadFile(ctx context.Context, size int64, dir, objName string, obj io.Reader, tags map[string]string) (uri string, err error)
	GetFile(ctx context.Context, objName string, w io.Writer) error
	CheckFileExists(ctx context.Context, objName string) (bool, error)
	PresignedGetObject(ctx context.Context, objName string, expires time.Duration) (*url.URL, error)
	ProxyPresignedGetObject(objName string, rw http.ResponseWriter, req *http.Request)
}
