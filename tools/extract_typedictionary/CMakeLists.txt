cmake_minimum_required (VERSION 3.2)
project (typedictionary-extractor)

add_custom_command(
  OUTPUT extract_typedictionary
  COMMAND cd ${CMAKE_CURRENT_SOURCE_DIR} && go build -o ${CMAKE_CURRENT_BINARY_DIR}/extract_typedictionary
  DEPENDS
  extract_typedictionary.go
)
add_custom_target(typedictionary-extractor DEPENDS extract_typedictionary)
