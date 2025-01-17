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

package packaging_test

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging_integration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"go.uber.org/multierr"
	"testing"

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	mp_post_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
)

var (
	piIDMpValid                      = "pi-id"
	piEntrypointMpValid              = "/usr/bin/test"
	piImageMpValid                   = "test:image"
	piArtifactNameMpValid            = "some-artifact-name.zip"
	connDockerTypeMpValid            = "docker-conn"
	connS3TypeMpValid                = "s3-conn"
	defaultTargetArgument1Connection = "default-conn-id"
	defaultResources                 = config.NewDefaultModelPackagingConfig().DefaultResources
	piArgumentsMpValid               = packaging.JsonSchema{
		Properties: []packaging.Property{
			{
				Name: "argument-1",
				Parameters: []packaging.Parameter{
					{
						Name:  "minimum",
						Value: float64(5),
					},
					{
						Name:  "type",
						Value: "number",
					},
				},
			},
			{
				Name: "argument-2",
				Parameters: []packaging.Parameter{
					{
						Name:  "type",
						Value: "string",
					},
				},
			},
		},
		Required: []string{"argument-1"},
	}
	piTargetsMpValid = []v1alpha1.TargetSchema{
		{
			Name: "target-1",
			ConnectionTypes: []string{
				string(connection.S3Type),
				string(connection.GcsType),
				string(connection.AzureBlobType),
			},
			Default:  defaultTargetArgument1Connection,
			Required: false,
		},
		{
			Name: "target-2",
			ConnectionTypes: []string{
				string(connection.DockerType),
			},
			Required: true,
		},
	}
	validNodeSelector = map[string]string{"mode": "valid"}
	validPackaging    = packaging.ModelPackaging{
		ID: "valid-id",
		Spec: packaging.ModelPackagingSpec{
			IntegrationName:  piIDMpValid,
			ArtifactName:     "test",
			OutputConnection: connS3TypeMpValid,
			Arguments: map[string]interface{}{
				"argument-1": 5,
			},
			Targets: []v1alpha1.Target{
				{
					Name:           "target-2",
					ConnectionName: connDockerTypeMpValid,
				},
			},
			NodeSelector: validNodeSelector,
		},
	}
)

type ModelPackagingValidationSuite struct {
	suite.Suite
	g            *GomegaWithT
	mpKubeClient kube_client.Client
	mpRepo       mp_repository.Repository
	piService    packagingIntegrationService
	connRepo     conn_repository.Repository
	validator    *pack_route.MpValidator
}

func (s *ModelPackagingValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *ModelPackagingValidationSuite) SetupSuite() {

	s.mpKubeClient = kube_client.NewClient(testNamespace, testNamespace, kubeClient, cfg)

	s.mpRepo = mp_post_repository.PackagingRepo{DB: db}
	piRepo := mp_post_repository.PackagingIntegrationRepository{DB: db}
	s.piService = packaging_integration.NewService(&piRepo)

	s.connRepo = conn_k8s_repository.NewRepository(testNamespace, kubeClient)

	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.NodePools = append(packagingConfig.NodePools, config.NodePool{NodeSelector: validNodeSelector})

	s.validator = pack_route.NewMpValidator(
		s.piService,
		s.connRepo,
		packagingConfig,
		config.NvidiaResourceName,
	)

	err := s.piService.CreatePackagingIntegration(&packaging.PackagingIntegration{
		ID: piIDMpValid,
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   piEntrypointMpValid,
			DefaultImage: piImageMpValid,
			Schema: packaging.Schema{
				Targets:   piTargetsMpValid,
				Arguments: piArgumentsMpValid,
			},
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}

	err = s.connRepo.SaveConnection(&connection.Connection{
		ID: connDockerTypeMpValid,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.DockerType,
		},
		Status: v1alpha1.ConnectionStatus{},
	})
	if err != nil {
		s.T().Fatal(err)
	}

	err = s.connRepo.SaveConnection(&connection.Connection{
		ID: connS3TypeMpValid,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.S3Type,
		},
		Status: v1alpha1.ConnectionStatus{},
	})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelPackagingValidationSuite) TearDownSuite() {
	if err := s.connRepo.DeleteConnection(connS3TypeMpValid); err != nil {
		s.T().Fatal(err)
	}

	if err := s.connRepo.DeleteConnection(connDockerTypeMpValid); err != nil {
		s.T().Fatal(err)
	}

	if err := s.piService.DeletePackagingIntegration(piIDMpValid); err != nil {
		s.T().Fatal(err)
	}
}

