#!/bin/bash

export GIN_ANNOTATION_FILE="test.route.entry.go"

cd ./_example/complex
gin-annotation ./

result1=`diff ./route.entry.go ./test.route.entry.go`

if ["${result1}" == ""];then
  echo "_example/complex success"
else
  echo "_example/complex failure"
  exit 1
fi

cd ../simple
gin-annotation ./

result2=`diff ./route.entry.go ./test.route.entry.go`

if ["${result2}" = ""];then
  echo "_example/simple success"
else
  echo "_example/simple failure"
  exit 1
fi


