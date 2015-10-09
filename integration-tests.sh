#!/usr/bin/env bash

BASEFILE=$0
TESTS=""
# Keep track of failures
FAILED=0

# Bin to use in tests
BIN=./build/wiki

# Setup all test functions that should be executed.
# All function starting with the word "test" will be part of the test suite.
function getTests() {
  for i in `cat $BASEFILE | sed -n "/^function test/s/function \([a-zA-Z0-9_]*\).*/\1/p"`; do
    TESTS="$TESTS $i"
  done
}

# Output a fail message
function fail() {
  echo ''
  echo 'FAIL'
  echo $1
  let "FAILED += 1"
}

# Show progress when tests pass
function pass() {
  echo -n '.'
}


# Output a message and return with no error when all tests has passed
function finished() {
  echo ''
  if [[ $FAILED -ne 0 ]]; then
    echo 'There were failed integration tests'
  else
    echo 'All integration tests passed'
  fi
  exit $FAILED
}

# Test that error message and usage is printed if no query
function testNoArguments() {
  $BIN > /dev/null 2>&1
  STATUS=$?
  if [[ $STATUS -ne 1 ]]; then
    fail 'Error message and usage not printed if no output'
  fi
  pass
}

# Test that standard usage prints a link to page
function testSimpleSearch() {
  OUTPUT="$($BIN golang)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  echo "$OUTPUT" | grep -q "Read more: https://en.wikipedia.org/wiki/Go"
  if [ $? -ne 0 ];then
    fail 'Standard usage did not output link to page'
  fi
  pass
}

# Test that short flag does not print a link to page
function testShortFlag() {
  OUTPUT="$($BIN -s golang)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  echo "$OUTPUT" | grep -q "Read more: https://en.wikipedia.org/wiki/Go"
  if [ $? -eq 0 ];then
    fail 'Short flag did output link to page'
  fi
  pass
}

# Test that language flag works
function testLanguageFlag() {
  OUTPUT="$($BIN -l sv c++)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  echo "$OUTPUT" | grep -q "Read more: https://sv.wikipedia.org/wiki/C"
  if [ $? -ne 0 ];then
    fail 'Language flag did not work'
  fi
  pass
}

# Test that language enviroment works
function testLanaguageEnv() {
  OUTPUT="$(WIKI_LANG=sv $BIN c++)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  echo "$OUTPUT" | grep -q "Read more: https://sv.wikipedia.org/wiki/C"
  if [ $? -ne 0 ];then
    fail 'Language flag did not work'
  fi
  pass
}

# Test that no color flag works
function testNoColorFlag() {
  OUTPUT="$($BIN -n golang)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  echo "$OUTPUT" | grep -q "\[32m"
  if [ $? -eq 0 ];then
    fail 'No color flag did not work'
  fi
  pass
}

# Test that url flag works
function testURLFlag() {
  OUTPUT="$($BIN -u http://localhost:8080/w/api.php golang 2>&1)"
  STATUS=$?
  if [[ $STATUS -eq 0 ]]; then
    fail 'Got success exit code'
  fi
  echo "$OUTPUT" | grep -q "Could not execute request Get http://localhost:8080"
  if [ $? -ne 0 ];then
    fail 'Url flag did not work'
  fi
  pass
}

# Test URL passed as enviroment
function testURLEnv () {
  OUTPUT="$(WIKI_URL=http://localhost:8080/w/api.php $BIN golang 2>&1)"
  STATUS=$?
  if [[ $STATUS -eq 0 ]]; then
    fail 'Got success exit code'
  fi
  echo "$OUTPUT" | grep -q "Could not execute request Get http://localhost:8080"
  if [ $? -ne 0 ];then
    fail 'Url flag did not work'
  fi
  pass
}

# Test that no-check-certificate flag works
function testNoCheckCertificateFlag() {
  OUTPUT="$($BIN -no-check-certificate golang)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  echo "$OUTPUT" | grep -q "Read more: https://en.wikipedia.org/wiki/Go"
  if [ $? -ne 0 ];then
    fail 'Standard usage did not output link to page'
  fi
  pass
}

# Test that short flag works
function testShortFlag2() {
  OUTPUT="$($BIN -short golang)"
  STATUS=$?
  if [[ $STATUS -ne 0 ]]; then
    fail 'Did not get success exit code'
  fi
  OUTPUT2="$(echo $OUTPUT | grep -c '.')"
  if [ $OUTPUT2 -ne 1 ];then
    fail 'Short flag did not work'
  fi
  pass
}

function main() {
  getTests

  for t in $TESTS; do
      $t
  done
  finished
}

main
