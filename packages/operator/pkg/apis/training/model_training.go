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

package training

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"time"
)

type ModelTraining struct {
	// Model training ID
	ID string `json:"id"`
	// Deletion mark
	DeletionMark bool `json:"deletionMark,omitempty" swaggerignore:"true"`
	// CreatedAt
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Model training specification
	Spec v1alpha1.ModelTrainingSpec `json:"spec,omitempty"`
	// Model training status
	Status v1alpha1.ModelTrainingStatus `json:"status,omitempty"`
}
