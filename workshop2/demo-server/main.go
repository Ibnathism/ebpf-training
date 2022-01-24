/*
 * Copyright 2018- The Pixie Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort    = "8080"
	maxPayloadSize = 10 * 1024 * 1024 // 10 MB
)

var (
	// source is a static, global rand object.
	source          *rand.Rand
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$"
	portEnvVariable = "PORT"
)

// customResponse holds the requested size for the response payload.
type customResponse struct {
	Size int `json:"size"`
}

func init() {
	source = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.POST("/customResponse", postCustomResponse)

	port := os.Getenv(portEnvVariable)
	if port == "" {
		port = defaultPort
	}

	fmt.Printf("listening on 0.0.0.0:%s\n", port)
	if err := engine.Run(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatal(err)
	}
}

func postCustomResponse(context *gin.Context) {
	var customResp customResponse
	if err := context.BindJSON(&customResp); err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if customResp.Size > maxPayloadSize {
		_ = context.AbortWithError(http.StatusBadRequest, fmt.Errorf("requested size %d is bigger than max allowed %d", customResp, maxPayloadSize))
		return
	}

	context.JSON(http.StatusOK, map[string]string{"answer": randStringBytes(customResp.Size)})
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[source.Intn(len(letterBytes))]
	}
	return string(b)
}
