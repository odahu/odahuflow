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

package packaging

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
)

func ConfigureRoutes(
	routeGroup *gin.RouterGroup,
	packService mp_repository.Service,
	packStorage mp_repository.Repository,
	piStorage mp_repository.PackagingIntegrationRepository,
	connStorage conn_repository.Repository,
	config config.ModelPackagingConfig,
	gpuResourceName string) {

	mtController := ModelPackagingController{
		packService: packService,
		packStorage: packStorage,
		validator: NewMpValidator(
			piStorage,
			connStorage,
			config.OutputConnectionID,
			gpuResourceName,
			config.DefaultResources,
		),
	}

	routeGroup.GET(GetModelPackagingURL, mtController.getMP)
	routeGroup.GET(GetAllModelPackagingURL, mtController.getAllMPs)
	routeGroup.POST(CreateModelPackagingURL, mtController.createMP)
	routeGroup.GET(GetModelPackagingLogsURL, mtController.getModelPackagingLog)
	routeGroup.PUT(UpdateModelPackagingURL, mtController.updateMP)
	routeGroup.PUT(SaveModelPackagingResultURL, mtController.saveMPResults)
	routeGroup.DELETE(DeleteModelPackagingURL, mtController.deleteMP)

}