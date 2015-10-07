#!/usr/bin/env bash

# Keep track of failures
FAILED=0

# Output a fail message
fail() {
  echo ''
  echo 'FAIL'
  echo $1
  let "FAILED += 1"
}

# Show progress when tests pass
pass() {
  echo -n '.'
}

# Output a message and return with no error when all tests has passed
finished() {
  echo ''
  if [[ $FAILED -ne 0 ]]; then
    echo 'There were failed integration tests'
  else
    echo 'All integration tests passed'
  fi
  exit $FAILED
}

# Bin to use in tests
BIN=./build/wiki

# Test that error message and usage is printed if no query
$BIN > /dev/null 2>&1
STATUS=$?
if [[ $STATUS -ne 1 ]]; then
  fail 'Error message and usage not printed if no output'
  exit 1
fi
pass

# Test that standard usage prints a link to page
OUTPUT="$($BIN golang)"
STATUS=$?
if [[ $STATUS -ne 0 ]]; then
  fail 'Did not get success exit code'
  exit 1
fi
echo "$OUTPUT" | grep -q "Read more: https://en.wikipedia.org/wiki/Go"
if [ $? -ne 0 ];then
  fail 'Standard usage did not output link to page' 
fi
pass

# Test that short flag does not print a link to page
OUTPUT="$($BIN -s golang)"
STATUS=$?
if [[ $STATUS -ne 0 ]]; then
  fail 'Did not get success exit code'
  exit 1
fi
echo "$OUTPUT" | grep -q "Read more: https://en.wikipedia.org/wiki/Go"
if [ $? -eq 0 ];then
  fail 'Short flag did output link to page' 
fi
pass

# Test that language flag works
OUTPUT="$($BIN -l sv c++)"
STATUS=$?
if [[ $STATUS -ne 0 ]]; then
  fail 'Did not get success exit code'
  exit 1
fi
echo "$OUTPUT" | grep -q "Read more: https://sv.wikipedia.org/wiki/C"
if [ $? -ne 0 ];then
  fail 'Language flag did not work' 
fi
pass

# Test that no color flag works
OUTPUT="$($BIN -n golang)"
STATUS=$?
if [[ $STATUS -ne 0 ]]; then
  fail 'Did not get success exit code'
  exit 1
fi
echo "$OUTPUT" | grep -q "\[32m"
if [ $? -eq 0 ];then
  fail 'No color flag did not work' 
fi
pass

# Test that url flag works
OUTPUT="$($BIN -u http://localhost:8080/w/api.php golang 2>&1)"
STATUS=$?
if [[ $STATUS -eq 0 ]]; then
  fail 'Got success exit code'
  exit 1
fi
echo "$OUTPUT" | grep -q "Could not execute request Get http://localhost:8080"
if [ $? -ne 0 ];then
  fail 'Url flag did not work' 
fi
pass

# Test that no-check-certificate flag works
OUTPUT="$($BIN -no-check-certificate golang)"
STATUS=$?
if [[ $STATUS -ne 0 ]]; then
  fail 'Did not get success exit code'
  exit 1
fi
echo "$OUTPUT" | grep -q "Read more: https://en.wikipedia.org/wiki/Go"
if [ $? -ne 0 ];then
  fail 'Standard usage did not output link to page' 
fi
pass

# Test that short flag works
OUTPUT="$($BIN -short golang)"
STATUS=$?
if [[ $STATUS -ne 0 ]]; then
  fail 'Did not get success exit code'
  exit 1
fi
OUTPUT2="echo "$OUTPUT" | grep -c '.'"
if [ $OUTPUT2 -ne 1 ];then
  fail 'Short flag did not work' 
fi
pass

finished
