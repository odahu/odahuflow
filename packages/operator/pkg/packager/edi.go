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

package packager

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
)

// The function extracts data from a repository and creates the packaging entity.
func (p *Packager) getPackaging() (*packaging.K8sPackager, error) {
	modelPackaging, err := p.packagingClient.GetModelPackaging(p.modelPackagingID)
	if err != nil {
		return nil, err
	}
	packagingIntegration, err := p.packagingClient.GetPackagingIntegration(modelPackaging.Spec.IntegrationName)
	if err != nil {
		return nil, err
	}

	targets := make([]packaging.PackagerTarget, 0, len(modelPackaging.Spec.Targets))
	for _, target := range modelPackaging.Spec.Targets {
		conn, err := p.connClient.GetConnection(target.ConnectionName)
		if err != nil {
			return nil, err
		}

		// Since connRepo here is actually an HTTP client, it returns connection with some fields base64-encoded
		if err := conn.DecodeBase64Fields(); err != nil {
			return nil, err
		}

		targets = append(targets, packaging.PackagerTarget{
			Name:       target.Name,
			Connection: *conn,
		})
	}

	modelHolder, err := p.connClient.GetConnection(modelPackaging.Spec.OutputConnection)
	if err != nil {
		return nil, err
	}

	// Since connRepo here is actually an HTTP client, it returns connection with some fields base64-encoded
	if err := modelHolder.DecodeBase64Fields(); err != nil {
		return nil, err
	}

	return &packaging.K8sPackager{
		ModelHolder:          modelHolder,
		ModelPackaging:       modelPackaging,
		PackagingIntegration: packagingIntegration,
		TrainingZipName:      modelPackaging.Spec.ArtifactName,
		Targets:              targets,
	}, nil
}
