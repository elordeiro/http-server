package main

import (
	"os"
	"strconv"
	"strings"
)

func (server *Server) Handle(req *Request) *Response {
	switch req.Method {
	case MethodGet:
		return server.Get(req)
	default:
		return req.NotFound(server)
	}
}

func (server *Server) Get(req *Request) *Response {
	if req.Path == "/" {
		if handler, ok := server.Paths["/"]; ok {
			return handler(req, server)
		} else {
			return req.NotFound(server)
		}
	}

	parts := strings.Split(req.Path, "/")
	if handler, ok := server.Paths["/"+parts[1]]; ok {
		return handler(req, server)
	}
	return req.NotFound(server)
}

func (req *Request) Ok(s *Server) *Response {
	return &Response{
		Version: req.Version,
		Status:  StatusOk,
		Reason:  "OK",
		Headers: map[string]string{},
		Body:    "",
	}
}

func (req *Request) Echo(s *Server) *Response {
	headers := map[string]string{}
	headers["Content-Type"] = "text/plain"
	body := strings.TrimPrefix(req.Path, "/echo/")
	headers["Content-Length"] = strconv.Itoa(len(body))
	return &Response{
		Version: req.Version,
		Status:  StatusOk,
		Reason:  "OK",
		Headers: headers,
		Body:    body,
	}
}

func (req *Request) UserAgent(s *Server) *Response {
	headers := map[string]string{}
	headers["Content-Type"] = "text/plain"
	body := req.Headers["user-agent"]
	headers["Content-Length"] = strconv.Itoa(len(body))
	return &Response{
		Version: req.Version,
		Status:  StatusOk,
		Reason:  "OK",
		Headers: headers,
		Body:    body,
	}
}

func (req *Request) Files(s *Server) *Response {
	headers := map[string]string{}
	headers["Content-Type"] = "application/octet-stream"
	path := s.Directory + strings.TrimPrefix(req.Path, "/files/")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return req.NotFound(s)
	}

	file, err := os.Open(path)
	if err != nil {
		return req.NotFound(s)
	}

	content := make([]byte, 1024)
	n, err := file.Read(content)
	if err != nil {
		return req.NotFound(s)
	}

	headers["Content-Length"] = strconv.Itoa(n)
	return &Response{
		Version: req.Version,
		Status:  StatusOk,
		Reason:  "OK",
		Headers: headers,
		Body:    string(content),
	}
}

func (req *Request) NotFound(s *Server) *Response {
	return &Response{
		Version: req.Version,
		Status:  StatusNotFound,
		Reason:  "Not Found",
		Headers: map[string]string{},
	}
}
