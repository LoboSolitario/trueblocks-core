# minimum cmake version supported
cmake_minimum_required (VERSION 3.0)

# script project
project (build_scripts)

# ----------------------------------------------------------------------------------------
# Folders for various purposes
set (BIN_TEST_FOLDER "${BIN_TEST_DIR}")

# user facing stuff
file (COPY ${CMAKE_CURRENT_SOURCE_DIR}/pylint.py      DESTINATION ${BIN_TEST_FOLDER} FILE_PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE)

# compare and copy test cases
file (REMOVE ${BIN_TEST_FOLDER}/test.sh)
file (REMOVE ${BIN_TEST_FOLDER}/test-all.sh)
file (REMOVE ${BIN_TEST_FOLDER}/test-api.sh)

file (COPY ${CMAKE_CURRENT_SOURCE_DIR}/test.sh        DESTINATION ${BIN_TEST_FOLDER} FILE_PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE)
file (COPY ${CMAKE_CURRENT_SOURCE_DIR}/test-all.sh    DESTINATION ${BIN_TEST_FOLDER} FILE_PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE)
file (COPY ${CMAKE_CURRENT_SOURCE_DIR}/test-api.sh    DESTINATION ${BIN_TEST_FOLDER} FILE_PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE)
