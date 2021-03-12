# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from odahuflow.sdk.models.base_model_ import Model
from odahuflow.sdk.models import util


class ConnectionReference(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, connection: str=None, path: str=None):  # noqa: E501
        """ConnectionReference - a model defined in Swagger

        :param connection: The connection of this ConnectionReference.  # noqa: E501
        :type connection: str
        :param path: The path of this ConnectionReference.  # noqa: E501
        :type path: str
        """
        self.swagger_types = {
            'connection': str,
            'path': str
        }

        self.attribute_map = {
            'connection': 'connection',
            'path': 'path'
        }

        self._connection = connection
        self._path = path

    @classmethod
    def from_dict(cls, dikt) -> 'ConnectionReference':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ConnectionReference of this ConnectionReference.  # noqa: E501
        :rtype: ConnectionReference
        """
        return util.deserialize_model(dikt, cls)

    @property
    def connection(self) -> str:
        """Gets the connection of this ConnectionReference.

        Next connection types are supported  # noqa: E501

        :return: The connection of this ConnectionReference.
        :rtype: str
        """
        return self._connection

    @connection.setter
    def connection(self, connection: str):
        """Sets the connection of this ConnectionReference.

        Next connection types are supported  # noqa: E501

        :param connection: The connection of this ConnectionReference.
        :type connection: str
        """

        self._connection = connection

    @property
    def path(self) -> str:
        """Gets the path of this ConnectionReference.

        User can override path otherwise Connection path will be used  # noqa: E501

        :return: The path of this ConnectionReference.
        :rtype: str
        """
        return self._path

    @path.setter
    def path(self, path: str):
        """Sets the path of this ConnectionReference.

        User can override path otherwise Connection path will be used  # noqa: E501

        :param path: The path of this ConnectionReference.
        :type path: str
        """

        self._path = path