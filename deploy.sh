
#!/bin/bash

set -e

STACK_NAME=$1
DEPLOYMENT_BUCKET=$2
TEST_STACK_NAME=$3

APIURL=`aws cloudformation describe-stacks \
            --stack-name $TEST_STACK_NAME \
            --query "Stacks[0].Outputs[0].{OutputValue:OutputValue}" \
            --output text`


make DEPLOYMENT_BUCKET=${DEPLOYMENT_BUCKET}
make package-aws DEPLOYMENT_BUCKET=${DEPLOYMENT_BUCKET}
make deploy STACK_NAME=${STACK_NAME} ENV=ci QueueName=SQSRMFCI RouteURL="${TEST_STACK_NAME}/test/api/hello"