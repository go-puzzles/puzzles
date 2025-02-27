// File:		config.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package minio

import "errors"

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string

	Bucket string
}

func (c *MinioConfig) Validate() error {
	if c.Endpoint == "" || c.AccessKey == "" || c.SecretKey == "" || c.Bucket == "" {
		return errors.New("invalid minio config")
	}

	return nil
}
