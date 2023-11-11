from lib import BaseTest
from unittest import TestCase
from urllib.parse import urlparse
import boto3
import os
import tempfile
import base64


FIXTURE_NAMESPACE = 'test-namespace'
FIXTURE_TYPE = 'test-type'
FIXTURE_VERSION = '1.2.3'


class ProvidersTest(BaseTest):

    def __init__(self, **kwargs):
        super(
            self.__class__,
            self
        ).__init__(**kwargs)


    def test(self):
        print(f'providers tests - starting')

        object_methods = [method_name for method_name in dir(self) if callable(getattr(self, method_name)) and method_name.startswith('test_')]
        for method in sorted(object_methods):
            print(f'method: {method}')
            method_invocation = getattr(self, method)
            method_invocation()


    def test_discovery(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json', pass_auth=False)
        assert r.status_code == 401

        r = self.get(f'{self.domain}/.well-known/terraform.json')
        assert 'providers.v1' in r.json()
        assert r.json()['providers.v1'] != ''
        
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        assert 'custom.v1' in r.json()
        assert r.json()['custom.v1'] != ''
        
        print(f'test_discovery - passed')


    def test_upload(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        custom_api_prefix = r.json()['custom.v1']

        r = self.get(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.1/linux/amd64/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/provider-linux_amd64.zip', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        r = self.get(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.1/darwin/amd64/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/provider-darwin_amd64.zip', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        with open(f'{self.fixtures_path}/provider-sha256sums.sig.keyid', 'r') as f:
            data_keyid = f.read()

        with open(f'{self.fixtures_path}/provider-1.2.1-sha256sums', 'w') as f:
            f.write("cf1ad324b97c32635eb35d11a1c9b2d2218a1437c260884a923548266cb1aea2  terraform-provider-test-type_1.2.1_darwin_amd64.zip\n")
            f.write("48de8aa0baf80bbe4aa4e3f623a9a0231e9436afb62f25ec221fee682c013324  terraform-provider-test-type_1.2.1_linux_amd64.zip\n")
        with open(f'{self.fixtures_path}/provider-1.2.1-sha256sums', 'r') as f:
            data_sha256sums = f.read()

        with open(f'{self.fixtures_path}/provider-1.2.1-sha256sums.sig', 'rb') as f:
            data_sha256sums_sig = f.read()

        with open(f'{self.fixtures_path}/provider-sha256sums.sig.pub', 'r') as f:
            data_sha256sums_sig_pub = f.read()
        data = {
            'keyId': data_keyid.strip(),
            'sha256Sums': data_sha256sums,
            'sha256SumsSig': base64.b64encode(data_sha256sums_sig).decode("utf-8"),
            'sha256SumsSigPub': data_sha256sums_sig_pub,
        }
        r = self.post_json(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.1/checksums/upload', json=data)
        assert r.status_code == 200
        assert r.json()['status'] == 'ok'


    def test_upload_409(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        custom_api_prefix = r.json()['custom.v1']

        r = self.get(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.2/linux/amd64/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/provider-linux_amd64.zip', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        r = self.get(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.2/darwin/amd64/upload')
        assert 'url' in r.json()
        assert r.json()['url'] != ''

        upload_url = r.json()['url']

        with open(f'{self.fixtures_path}/provider-darwin_amd64.zip', 'rb') as f:
            r = self.put_binary(upload_url, data=f, pass_auth=False)

        assert r.status_code == 200


        with open(f'{self.fixtures_path}/provider-sha256sums.sig.keyid', 'r') as f:
            data_keyid = f.read()

        with open(f'{self.fixtures_path}/provider-1.2.2-sha256sums', 'w') as f:
            f.write("cf1ad324b97c32635eb35d11a1c9b2d2218a1437c260884a923548266cb1aea2  terraform-provider-test-type_1.2.2_darwin_amd64.zip\n")
            f.write("48de8aa0baf80bbe4aa4e3f623a9a0231e9436afb62f25ec221fee682c013324  terraform-provider-test-type_1.2.2_linux_amd64.zip\n")
        with open(f'{self.fixtures_path}/provider-1.2.2-sha256sums', 'r') as f:
            data_sha256sums = f.read()

        with open(f'{self.fixtures_path}/provider-1.2.2-sha256sums.sig', 'rb') as f:
            data_sha256sums_sig = f.read()

        with open(f'{self.fixtures_path}/provider-sha256sums.sig.pub', 'r') as f:
            data_sha256sums_sig_pub = f.read()
        data = {
            'keyId': data_keyid.strip(),
            'sha256Sums': data_sha256sums,
            'sha256SumsSig': base64.b64encode(data_sha256sums_sig).decode("utf-8"),
            'sha256SumsSigPub': data_sha256sums_sig_pub,
        }
        r = self.post_json(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.2/checksums/upload', json=data)
        assert r.status_code == 200
        assert r.json()['status'] == 'ok'


        r = self.get(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/1.2.2/linux/amd64/upload')
        assert r.status_code == 409


    def test_upload_404(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        custom_api_prefix = r.json()['custom.v1']

        with open(f'{self.fixtures_path}/provider-sha256sums.sig.keyid', 'r') as f:
            data_keyid = f.read()
        with open(f'{self.fixtures_path}/provider-sha256sums', 'r') as f:
            data_sha256sums = f.read()
        with open(f'{self.fixtures_path}/provider-sha256sums.sig', 'rb') as f:
            data_sha256sums_sig = f.read()
        with open(f'{self.fixtures_path}/provider-sha256sums.sig.pub', 'r') as f:
            data_sha256sums_sig_pub = f.read()
        data = {
            'keyId': data_keyid.strip(),
            'sha256Sums': data_sha256sums,
            'sha256SumsSig': base64.b64encode(data_sha256sums_sig).decode("utf-8"),
            'sha256SumsSigPub': data_sha256sums_sig_pub,
        }
        r = self.post_json(f'{self.domain}/{custom_api_prefix}/providers/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/4.0.4/checksums/upload', json=data)
        assert r.status_code == 404
        assert r.json()['status'] == 'Not found'

        print(f'test_upload - passed')


    def test_versions(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        providers_api_prefix = r.json()['providers.v1']

        r = self.get(f'{self.domain}/{providers_api_prefix}/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/versions')
        assert 'versions' in r.json()
        assert len(r.json()['versions']) == 4
        assert 'version' in r.json()['versions'][0]
        assert r.json()['versions'][0]['version'] == '1.2.1'
        assert len(r.json()['versions'][0]['platforms']) == 2

        assert r.json()['versions'][1]['version'] == '1.2.2'
        assert r.json()['versions'][2]['version'] == '2.0.1'
        assert r.json()['versions'][3]['version'] == '2.0.2'


    def test_download(self):
        r = self.get(f'{self.domain}/.well-known/terraform.json')
        providers_api_prefix = r.json()['providers.v1']

        r = self.get(f'{self.domain}/{providers_api_prefix}/{FIXTURE_NAMESPACE}/{FIXTURE_TYPE}/2.0.1/download/linux/amd64')
        assert r.status_code == 200
        
        # # ...
        # assert urlparse(r.headers['X-Terraform-Get']) != None

        # tmp = tempfile.NamedTemporaryFile(delete=False)
        # try:
        #     self.download(r.headers['X-Terraform-Get'], tmp.name)
        # finally:
        #     tmp.close()
        #     os.unlink(tmp.name)
