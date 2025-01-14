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
Robot test library - utils
"""
import os
import datetime
import socket
import time
import json
import urllib.parse
import requests

from odahuflow.robot.libraries.auth_client import get_authorization_headers

DOCKER_HOST_ENV = "DOCKER_HOST"

class Utils:
    """
    Utilities for robot tests
    """

    @staticmethod
    def request_as_authorized_user(service_url, token=None):
        """
        Request resource as authorized user

        :param service_url: target URL
        :param token: JWT bearer token
        :type service_url: str
        :return: final response - Response
        """
        return requests.get(service_url, headers=get_authorization_headers(token))

    @staticmethod
    def request_as_unauthorized_user(service_url):
        """
        Request resource as authorized user

        :param service_url: target URL
        :type service_url: str
        :return: final response - Response
        """
        return requests.get(service_url)

    @staticmethod
    def check_domain_exists(domain):
        """
        Check that domain (DNS name) has been registered

        :param domain: domain name (A record)
        :type domain: str
        :raises: Exception
        :return: None
        """
        try:
            return socket.gethostbyname(domain)
        except socket.gaierror as exception:
            if exception.errno == -2:
                raise Exception(f'Unknown domain name: {domain}') from exception
            else:
                raise

    @staticmethod
    def check_remote_file_exists(url, login=None, password=None):
        """
        Check that remote file exists (through HTTP/HTTPS)

        :param url: remote file URL
        :type url: str
        :param login: login
        :type login: str or None
        :param password: password
        :type password: str or None
        :raises: Exception
        :return: None
        """
        credentials = None
        if login and password:
            credentials = login, password

        response = requests.get(url,
                                stream=True,
                                verify=False,
                                auth=credentials)
        if response.status_code >= 400 or response.status_code < 200:
            raise Exception(f'Returned wrong status code: {response.status_code}')

        response.close()

    @staticmethod
    def sum_up(*values):
        """
        Sum up arguments

        :param values: Values to sum up
        :type values: int[]
        :return: Sum
        :rtype: int
        """
        result = 0
        for value in values:
            result += value
        return result

    @staticmethod
    def subtract(minuend, *values):
        """
        Subtract arguments from minuend

        :param minuend: A Minuend
        :type minuend: int
        :param values: Values to subtract from minuend
        :type values: int[]
        :rtype: int
        """
        result = minuend
        for value in values:
            result -= value
        return result

    @staticmethod
    def parse_api_inspect_columns_info(api_output):
        """
        Parse API inspect output

        :param api_output: API stdout
        :type api_output: str
        :return: list[list[str]] -- parsed API output
        """
        lines = api_output.splitlines()
        if len(lines) < 2:
            return []

        return [[item.strip() for item in line.split('|') if item] for line in lines[1:]]

    @staticmethod
    def find_model_information_in_api(parsed_api_output, model_name):
        """
        Get specific model API output

        :param parsed_api_output: parsed API output
        :type parsed_api_output: list[list[str]]
        :param model_name: model deployment name
        :type model_name: str
        :return: list[str] -- parsed API output for specific model
        """
        founded = [info for info in parsed_api_output if info[0] == model_name]
        if not founded:
            raise Exception(f'Info about model {model_name} not found')

        return founded[0]


    @staticmethod
    def check_valid_http_response(url, token=None):
        """
        Check if model return valid code for get request

        :param url: url with model_name for checking
        :type url: str
        :param token: token for the authorization
        :type token: str
        :return:  str -- response text
        """
        tries = 6
        error = None
        for i in range(tries):
            try:
                if token:
                    headers = {'Authorization': f'Bearer {token}'}
                    response = requests.get(url, timeout=10, headers=headers)
                else:
                    response = requests.get(url, timeout=10)

                if response.status_code == 200:
                    return response.text
                elif i >= 5:
                    raise Exception(f'Returned wrong status code: {response.status_code}')
                elif response.status_code >= 400 or response.status_code < 200:
                    print(f'Response code = {response.status_code}, sleep and try again')
                    time.sleep(3)
            except requests.exceptions.Timeout as e:
                error = e
                time.sleep(3)
        if error:
            raise error
        else:
            raise Exception('Unexpected case happen!')

    @staticmethod
    def execute_post_request_as_authorized_user(url, data=None, json_data=None):
        """
        Execute post request as authorized user

        :param url: url for request
        :type url: str
        :param data: data to send in request
        :type data: dict
        :param json_data: json data to send in request
        :type json_data: dict
        :return:  str -- response text
        """
        response = requests.post(url, json=json_data, data=data, headers=get_authorization_headers())

        return {"text": response.text, "code": response.status_code}

    @staticmethod
    def get_component_auth_page(url, token=None):
        """
        Get component main auth page
        :param url: component url
        :type url: str
        :param token: token for the authorization
        :type url: str
        :return:  response_code and response_text
        :rtype: dict
        """
        if token:
            headers = {'Authorization': f'Bearer {token}'}
            response = requests.get(url, timeout=10, headers=headers)
        else:
            response = requests.get(url, timeout=10)

        return {"response_code": response.status_code, "response_text": response.text}

    @staticmethod
    def parse_json_string(string):
        """
        Parse JSON string

        :param string: JSON string
        :type string: str
        :return: dict -- object
        """
        return json.loads(string)

    @staticmethod
    def get_current_time(time_template):
        """
        Get templated time

        :param time_template: time template
        :type time_template: str
        :return: None or str -- time from template
        """
        return datetime.datetime.utcnow().strftime(time_template)

    @staticmethod
    def get_future_time(offset, time_template):
        """
        Get templated time on `offset` seconds in future

        :param offset: time offset from current time in seconds
        :type offset: int
        :param time_template: time template
        :type time_template: str
        :return: str -- time from template
        """
        return (datetime.datetime.utcnow() +
                datetime.timedelta(seconds=offset)).strftime(time_template)

    @staticmethod
    def reformat_time(time_str, initial_format, target_format):
        """
        Convert date/time string from initial_format to target_format
        :param time_str: date/time string
        :type time_str: str
        :param initial_format: initial format of date/time string
        :type initial_format: str
        :param target_format: format to convert date/time object to
        :type target_format: str
        :return: str -- date/time string according to target_format
        """
        datetime_obj = datetime.datetime.strptime(time_str, initial_format)
        return datetime.datetime.strftime(datetime_obj, target_format)

    @staticmethod
    def get_timestamp_from_string(time_string, string_format):
        """
        Get timestamp from date/time string
        :param time_string: date/time string
        :type time_string: str
        :param string_format: format of time_string
        :type string_format: str
        :return: float -- timestamp
        """
        return datetime.datetime.strptime(time_string, string_format).timestamp()

    @staticmethod
    def wait_up_to_second(second, time_template=None):
        """
        Wait up to second then generate time from template

        :param second: target second (0..59)
        :type second: int
        :param time_template: (Optional) time template
        :type time_template: str
        :return: None or str -- time from template
        """
        current_second = datetime.datetime.now().second
        target_second = int(second)

        if current_second > target_second:
            sleep_time = 60 - (current_second - target_second)
        else:
            sleep_time = target_second - current_second

        if sleep_time:
            print(f'Waiting {sleep_time} second(s)')
            time.sleep(sleep_time)

        if time_template:
            return Utils.get_current_time(time_template)

    @staticmethod
    def order_list_of_dicts_by_key(list_of_dicts, field_key):
        """
        Order list of dictionaries by key as integer

        :param list_of_dicts: list of dictionaries
        :type list_of_dicts: List[dict]
        :param field_key: key to be ordered by
        :type field_key: str
        :return: List[dict] -- ordered list
        """
        return sorted(list_of_dicts, key=lambda item: int(item[field_key]))

    @staticmethod
    def concatinate_list_of_dicts_field(list_of_dicts, field_key):
        """
        Concatinate list of dicts field to string

        :param list_of_dicts: list of dictionaries
        :type list_of_dicts: List[dict]
        :param field_key: key of field to be concatinated
        :type field_key: str
        :return: str -- concatinated string
        """
        return ''.join([item[field_key] for item in list_of_dicts])

    @staticmethod
    def repeat_string_n_times(string, count):
        """
        Repeat string N times

        :param string: string to be repeated
        :type string: str
        :param count: count
        :type count: int
        :return: str -- result string
        """
        return string * int(count)

    @staticmethod
    def get_local_model_host():
        docker_host = os.getenv(DOCKER_HOST_ENV)
        if docker_host:
            return f"http://{urllib.parse.urlparse(docker_host).hostname}"
        else:
            return "http://0"
