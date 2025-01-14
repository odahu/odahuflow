/*
 * Copyright 2019 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package routes_test

import (
	"fmt"
	"net/http"
	"testing"

	odahuflow_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
)

type UtilsSuite struct {
	suite.Suite
	g *GomegaWithT
}

func (s *UtilsSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}

// TODO: Remove the test after implementing custom exceptions for ALL entity's repositories
func (s *UtilsSuite) TestKubernetesErrors() {
	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		errors.NewNotFound(v1beta2.Resource("statefulset"), "test"),
	)).Should(Equal(http.StatusNotFound))

	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		errors.NewAlreadyExists(v1beta2.Resource("statefulset"), "test"),
	)).Should(Equal(409))
}

func (s *UtilsSuite) TestOdahuflowErrors() {
	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		odahuflow_errors.NotFoundError{},
	)).Should(Equal(http.StatusNotFound))

	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		odahuflow_errors.AlreadyExistError{},
	)).Should(Equal(http.StatusConflict))

	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		odahuflow_errors.SerializationError{},
	)).Should(Equal(http.StatusInternalServerError))

	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		odahuflow_errors.ForbiddenError{},
	)).Should(Equal(http.StatusForbidden))
}

func (s *UtilsSuite) TestUnknownError() {
	s.g.Expect(odahuflow_errors.CalculateHTTPStatusCode(
		fmt.Errorf("some exception"),
	)).Should(Equal(http.StatusInternalServerError))
}
