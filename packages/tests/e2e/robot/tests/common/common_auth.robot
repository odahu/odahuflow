*** Settings ***
Documentation       Check if all core components are secured
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Force Tags          common  security

*** Keywords ***
Url stay the same after log in
    [Arguments]  ${service_url}
    ${resp}=  Wait Until Keyword Succeeds  2m  5 sec  Wait Until Keyword Succeeds  2m  5 sec  Request as authorized user  ${service_url}  ${AUTH_TOKEN}
    should be equal  ${service_url}  ${resp.url}

Authorization should raise auth error if user is not authorized
    [Arguments]  ${service_url}
    ${resp}=  Wait Until Keyword Succeeds  2m  5 sec  Request as unauthorized user  ${service_url}
    Log              Response for ${service_url} is ${resp}
    Should contain   ${resp.text}    Log in to
    Should contain   ${resp.text}    Username or email
    Should contain   ${resp.text}    Password

*** Test Cases ***
Service url stay the same after log in
    [Tags]  apps
    [Documentation]  Service url stay the same after log in
    [Template]    Url stay the same after log in
    service_url=${EDGE_URL}/swagger/index.html
    service_url=${API_URL}/swagger/index.html
    service_url=${GRAFANA_URL}/?orgId=1&x=2
    service_url=${PROMETHEUS_URL}/graph?x=2&y=3
    service_url=${ALERTMANAGER_URL}/?orgId=1&x=2
    service_url=${JUPYTERLAB_URL}/lab
    service_url=${MLFLOW_URL}/
    service_url=${AIRFLOW_URL}/admin/

Invalid credentials raise Auth error
    [Tags]  apps  e2t
    [Documentation]  Invalid credentials raise Auth error
    [Template]    Authorization should raise auth error if user is not authorized
    service_url=${EDGE_URL}/swagger/index.html
    service_url=${API_URL}/swagger/index.html
    service_url=${GRAFANA_URL}/?orgId=1&x=2
    service_url=${PROMETHEUS_URL}/graph?x=2&y=3
    service_url=${ALERTMANAGER_URL}/?orgId=1&x=2
    service_url=${JUPYTERLAB_URL}/lab
    service_url=${MLFLOW_URL}/?a=1
    service_url=${AIRFLOW_URL}/admin/
