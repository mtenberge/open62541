cmake_minimum_required (VERSION 3.2)
project (typedef-extractor)

add_custom_command(
  OUTPUT extract_typedefs
  COMMAND cd ${CMAKE_CURRENT_SOURCE_DIR} && go build -o ${CMAKE_CURRENT_BINARY_DIR}/extract_typedefs
  DEPENDS
  extract_typedefs.go
)
add_custom_target(typedef-extractor DEPENDS extract_typedefs)

