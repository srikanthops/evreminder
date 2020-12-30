#!/bin/bash

last_tag=$(git tag --sort=-creatordate | head -n 1)
dist_from_last_tag=$(git rev-list "${last_tag}"..HEAD --count)
if [[ $dist_from_last_tag == 0 ]]; then
    # tagged build
    echo "${last_tag}"
elif [[ ${CIRCLE_SHA1} != "" ]]; then
  echo "${last_tag}-${dist_from_last_tag}-${CIRCLE_SHA1}"
else
    echo "${last_tag}-${dist_from_last_tag}-dev"
fi

