#!/usr/bin/python

import os
import sys
import yaml

if len(sys.argv) != 3:
    print("error in input parameters")

input_file = sys.argv[1]
output_file = sys.argv[2]

in_conf = None
out_conf = []

try:
    with open(input_file, 'r') as stream:
        try:
            conf = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print(e)
except FileNotFoundError as e:
    print(e)

i = 0

for elem in conf:
    try:
        rule = elem['rule']
        out_conf.append({'rule': rule, 'rule_id': i})
        i += 1
    except KeyError as e:
        continue

try:
    with open(output_file, 'w') as stream:
        try:
            stream.write(yaml.dump(out_conf))
        except yaml.YAMLError as e:
            print(e)
except FileNotFoundError as e:
    print(e)


