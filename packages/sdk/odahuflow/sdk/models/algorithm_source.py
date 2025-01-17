# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models.object_storage import ObjectStorage  # noqa: F401,E501
from odahuflow.sdk.models.vcs import VCS  # noqa: F401,E501
from odahuflow.sdk.models import util


class AlgorithmSource(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, object_storage: ObjectStorage=None, vcs: VCS=None):  # noqa: E501
        """AlgorithmSource - a model defined in Swagger

        :param object_storage: The object_storage of this AlgorithmSource.  # noqa: E501
        :type object_storage: ObjectStorage
        :param vcs: The vcs of this AlgorithmSource.  # noqa: E501
        :type vcs: VCS
        """
        self.swagger_types = {
            'object_storage': ObjectStorage,
            'vcs': VCS
        }

        self.attribute_map = {
            'object_storage': 'objectStorage',
            'vcs': 'vcs'
        }

        self._object_storage = object_storage
        self._vcs = vcs

    @classmethod
    def from_dict(cls, dikt) -> 'AlgorithmSource':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The AlgorithmSource of this AlgorithmSource.  # noqa: E501
        :rtype: AlgorithmSource
        """
        return util.deserialize_model(dikt, cls)

    @property
    def object_storage(self) -> ObjectStorage:
        """Gets the object_storage of this AlgorithmSource.


        :return: The object_storage of this AlgorithmSource.
        :rtype: ObjectStorage
        """
        return self._object_storage

    @object_storage.setter
    def object_storage(self, object_storage: ObjectStorage):
        """Sets the object_storage of this AlgorithmSource.


        :param object_storage: The object_storage of this AlgorithmSource.
        :type object_storage: ObjectStorage
        """

        self._object_storage = object_storage

    @property
    def vcs(self) -> VCS:
        """Gets the vcs of this AlgorithmSource.


        :return: The vcs of this AlgorithmSource.
        :rtype: VCS
        """
        return self._vcs

    @vcs.setter
    def vcs(self, vcs: VCS):
        """Sets the vcs of this AlgorithmSource.


        :param vcs: The vcs of this AlgorithmSource.
        :type vcs: VCS
        """

        self._vcs = vcs
