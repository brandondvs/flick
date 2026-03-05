#!/bin/bash

BASE_URL="${1:-http://localhost:8080}"
PASSED=0
FAILED=0

red="\033[0;31m"
green="\033[0;32m"
reset="\033[0m"

assert_status() {
    local description="$1"
    local expected="$2"
    local actual="$3"
    local body="$4"

    if [ "$expected" -eq "$actual" ]; then
        echo -e "${green}PASS${reset} $description (status $actual)"
        ((PASSED++))
    else
        echo -e "${red}FAIL${reset} $description (expected $expected, got $actual)"
        echo "  response: $body"
        ((FAILED++))
    fi
}

assert_json_field() {
    local description="$1"
    local field="$2"
    local expected="$3"
    local body="$4"

    local actual
    actual=$(echo "$body" | grep -o "\"$field\":[^,}]*" | head -1 | cut -d: -f2- | tr -d '"' | tr -d ' ')

    if [ "$actual" = "$expected" ]; then
        echo -e "${green}PASS${reset} $description ($field=$actual)"
        ((PASSED++))
    else
        echo -e "${red}FAIL${reset} $description (expected $field=$expected, got $field=$actual)"
        echo "  response: $body"
        ((FAILED++))
    fi
}

echo "=== Feature Flag CRUD Lifecycle Test ==="
echo "Target: $BASE_URL"
echo ""

# 1. List flags (should be empty or baseline)
echo "--- LIST (empty) ---"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/flags")
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "GET /flags returns 200" 200 "$status" "$body"

# 2. Create a flag
echo ""
echo "--- CREATE ---"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/flags" \
    -H "Content-Type: application/json" \
    -d '{"name":"test-flag","enabled":false}')
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "POST /flags returns 201" 201 "$status" "$body"
assert_json_field "created flag name" "name" "test-flag" "$body"
assert_json_field "created flag enabled" "enabled" "false" "$body"

# 3. Create duplicate (should fail)
echo ""
echo "--- CREATE DUPLICATE ---"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/flags" \
    -H "Content-Type: application/json" \
    -d '{"name":"test-flag","enabled":true}')
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "POST /flags duplicate returns 409" 409 "$status" "$body"

# 4. Get the flag
echo ""
echo "--- READ ---"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/flags/test-flag")
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "GET /flags/test-flag returns 200" 200 "$status" "$body"
assert_json_field "read flag name" "name" "test-flag" "$body"
assert_json_field "read flag enabled" "enabled" "false" "$body"

# 5. Update the flag
echo ""
echo "--- UPDATE ---"
response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/flags/test-flag" \
    -H "Content-Type: application/json" \
    -d '{"enabled":true}')
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "PUT /flags/test-flag returns 200" 200 "$status" "$body"
assert_json_field "updated flag enabled" "enabled" "true" "$body"

# 6. Confirm update persisted
echo ""
echo "--- READ AFTER UPDATE ---"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/flags/test-flag")
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "GET /flags/test-flag returns 200" 200 "$status" "$body"
assert_json_field "flag still enabled" "enabled" "true" "$body"

# 7. Delete the flag
echo ""
echo "--- DELETE ---"
response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/flags/test-flag")
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "DELETE /flags/test-flag returns 204" 204 "$status" "$body"

# 8. Confirm deletion
echo ""
echo "--- READ AFTER DELETE ---"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/flags/test-flag")
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "GET /flags/test-flag returns 404" 404 "$status" "$body"

# 9. Delete non-existent flag
echo ""
echo "--- DELETE NON-EXISTENT ---"
response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/flags/test-flag")
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "DELETE /flags/test-flag returns 404" 404 "$status" "$body"

# 10. Update non-existent flag
echo ""
echo "--- UPDATE NON-EXISTENT ---"
response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/flags/test-flag" \
    -H "Content-Type: application/json" \
    -d '{"enabled":true}')
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "PUT /flags/test-flag returns 404" 404 "$status" "$body"

# 11. Create with missing name
echo ""
echo "--- CREATE MISSING NAME ---"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/flags" \
    -H "Content-Type: application/json" \
    -d '{"enabled":true}')
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "POST /flags missing name returns 400" 400 "$status" "$body"

# 12. Create with invalid body
echo ""
echo "--- CREATE INVALID BODY ---"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/flags" \
    -H "Content-Type: application/json" \
    -d 'not json')
body=$(echo "$response" | sed '$d')
status=$(echo "$response" | tail -1)
assert_status "POST /flags invalid body returns 400" 400 "$status" "$body"

# Summary
echo ""
echo "=== Results ==="
echo -e "${green}Passed: $PASSED${reset}"
echo -e "${red}Failed: $FAILED${reset}"

if [ "$FAILED" -gt 0 ]; then
    exit 1
fi
