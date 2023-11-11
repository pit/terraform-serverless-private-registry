#!/usr/bin/env python3

import sys
import os
import requests
import argparse
import lib
import logging
import sh


parser = argparse.ArgumentParser(description='Test terraform-registry api.')
parser.add_argument('--bucket', required=True, help='bucket name')
parser.add_argument('--debug', action='store_true', help='enable debug logging')

args = parser.parse_args()
print(args)

handler = logging.StreamHandler(sys.stdout)
logger = logging.getLogger(__name__)
if args.debug:
    print('Enabling debug mode')

    formatter = logging.Formatter("%(asctime)s %(name)s [%(levelname)-7s] %(message)s")
    handler.setLevel(logging.DEBUG)
    handler.setFormatter(formatter)
    logger.setLevel(logging.DEBUG)
else:
    formatter = logging.Formatter("%(message)s")
    handler.setLevel(logging.WARNING)
    handler.setFormatter(formatter)
    logger.setLevel(logging.WARNING)
logger.addHandler(handler)


fixtures_path = f'{os.path.dirname(os.path.realpath(__file__))}/fixtures/'

sh.aws('s3', 'rm', f's3://{args.bucket}', '--recursive')
sh.aws('s3', 'cp', f'{fixtures_path}/module-test.tar.gz', f's3://{args.bucket}/modules/test-namespace/test-name/test-provider/2.0.0/test-namespace-test-name-test-provider-2.0.0.tar.gz', '--acl', 'bucket-owner-full-control')

sh.aws('s3', 'cp', f'{fixtures_path}/provider-darwin_amd64.zip', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_darwin_amd64.zip', '--metadata', 'X-Sha256=cf1ad324b97c32635eb35d11a1c9b2d2218a1437c260884a923548266cb1aea2', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-linux_amd64.zip', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_linux_amd64.zip', '--metadata', 'X-Sha256=48de8aa0baf80bbe4aa4e3f623a9a0231e9436afb62f25ec221fee682c013324', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-2.0.1-sha256sums', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_SHA256SUMS', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-2.0.1-sha256sums.sig', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_SHA256SUMS.sig', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-sha256sums.sig.pub', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_SHA256SUMS.sig.pub', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-sha256sums.sig.keyid', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_SHA256SUMS.sig.keyid', '--acl', 'bucket-owner-full-control')

sh.aws('s3', 'cp', f'{fixtures_path}/provider-darwin_amd64.zip', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.2/terraform-provider-test-type_2.0.2_darwin_amd64.zip', '--metadata', 'X-Sha256=cf1ad324b97c32635eb35d11a1c9b2d2218a1437c260884a923548266cb1aea2', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-linux_amd64.zip', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.2/terraform-provider-test-type_2.0.2_linux_amd64.zip', '--metadata', 'X-Sha256=48de8aa0baf80bbe4aa4e3f623a9a0231e9436afb62f25ec221fee682c013324', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-2.0.2-sha256sums', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.2_SHA256SUMS', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-2.0.2-sha256sums.sig', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.2_SHA256SUMS.sig', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-sha256sums.sig.pub', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.2/terraform-provider-test-type_2.0.2_SHA256SUMS.sig.pub', '--acl', 'bucket-owner-full-control')
sh.aws('s3', 'cp', f'{fixtures_path}/provider-sha256sums.sig.keyid', f's3://{args.bucket}/providers/test-namespace/test-type/2.0.2/terraform-provider-test-type_2.0.2_SHA256SUMS.sig.keyid', '--acl', 'bucket-owner-full-control')