func TestModelPackagingValidationSuite(t *testing.T) {
	suite.Run(t, new(ModelPackagingValidationSuite))
}

func (s *ModelPackagingValidationSuite) TestMpIDExplicitly() {
	id := "some-id"
	mp := &packaging.ModelPackaging{
		ID:   id,
		Spec: packaging.ModelPackagingSpec{},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.ID).Should(Equal(id))
}

func (s *ModelPackagingValidationSuite) TestMpImage() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
		},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.Image).ShouldNot(BeEmpty())
	s.g.Expect(mp.Spec.Image).Should(Equal(piImageMpValid))
}

func (s *ModelPackagingValidationSuite) TestMpImageExplicitly() {
	image := "some-image"
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Image:           image,
		},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.Image).Should(Equal(image))
}

func (s *ModelPackagingValidationSuite) TestMpArtifactName() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			ArtifactName: piArtifactNameMpValid,
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).ShouldNot(ContainSubstring(pack_route.TrainingIDOrArtifactNameErrorMessage))
}

func (s *ModelPackagingValidationSuite) TestMpArtifactNameMissed() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.TrainingIDOrArtifactNameErrorMessage))
}

func (s *ModelPackagingValidationSuite) TestMpIntegrationNameEmpty() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyIntegrationNameErrorMessage))
	// "not found" substring here is expected in case when packaging integration can't be found by
	// .Spec.IntegrationName, but because .Spec.IntegrationName is empty we shouldn't try to find
	// it at all
	s.g.Expect(err.Error()).Should(Not(ContainSubstring("not found")))
}

func (s *ModelPackagingValidationSuite) TestMpIntegrationNotFound() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: "some-packaging-name",
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(
		"packaging integration with name .spec.integrationName = \"some-packaging-name\" is not found"))
}

func (s *ModelPackagingValidationSuite) TestMpNotValidArgumentsSchema() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Arguments: map[string]interface{}{
				"argument-1": 4,
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("argument-1: Must be greater than or equal to 5"))
}

func (s *ModelPackagingValidationSuite) TestMpAdditionalArguments() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Arguments: map[string]interface{}{
				"argument-3": "value",
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("Additional property argument-3 is not allowed"))
}

func (s *ModelPackagingValidationSuite) TestMpRequiredArguments() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Arguments: map[string]interface{}{
				"argument-2": "value",
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("argument-1 is required"))
}

func (s *ModelPackagingValidationSuite) TestMpRequiredTargets() {
	ti := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           "target-1",
					ConnectionName: connS3TypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(ti)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("[target-2] are required targets"))
}

func (s *ModelPackagingValidationSuite) TestMpDefaultTargets() {
	ti := &packaging.ModelPackaging{
		ID: "valid-id",
		Spec: packaging.ModelPackagingSpec{
			IntegrationName:  piIDMpValid,
			ArtifactName:     "test",
			OutputConnection: connS3TypeMpValid,
			Arguments: map[string]interface{}{
				"argument-1": 5,
			},
			Targets: []v1alpha1.Target{
				{
					Name:           "target-2",
					ConnectionName: connDockerTypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(ti)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(ti.Spec.Targets).Should(HaveLen(2))
	s.g.Expect(ti.Spec.Targets[0].Name).Should(Equal("target-2"))
	s.g.Expect(ti.Spec.Targets[0].ConnectionName).Should(Equal(connDockerTypeMpValid))
	s.g.Expect(ti.Spec.Targets[1].Name).Should(Equal("target-1"))
	s.g.Expect(ti.Spec.Targets[1].ConnectionName).Should(Equal(defaultTargetArgument1Connection))
}

func (s *ModelPackagingValidationSuite) TestMpNotFoundTargets() {
	targetNotFoundName := "target-not-found"
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           targetNotFoundName,
					ConnectionName: connS3TypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		pack_route.TargetNotFoundErrorMessage, targetNotFoundName, piIDMpValid,
	)))
}

func (s *ModelPackagingValidationSuite) TestMpTargetConnNotFound() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           "target-1",
					ConnectionName: "conn-not-found",
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("not found"))
}

