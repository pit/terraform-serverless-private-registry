#!/usr/bin/env python3

import sys
import os
import requests
import argparse
import lib
import logging


parser = argparse.ArgumentParser(description='Test terraform-registry api.')
parser.add_argument('--domain', required=True, help='api base domain')
parser.add_argument('--username', required=True, help='api auth username')
parser.add_argument('--password', required=True, help='api auth password')
parser.add_argument('--debug', action='store_true', help='enable debug logging')
parser.add_argument('--test', default='providers,modules', help='what to test, comma separated')

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

if 'modules' in args.test:
    modules_test = lib.ModulesTest(domain=args.domain, username=args.username, password=args.password, fixtures_path=fixtures_path, debug=args.debug, logger=logger)
    modules_test.test()

if 'providers' in args.test:
    providers_test = lib.ProvidersTest(domain=args.domain, username=args.username, password=args.password, fixtures_path=fixtures_path, debug=args.debug, logger=logger)
    providers_test.test()
