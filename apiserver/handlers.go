/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package apiserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *ApiServer) initRouter() {
	s.root.Add(http.MethodGet, "/status", s.Status)
	s.root.Add(http.MethodGet, "/", s.Welcome)
}

func (s *ApiServer) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
	return nil
}
func (s *ApiServer) Welcome(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Welcome to PeaEdge",
	})
	return nil
}