func (s *ModelPackagingValidationSuite) TestMpTargetConnWrongType() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           "target-1",
					ConnectionName: connDockerTypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		pack_route.NotValidConnTypeErrorMessage, "target-1", connection.DockerType, piIDMpValid,
	)))
}

func (s *ModelPackagingValidationSuite) TestMpGenerateDefaultResources() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
		},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.Resources).ShouldNot(BeNil())
	s.g.Expect(mp.Spec.Resources).Should(Equal(&defaultResources))
}

func (s *ModelPackagingValidationSuite) TestMpResourcesValidation() {
	wrongResourceValue := "wrong res"
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Resources: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					Memory: &wrongResourceValue,
					GPU:    &wrongResourceValue,
					CPU:    &wrongResourceValue,
				},
				Requests: &v1alpha1.ResourceList{
					Memory: &wrongResourceValue,
					GPU:    &wrongResourceValue,
					CPU:    &wrongResourceValue,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())

	errorMessage := err.Error()
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of memory request is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of cpu request is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of memory limit is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of cpu limit is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of gpu limit is failed: quantities must match the regular expression"))
}

func (s *ModelPackagingValidationSuite) TestOutputConnection() {
	testMpOutConnDefault := testOutConnDefault
	testMpOutConn := testOutConn
	testMpOutConnNotFound := testOutConnNotFound

	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	// If configuration output connection is not set then user must specify it as ModelTraining parameter
	err := pack_route.NewMpValidator(
		s.piService,
		s.connRepo,
		config.NewDefaultModelPackagingConfig(),
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(validation.EmptyValueStringError, "OutputConnection")))

	// If configuration output connection is set and user has not passed output connection as training
	// parameter then output connection value from configuration will be used as default
	packConfig := config.NewDefaultModelPackagingConfig()
	packConfig.OutputConnectionID = testMpOutConnDefault
	_ = pack_route.NewMpValidator(
		s.piService,
		s.connRepo,
		packConfig,
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.OutputConnection).Should(Equal(testMpOutConnDefault))

	// If configuration output connection is set but user also has passed output connection as training
	// parameter then user value is used
	packConfig = config.NewDefaultModelPackagingConfig()
	packConfig.OutputConnectionID = "default-output-connection"
	mp.Spec.OutputConnection = testMpOutConn
	_ = pack_route.NewMpValidator(
		s.piService,
		s.connRepo,
		config.NewDefaultModelPackagingConfig(),
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.OutputConnection).Should(Equal(testMpOutConn))

	// If connection kubePackClient doesn't contain connection with passed ID validation must raise NotFoundError
	mp = &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{OutputConnection: testMpOutConnNotFound},
	}
	err = pack_route.NewMpValidator(
		s.piService,
		s.connRepo,
		config.NewDefaultModelPackagingConfig(),
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("entity %q is not found", testMpOutConnNotFound))
}

func (s *ModelPackagingValidationSuite) TestValidateID() {
	mp := &packaging.ModelPackaging{
		ID: "not-VALID-id-",
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

// Tests that nil node selector is considered valid
func (s *ModelPackagingValidationSuite) TestValidateNodeSelector_nil() {
	mp := validPackaging
	mp.Spec.NodeSelector = nil
	err := s.validator.ValidateAndSetDefaults(&mp)
	s.Assertions.Nil(err)
}

// Packaging object has valid node selector that exists in config
func (s *ModelPackagingValidationSuite) TestValidateNodeSelector_Valid() {
	mp := validPackaging
	err := s.validator.ValidateAndSetDefaults(&mp)
	s.Assertions.Nil(err)
}

// Packaging object has invalid node selector that does not exist in config
// Expect validator to return exactly one error
func (s *ModelPackagingValidationSuite) TestValidateNodeSelector_Invalid() {
	mp := validPackaging
	mp.Spec.NodeSelector = map[string]string{"mode": "invalid"}
	err := s.validator.ValidateAndSetDefaults(&mp)
	s.Assertions.NotNil(err)
	s.Assertions.Len(multierr.Errors(err), 1)
}
