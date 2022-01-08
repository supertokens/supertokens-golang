/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

/*
 The contents of this file are taken from

 https://github.com/zenazn/goji/blob/master/web/mutil/writer_proxy.go

 and

 https://medium.com/@cep21/interface-wrapping-method-erasure-c523b3549912
*/

package supertokens

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type DoneWriter interface {
	http.ResponseWriter
	IsDone() bool
}

// WrapWriter wraps an http.ResponseWriter, returning a proxy that allows you to
// hook into various parts of the response process.
func MakeDoneWriter(w http.ResponseWriter) DoneWriter {
	_, cn := w.(http.CloseNotifier)
	_, fl := w.(http.Flusher)
	_, hj := w.(http.Hijacker)
	_, rf := w.(io.ReaderFrom)

	bw := basicWriter{ResponseWriter: w}
	if cn && fl && hj && rf {
		return &fancyWriter{bw}
	}
	if fl {
		return &flushWriter{bw}
	}
	return &bw
}

type basicWriter struct {
	http.ResponseWriter
	done bool
}

func (w *basicWriter) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}

func (w *basicWriter) IsDone() bool {
	return w.done
}

/////////////////////////////////////////

// fancyWriter is a writer that additionally satisfies http.CloseNotifier,
// http.Flusher, http.Hijacker, and io.ReaderFrom.
type fancyWriter struct {
	basicWriter
}

// Deprecated: the CloseNotifier interface predates Go's context package.
// New code should use Request.Context instead.
func (f *fancyWriter) CloseNotify() <-chan bool {
	cn := f.basicWriter.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}
func (f *fancyWriter) Flush() {
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
func (f *fancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}
func (f *fancyWriter) ReadFrom(r io.Reader) (int64, error) {
	rf := f.basicWriter.ResponseWriter.(io.ReaderFrom)
	return rf.ReadFrom(r)
}

/////////////////////////////////////////

type flushWriter struct {
	basicWriter
}

func (f *flushWriter) Flush() {
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

/////////////////////////////////////////

// type checking to make sure that we have implemented all the interface functions correctly.
var (
	_ http.CloseNotifier = &fancyWriter{}
	_ http.Flusher       = &fancyWriter{}
	_ http.Hijacker      = &fancyWriter{}
	_ io.ReaderFrom      = &fancyWriter{}
	_ http.Flusher       = &flushWriter{}
)
