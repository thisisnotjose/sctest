#!/bin/bash -ex

this_dir=$(dirname $0)
export totalEvents=100000
export concurrencyLevel=100
time java -server -Xmx1G -jar "${this_dir}/follower-maze-2.0.jar"
