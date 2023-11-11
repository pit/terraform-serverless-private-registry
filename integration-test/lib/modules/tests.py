from lib import BaseTest
from unittest import TestCase
from urllib.parse import urlparse
import boto3
import os
import tempfile


FIXTURE_NAMESPACE = 'test-namespace'
FIXTURE_NAME = 'test-name'
FIXTURE_PROVIDER = 'test-provider'
FIXTURE_VERSION = '1.2.3'


class ModulesTest(BaseTest):

    def __init__(self, **kwargs):
        super(
            self.__class__,
            self
        ).__init__(**kwargs)


    def test(self):
        print(f'modules tests - starting')

        object_methods = [method_name for method_name in dir(self) if callable(getattr(self, method_name)) and method_name.startswith('test_')]
        for method in sorted(object_methods):
            print(f'method: {method}')
            method_invocation = getattr(self, method)
            method_invocation()


    def test_discovery(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json', pass_auth=False)
        assert r.status_code == 401

        r = self.get(f'{self.domain}/.well-known/terraform.json')
        assert 'modules.v1' in r.json()
        assert r.json()['modules.v1'] != ''
        
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        assert 'custom.v1' in r.json()
        assert r.json()['custom.v1'] != ''
        
        print(f'test_discovery - passed')


    def test_upload(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        custom_api_prefix = r.json()['custom.v1']
        modules_api_prefix = r.json()['modules.v1']


        r = self.get(f'{self.domain}/{modules_api_prefix}/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/versions')
        assert 'modules' in r.json()
        assert 'versions' in r.json()['modules'][0]
        assert len(r.json()['modules'][0]['versions']) == 1
        assert 'version' in r.json()['modules'][0]['versions'][0]
        assert r.json()['modules'][0]['versions'][0]['version'] == '2.0.0'

        r = self.get(f'{self.domain}/{custom_api_prefix}/modules/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/1.2.1/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/module-test.tar.gz', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        r = self.get(f'{self.domain}/{custom_api_prefix}/modules/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/1.2.2/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/module-test.tar.gz', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        r = self.get(f'{self.domain}/{custom_api_prefix}/modules/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/1.2.3/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/module-test.tar.gz', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        r = self.get(f'{self.domain}/{custom_api_prefix}/modules/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/1.2.2/upload')
        assert r.status_code == 409
        assert r.json()['status'] == 'Already exists'
        assert r.json()['details'] == 'Module test-namespace/test-name/test-provider version 1.2.2 already exists'


        r = self.get(f'{self.domain}/{custom_api_prefix}/modules/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/2.0.0/upload')
        assert r.status_code == 409
        assert r.json()['status'] == 'Already exists'
        assert r.json()['details'] == 'Module test-namespace/test-name/test-provider version 2.0.0 already exists'

        r = self.get(f'{self.domain}/{modules_api_prefix}/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/versions')
        assert 'modules' in r.json()
        assert 'versions' in r.json()['modules'][0]
        assert len(r.json()['modules'][0]['versions']) == 4
        assert 'version' in r.json()['modules'][0]['versions'][0]
        assert r.json()['modules'][0]['versions'][0]['version'] == '1.2.1'
        assert r.json()['modules'][0]['versions'][1]['version'] == '1.2.2'
        assert r.json()['modules'][0]['versions'][2]['version'] == '1.2.3'
        assert r.json()['modules'][0]['versions'][3]['version'] == '2.0.0'

        
        print(f'test_upload - passed')


    def test_versions(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        modules_api_prefix = r.json()['modules.v1']

        r = self.get(f'{self.domain}/{modules_api_prefix}/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/versions')
        assert 'modules' in r.json()
        assert 'versions' in r.json()['modules'][0]
        assert len(r.json()['modules'][0]['versions']) == 4
        assert 'version' in r.json()['modules'][0]['versions'][0]
        assert r.json()['modules'][0]['versions'][0]['version'] == '1.2.1'
        assert r.json()['modules'][0]['versions'][1]['version'] == '1.2.2'
        assert r.json()['modules'][0]['versions'][2]['version'] == '1.2.3'
        assert r.json()['modules'][0]['versions'][3]['version'] == '2.0.0'


    def test_download(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        modules_api_prefix = r.json()['modules.v1']

        r = self.get(f'{self.domain}/{modules_api_prefix}/{FIXTURE_NAMESPACE}/{FIXTURE_NAME}/{FIXTURE_PROVIDER}/2.0.0/download')
        assert r.status_code == 204
        assert r.text == ''
        assert r.headers['X-Terraform-Get'] != ''
        assert urlparse(r.headers['X-Terraform-Get']) != None

        tmp = tempfile.NamedTemporaryFile(delete=False)
        try:
            self.download(r.headers['X-Terraform-Get'], tmp.name)
        finally:
            tmp.close()
            os.unlink(tmp.name)
