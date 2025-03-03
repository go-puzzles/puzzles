// File:		minio.go
// Created by:	Hoven
// Created on:	2024-11-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/poss"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

var _ poss.IOSS = (*MinioOss)(nil)

type MinioOss struct {
	*MinioConfig
	client *minio.Client
}

func NewMinioOss(conf *MinioConfig) *MinioOss {
	m := &MinioOss{
		MinioConfig: conf,
	}

	discoverAddr := discover.GetAddress(conf.Endpoint)
	conf.Endpoint = discoverAddr

	var err error
	m.client, err = minio.New(discoverAddr, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: false,
	})
	plog.PanicError(err)

	exists, err := m.client.BucketExists(context.TODO(), conf.Bucket)
	plog.PanicError(err)
	if !exists {
		plog.Fatalf("bucket %s not exists", conf.Bucket)
	}

	return m
}

func (m *MinioOss) Init(router gin.IRouter) {
	minioGroup := router.Group("minio")
	minioGroup.GET(":sourceType/:sourceName", pgin.RequestHandler(m.getMinioSourceHandler))
}

type GetMinioSourceRequest struct {
	SourceType string `uri:"sourceType" binding:"required"`
	SourceName string `uri:"sourceName" binding:"required"`
}

func (m *MinioOss) getMinioSourceHandler(ctx *gin.Context, req *GetMinioSourceRequest) {
	objName := fmt.Sprintf("%s/%s", req.SourceType, req.SourceName)
	if err := m.GetFile(ctx, objName, ctx.Writer); err != nil {
		plog.Errorc(ctx, "get minio source(%s) error: %v", objName, err)
		ctx.Status(http.StatusBadRequest)
		return
	}
}

func (m *MinioOss) getFileExt(file string) string {
	return filepath.Ext(file)
}

func (m *MinioOss) generateObjName(obj string) string {
	fileName := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), uuid.New().String())
	ext := m.getFileExt(obj)
	return fileName + ext
}

func (m *MinioOss) UploadFile(ctx context.Context, size int64, dir, objName string, obj io.Reader) (uri string, err error) {
	rawObjName := m.generateObjName(objName)

	putOpt := minio.PutObjectOptions{
		UserTags: map[string]string{},
	}

	newObjName := fmt.Sprintf("%s/%s", dir, rawObjName)
	_, err = m.client.PutObject(ctx, m.Bucket, newObjName, obj, size, putOpt)
	if err != nil {
		return "", errors.Wrap(err, "uploadMinio")
	}

	// 1731850656800-d887240b-0177-44c7-853d-69f14b7cf874.jpeg
	return rawObjName, nil
}

func (m *MinioOss) GetFile(ctx context.Context, objName string, w io.Writer) error {
	object, err := m.client.GetObject(ctx, m.Bucket, objName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()

	_, err = io.Copy(w, object)
	if err != nil {
		return errors.Wrap(err, "getMinioObject")
	}

	return nil
}

func (m *MinioOss) CheckFileExists(ctx context.Context, objName string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.Bucket, objName, minio.StatObjectOptions{})
	if err != nil {
		minioErrResp := new(minio.ErrorResponse)
		if errors.As(err, &minioErrResp) && minioErrResp.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, errors.Wrap(err, "statObject")
	}

	return true, nil
}

func (m *MinioOss) PresignedGetObject(ctx context.Context, objName string, expires time.Duration) (*url.URL, error) {
	u, err := m.client.PresignedGetObject(ctx, m.Bucket, objName, expires, url.Values{})
	if err != nil {
		return nil, errors.Wrap(err, "presignedGetObject")
	}

	plog.Debugf("presigned origin url: %s", u.String())

	return u, nil
}

func (m *MinioOss) ProxyPresignedGetObject(objName string, rw http.ResponseWriter, req *http.Request) {
	minioEndpoint := m.client.EndpointURL()
	proxy := httputil.NewSingleHostReverseProxy(minioEndpoint)
	req.URL.Path = fmt.Sprintf("%s/%s", m.Bucket, objName)
	req.Host = minioEndpoint.Host

	plog.Debugf("Proxy get object url: %s", req.URL.String())
	proxy.ServeHTTP(rw, req)
}
