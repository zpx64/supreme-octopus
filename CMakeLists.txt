cmake_minimum_required(VERSION 3.20)
cmake_policy(SET CMP0069 NEW)
project(supreme-octopus CXX)

# FetchContent support
include(FetchContent)

# i am too lazy to change my code
if(CMAKE_BUILD_TYPE STREQUAL "Debug")
  set(DEBUG true)
else()
  set(DEBUG false)
endif()

# hard code clang(
#set(CMAKE_C_COMPILER "clang")
#set(CMAKE_CXX_COMPILER "clang++")

set(SOURCE_EXE src/main.cc)

add_compile_options(
    -Wall -Wextra -pedantic
    -Wformat=2 -Wno-unused-parameter -Wshadow
    -Wwrite-strings -Wredundant-decls
    -Wmissing-include-dirs -Wno-format-nonliteral
)
if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
  add_compile_options(
    -Wlogical-op
  )
endif()

if(DEBUG)
  add_compile_options(
    -O0 -g3 -fstack-protector -ftrapv -fwrapv
    -fsanitize=address,undefined
    -DDEBUG -DTRACE
  )
  add_link_options(-fsanitize=address,undefined)
else()
  add_compile_options(
    -O3 -DNDEBUG -DNTRACE
  )
endif()

# ---libs---
#find_package(OpenSSL REQUIRED)
#find_package(Brotli REQUIRED)
#find_package(ZLIB REQUIRED)

set(HTTPLIB_COMPILE true)
set(HTTPLIB_REQUIRE_BROTLI true)
set(HTTPLIB_REQUIRE_ZLIB true)
set(HTTPLIB_REQUIRE_OPENSSL true)

FetchContent_Declare(
        httplib
        GIT_REPOSITORY "https://github.com/yhirose/cpp-httplib"
        GIT_TAG "97ae6733ed1d39a074a063db21aae460ec78181d"
)
FetchContent_MakeAvailable(httplib)
# ----------

add_executable(${PROJECT_NAME} ${SOURCE_EXE})
set_target_properties(${PROJECT_NAME} PROPERTIES CXX_STANDARD 23)
set_target_properties(${PROJECT_NAME} PROPERTIES CXX_STANDARD_REQUIRED true)
set_target_properties(${PROJECT_NAME} PROPERTIES CXX_EXTENSIONS false)

target_link_libraries(${PROJECT_NAME}
  PRIVATE
  httplib::httplib
  -static-libgcc -static-libstdc++
)

if(NOT DEBUG)
  set_target_properties(${PROJECT_NAME} PROPERTIES INTERPROCEDURAL_OPTIMIZATION true)
endif()
