#
#    Copyright 2017 EPAM Systems
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#
"""
Robot test library - kubernetes dashboard
"""
import time

import kubernetes
import kubernetes.client
import kubernetes.config
import kubernetes.config.config_exception
from kubernetes.client import V1DeleteOptions, V1Pod, V1ObjectMeta, V1PodSpec, V1Toleration, V1Container, \
    V1ResourceRequirements
from kubernetes.client.rest import ApiException
from odahuflow.robot.utils import wait_until

FAT_POD_MEMORY = "4Gi"
FAT_POD_CPU = "4"
FAT_POD_IMAGE = 'alpine:3.9.3'
FAT_POD_NAME = "fat-pod-name"


def _build_client() -> kubernetes.client.ApiClient:
    kubernetes.config.load_kube_config()

    return kubernetes.client.ApiClient()


class K8s:
    """
    Kubernetes dashboard client for robot tests
    """

    ROBOT_LIBRARY_SCOPE = 'TEST SUITE'

    def __init__(self, namespace):
        """
        Init client
        :param namespace: default k8s namespace
        """
        self._context = None
        self._information = None
        self._odahuflow_group = 'odahuflow.odahu.org'
        self._model_training_version = 'v1alpha1'
        self._namespace = namespace
        self._model_training_plural = 'modeltrainings'
        self._model_training_info = self._odahuflow_group, self._model_training_version, \
                                    self._namespace, self._model_training_plural
        self._model_deployment_info = self._odahuflow_group, 'v1alpha1', \
                                      self._namespace, 'modeldeployments'
        self._default_vcs = 'odahuflow'

    def start_fat_pod(self, node_taint_key, node_taint_value):
        """
        Start fat pod
        """
        core_api = kubernetes.client.CoreV1Api(_build_client())

        pod = V1Pod(
            api_version='v1',
            kind='Pod',
            metadata=V1ObjectMeta(
                name=FAT_POD_NAME,
                annotations={
                    "sidecar.istio.io/inject": "false"
                },
            ),
            spec=V1PodSpec(
                restart_policy='Never',
                priority=0,
                tolerations=[
                    V1Toleration(
                        key=node_taint_key,
                        operator="Equal",
                        value=node_taint_value,
                        effect="NoSchedule",
                    )
                ],
                containers=[
                    V1Container(
                        name=FAT_POD_NAME,
                        image=FAT_POD_IMAGE,
                        resources=V1ResourceRequirements(
                            limits={'cpu': FAT_POD_CPU, 'memory': FAT_POD_MEMORY},
                            requests={'cpu': FAT_POD_CPU, 'memory': FAT_POD_MEMORY}
                        ),
                        command=["echo"],
                        args=["I am a fat :("]
                    )
                ]
            )
        )

        core_api.create_namespaced_pod(self._namespace, pod)

    def delete_fat_pod(self):
        """
        Delete fat pod
        """
        core_api = kubernetes.client.CoreV1Api(_build_client())

        try:
            core_api.delete_namespaced_pod(
                name=FAT_POD_NAME, namespace=self._namespace,
                body=V1DeleteOptions(propagation_policy='Foreground', grace_period_seconds=0)
            )
        except ApiException as e:
            if e.status != 404:
                raise e

    def wait_fat_pod_completion(self):
        """
        Wait completion of fat pod
        """
        core_api = kubernetes.client.CoreV1Api(_build_client())

        pod_completed_lambda = lambda: core_api.read_namespaced_pod(
            FAT_POD_NAME, self._namespace
        ).status.phase == "Succeeded"
        if not wait_until(pod_completed_lambda, iteration_duration=10, iterations=120):
            raise Exception("Timeout")

    def get_model_training_logs(self, name):
        """
        Get model training logs
        :param name: name of a model training resource
        :return: status
        """
        core_api = kubernetes.client.CoreV1Api(_build_client())

        return core_api.read_namespaced_pod_log(f'{name}-training-pod', self._namespace, container='builder')

    def check_all_containers_terminated(self, name):
        """
        Check that all pod containers are terminated
        :param name: name of a model training resource
        :return: None
        """
        core_api = kubernetes.client.CoreV1Api(_build_client())

        pod = core_api.read_namespaced_pod(f'{name}-training-pod', self._namespace)
        for container in pod.status.container_statuses:
            if not container.state.terminated:
                raise Exception(f'Container {container.name} of {pod.metadata.name} pod is still alive')

    def get_model_training_status(self, name):
        """
        Get model training status
        :param name: name of a model training resource
        :return: status
        """
        crds = kubernetes.client.CustomObjectsApi()

        model_training = crds.get_namespaced_custom_object(*self._model_training_info, name.lower())
        print(f'Fetched model training: {model_training}')

        status = model_training.get('status')
        return status if status else {}

    def get_model_deployment_status(self, name):
        """
        Get model training status
        :param name: name of a model training resource
        :return: status
        """
        crds = kubernetes.client.CustomObjectsApi()

        md = crds.get_namespaced_custom_object(*self._model_deployment_info, name.lower())
        print(f'Fetched model training: {md}')

        status = md.get('status')
        return status if status else {}

    def get_number_of_md_replicas(self, name, expected_number_of_replicas):
        """
        Get number of model deploymet replicas
        :param name: resource name
        :param expected_number_of_replicas: expected state
        :return:
        """
        md_status = self.get_model_deployment_status(name)
        return self.get_deployment_replicas(md_status["deployment"], self._namespace) == expected_number_of_replicas

    def wait_model_deployment_replicas(self, name, expected_number_of_replicas):
        """
        Wait specific status of a model training resource
        :param name: resource name
        :param expected_number_of_replicas: expected state
        """
        if not wait_until(lambda: self.get_number_of_md_replicas(name, expected_number_of_replicas),
                          iteration_duration=10, iterations=24):  # 4 min
            raise Exception("Timeout")

    def deployment_is_running(self, deployment_name, namespace):
        """
        Check that specific named deployment is okay (no one pending or failed pod)

        :param deployment_name: name of replication controller
        :type deployment_name: str
        :param namespace: name of namespace to look in
        :type namespace: str
        :raises: Exception
        :return: None
        """
        apps_api = kubernetes.client.AppsV1Api(_build_client())

        deployment = apps_api.read_namespaced_deployment(deployment_name, namespace)

        if deployment.status.replicas != deployment.status.ready_replicas:
            raise Exception(f"Deployment '{deployment_name}' is not ready: {deployment.status.ready_replicas}/"
                            f"{deployment.status.replicas} replicas are running")

    def get_model_deployment(self, deployment_name, namespace):
        """
        Get dict of deployment by model name and version

        :param namespace: k8s namespace
        :param deployment_name: k8s deployment name
        :return: k8s deployment object
        """
        extension_api = kubernetes.client.AppsV1Api(_build_client())
        print(namespace, self._namespace, deployment_name)
        return extension_api.read_namespaced_deployment(deployment_name, namespace or self._namespace)

    def get_deployment_replicas(self, deployment_name, namespace='default'):
        """
        Get number of replicas for a specified deployment from Kubernetes API

        :param deployment_name: name of a deployment
        :type deployment_name: str
        :param namespace: name of a namespace to look in
        :type namespace: str
        :return: number of replicas for a specified deployment
        :rtype int
        """
        apps_api = kubernetes.client.AppsV1Api(_build_client())
        scale_data = apps_api.read_namespaced_deployment(deployment_name, namespace)
        print(f"Scale data for {deployment_name} in {namespace} ns is {scale_data}")
        if scale_data is not None:
            return 0 if not scale_data.status.available_replicas else scale_data.status.available_replicas
        else:
            return None

    def wait_nodes_scale_down(self, node_taint_key, node_taint_value, timeout=600, sleep=60):
        """
        Wait finish of last job build

        :raises: Exception
        :return: None
        """
        core_api = kubernetes.client.CoreV1Api(_build_client())

        timeout = int(timeout)
        sleep = int(sleep)
        start = time.time()
        time.sleep(sleep)

        while True:
            nodes_num = 0

            for node in core_api.list_node().items:
                if not node.spec.taints:
                    continue

                for taint in node.spec.taints:
                    if taint.key == node_taint_key and taint.value == node_taint_value:
                        nodes_num += 1
                        break

            elapsed = time.time() - start

            if nodes_num == 0:
                print(f'Scaled node was successfully unscaled after {elapsed} seconds')
                return
            elif elapsed > timeout > 0:
                raise Exception(f'Node was not unscaled after {timeout} seconds wait')
            else:
                print(f'Current node count {nodes_num}. Sleep {sleep} seconds and try again')
                time.sleep(sleep)
