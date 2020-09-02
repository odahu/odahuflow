//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package deployment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	md_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	"net/http"
	"reflect"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logMD = logf.Log.WithName("md-controller")

const (
	GetModelDeploymentURL    = "/model/deployment/:id"
	GetAllModelDeploymentURL = "/model/deployment"
	CreateModelDeploymentURL = "/model/deployment"
	UpdateModelDeploymentURL = "/model/deployment"
	DeleteModelDeploymentURL = "/model/deployment/:id"
	IDMdURLParam             = "id"
)

var (
	fieldsCache = map[string]int{}
)

func init() {
	elem := reflect.TypeOf(&md_repository.MdFilter{}).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tagName := elem.Field(i).Tag.Get(md_repository.TagKey)

		fieldsCache[tagName] = i
	}
}

type ModelDeploymentController struct {
	mdRepo      md_repository.Repository
	mdValidator *ModelDeploymentValidator
}

// @Summary Get a Model deployment
// @Description Get a Model deployment by id
// @Tags Deployment
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Model deployment id"
// @Success 200 {object} deployment.ModelDeployment
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/deployment/{id} [get]
func (mdc *ModelDeploymentController) getMD(c *gin.Context) {
	mdID := c.Param(IDMdURLParam)

	md, err := mdc.mdRepo.GetModelDeployment(mdID)
	if err != nil {
		logMD.Error(err, fmt.Sprintf("Retrieving %s model deployment", mdID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, md)
}

// @Summary Get list of Model deployments
// @Description Get list of Model deployments
// @Tags Deployment
// @Accept  json
// @Produce  json
// @Param size path int false "Number of entities in a response"
// @Param page path int false "Number of a page"
// @Success 200 {array} deployment.ModelDeployment
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/deployment [get]
func (mdc *ModelDeploymentController) getAllMDs(c *gin.Context) {
	f := &md_repository.MdFilter{}
	size, page, err := routes.URLParamsToFilter(c, f, fieldsCache)
	if err != nil {
		logMD.Error(err, "Malformed url parameters of model deployment request")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	mdList, err := mdc.mdRepo.GetModelDeploymentList(
		filter.ListFilter(f),
		filter.Size(size),
		filter.Page(page),
	)
	if err != nil {
		logMD.Error(err, "Retrieving list of model deployments")
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, mdList)
}

// @Summary Create a Model deployment
// @Description Create a Model  Results is created Model
// @Param md body deployment.ModelDeployment true "Create a Model deployment"
// @Tags Deployment
// @Accept  json
// @Produce  json
// @Success 201 {object} deployment.ModelDeployment
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/deployment [post]
func (mdc *ModelDeploymentController) createMD(c *gin.Context) {
	var md deployment.ModelDeployment

	if err := c.ShouldBindJSON(&md); err != nil {
		logMD.Error(err, "JSON binding of the model deployment is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mdc.mdValidator.ValidatesMDAndSetDefaults(&md); err != nil {
		logMD.Error(err, fmt.Sprintf("Validation of the model deployment is failed: %v", md))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mdc.mdRepo.CreateModelDeployment(&md); err != nil {
		logMD.Error(err, fmt.Sprintf("Creation of the model deployment: %+v", md))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusCreated, md)
}

// @Summary Update a Model deployment
// @Description Update a Model  Results is updated Model
// @Param md body deployment.ModelDeployment true "Update a Model deployment"
// @Tags Deployment
// @Accept  json
// @Produce  json
// @Success 200 {object} deployment.ModelDeployment
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/deployment [put]
func (mdc *ModelDeploymentController) updateMD(c *gin.Context) {
	var md deployment.ModelDeployment

	if err := c.ShouldBindJSON(&md); err != nil {
		logMD.Error(err, "JSON binding of the model deployment is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mdc.mdValidator.ValidatesMDAndSetDefaults(&md); err != nil {
		logMD.Error(err, fmt.Sprintf("Validation of the model deployment is failed: %v", md))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mdc.mdRepo.UpdateModelDeployment(&md); err != nil {
		logMD.Error(err, fmt.Sprintf("Update of the model deployment: %+v", md))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, md)
}

// @Summary Delete a Model deployment
// @Description Delete a Model deployment by id
// @Tags Deployment
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Model deployment id"
// @Success 200 {object} routes.HTTPResult
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/deployment/{id} [delete]
func (mdc *ModelDeploymentController) deleteMD(c *gin.Context) {
	mdID := c.Param(IDMdURLParam)

	if err := mdc.mdRepo.SetDeletionMark(mdID, true); err != nil {
		logMD.Error(err, fmt.Sprintf("Deletion of %s model deployment is failed", mdID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, routes.HTTPResult{Message: fmt.Sprintf("Model deployment %s was deleted", mdID)})
}