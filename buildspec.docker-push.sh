#!/bin/sh

if [ -z "$CODEBUILD_WEBHOOK_HEAD_REF" ]; then
  export IMAGE_TAG=`echo $CODEBUILD_WEBHOOK_HEAD_REF | rev | cut -d'/' -f1 | rev`
fi
if [ -z "$IMAGE_TAG" ]; then
  export IMAGE_TAG="latest"
fi
cat ~/.docker/config.json
echo docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME:$IMAGE_TAG
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME:$IMAGE_TAG
