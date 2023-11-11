import inspect
import logging
import re
import os
import requests
import sys
import tempfile


class BaseTest():
    def __init__(self, **kwargs):
        self.domain = kwargs['domain']
        self.username = kwargs['username']
        self.password = kwargs['password']
        self.fixtures_path = kwargs['fixtures_path']
        self.is_debug = kwargs.get('debug', False)
        self.http = requests.Session()
        self.log = kwargs['logger']


    def get(self, url, **kwargs):
        url = 'https://' + re.sub(r'\/+', '/', url[len('https://'):])
        if kwargs.get('pass_auth', True):
            r = self.http.get(url, headers=kwargs.get('headers', {}), auth=(self.username, self.password))
        else:
            r = self.http.get(url, headers=kwargs.get('headers', {}))
        
        self.log.debug(f'Caller func: {inspect.stack()[1][3]}()')
        self.log.debug(f'---8<---------------------------------')
        self.log.debug(f'Request: {r.request.method}')
        self.log.debug(f'GET {url}')
        for hdr_name,hdr_val in r.request.headers.items():
            self.log.debug(f'{hdr_name}: {hdr_val}')
        self.log.debug(f'--------------------------------------')

        self.log.debug(f'Response: {r.status_code}')
        for hdr_name,hdr_val in r.headers.items():
            self.log.debug(f'{hdr_name}: {hdr_val}')
        self.log.debug(f'')            
        self.log.debug(f'{r.text}')
        self.log.debug(f'--------------------------------->8---')
        return r


    def put_binary(self, url, **kwargs):
        url = 'https://' + re.sub(r'\/+', '/', url[len('https://'):])
        if kwargs.get('pass_auth', True):
            r = self.http.put(url, headers=kwargs.get('headers', {'Content-Type': 'application/binary'}), data=kwargs['data'], auth=(self.username, self.password))
        else:
            r = self.http.put(url, headers=kwargs.get('headers', {'Content-Type': 'application/binary'}), data=kwargs['data'])
        
        self.log.debug(f'Caller func: {inspect.stack()[1][3]}()')
        self.log.debug(f'---8<---------------------------------')
        self.log.debug(f'Request: {r.request.method}')
        self.log.debug(f'PUT {url}')
        for hdr_name,hdr_val in r.request.headers.items():
            self.log.debug(f'{hdr_name}: {hdr_val}')
        self.log.debug(f'<body binary data>')
        self.log.debug(f'--------------------------------------')

        self.log.debug(f'Response: {r.status_code}')
        for hdr_name,hdr_val in r.headers.items():
            self.log.debug(f'{hdr_name}: {hdr_val}')
        self.log.debug(f'')            
        self.log.debug(f'{r.text}')
        self.log.debug(f'--------------------------------->8---')
        return r

    def post_json(self, url, **kwargs):
        url = 'https://' + re.sub(r'\/+', '/', url[len('https://'):])
        if kwargs.get('pass_auth', True):
            r = self.http.post(url, headers=kwargs.get('headers', {'Content-Type': 'application/json'}), json=kwargs['json'], auth=(self.username, self.password))
        else:
            r = self.http.post(url, headers=kwargs.get('headers', {'Content-Type': 'application/json'}), json=kwargs['json'])
        
        self.log.debug(f'Caller func: {inspect.stack()[1][3]}()')
        self.log.debug(f'---8<---------------------------------')
        self.log.debug(f'Request: {r.request.method}')
        self.log.debug(f'POST {url}')
        for hdr_name,hdr_val in r.request.headers.items():
            self.log.debug(f'{hdr_name}: {hdr_val}')
        self.log.debug(f'')
        self.log.debug(f'{r.request.body.decode("utf-8")}')
        self.log.debug(f'--------------------------------------')

        self.log.debug(f'Response: {r.status_code}')
        for hdr_name,hdr_val in r.headers.items():
            self.log.debug(f'{hdr_name}: {hdr_val}')
        self.log.debug(f'')            
        self.log.debug(f'{r.text}')
        self.log.debug(f'--------------------------------->8---')
        return r


    def download(self, url, filename):        
        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            with open(filename, 'wb') as f:
                for chunk in r.iter_content(chunk_size=8192): 
                    f.write(chunk)
