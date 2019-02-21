####################
# Nodeset Examples #
####################

###################
# Custom XML      #
###################

set(FILE_CSV_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)
set(FILE_BSD_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)
set(FILE_NS_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)

if(NOT CMAKE_PROJECT_NAME STREQUAL PROJECT_NAME)
    # needed or cmake doesn't recognize dependencies of generated files
    set(PROJECT_BINARY_DIR ${CMAKE_BINARY_DIR})
endif()

# generate namespace from XML file
if(NOT CMAKE_PROJECT_NAME STREQUAL PROJECT_NAME)
    ua_generate_nodeset_and_datatypes(
        NAME "example"
        FILE_NS "${FILE_NS_DIRPREFIX}/server_nodeset.xml"
        DEPENDS "${CMAKE_SOURCE_DIR}/tools/schema/Opc.Ua.NodeSet2.Minimal.xml"
    )
else()
    # standalone examples build expects already installed Opc.Ua.NodeSet2.Minimal.xml
    ua_generate_nodeset_and_datatypes(
        NAME "example"
        FILE_NS "${FILE_NS_DIRPREFIX}/server_nodeset.xml"
    )
endif()

# The .csv file can be created from within UaModeler or manually
ua_generate_nodeid_header(
    NAME "example_nodeids"
    ID_PREFIX "EXAMPLE_NS"
    TARGET_SUFFIX "ids_example"
    FILE_CSV "${FILE_CSV_DIRPREFIX}/server_nodeset.csv"
)

add_example(server_nodeset server_nodeset.c
            ${UA_NODESET_EXAMPLE_SOURCES}
            ${PROJECT_BINARY_DIR}/src_generated/open62541/example_nodeids.h)

if(UA_NAMESPACE_ZERO STREQUAL "FULL")
    ua_generate_nodeset_and_datatypes(
        NAME "testnodeset"
        FILE_CSV "${FILE_CSV_DIRPREFIX}/testnodeset.csv"
        FILE_BSD "${FILE_BSD_DIRPREFIX}/testtypes.bsd"
        NAMESPACE_IDX 2
        FILE_NS "${FILE_NS_DIRPREFIX}/testnodeset.xml"
        INTERNAL
    )
    add_example(server_testnodeset server_testnodeset.c
                ${UA_NODESET_TESTNODESET_SOURCES}
                ${UA_TYPES_TESTNODESET_SOURCES})
    add_dependencies(server_testnodeset open62541-generator-ns-testnodeset)
endif()


###################
# PLCopen Nodeset #
###################

if(NOT CMAKE_PROJECT_NAME STREQUAL PROJECT_NAME)
    set(FILE_CSV_DIRPREFIX ${CMAKE_SOURCE_DIR}/deps/ua-nodeset)
    set(FILE_BSD_PLCOPEN_DIRPREFIX ${CMAKE_SOURCE_DIR}/deps/ua-nodeset)
    set(FILE_BSD_POWERLINK_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)
    set(FILE_BSD_SIMATIC_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)
    set(FILE_NS_DIRPREFIX ${CMAKE_SOURCE_DIR}/deps/ua-nodeset)
else()
    set(FILE_CSV_DIRPREFIX ${UA_NODESET_DIR})
    set(FILE_BSD_PLCOPEN_DIRPREFIX ${UA_NODESET_DIR})
    set(FILE_BSD_POWERLINK_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)
    set(FILE_BSD_SIMATIC_DIRPREFIX ${PROJECT_SOURCE_DIR}/nodeset)
    set(FILE_NS_DIRPREFIX ${UA_NODESET_DIR})
endif()

# PLCopen requires the full ns0 as basis
if(UA_NAMESPACE_ZERO STREQUAL "FULL")

    # Generate types and namespace for DI
    ua_generate_nodeset_and_datatypes(
        NAME "di"
        FILE_CSV "${FILE_CSV_DIRPREFIX}/DI/OpcUaDiModel.csv"
        FILE_BSD "${FILE_BSD_PLCOPEN_DIRPREFIX}/DI/Opc.Ua.Di.Types.bsd"
        NAMESPACE_IDX 2
        FILE_NS "${FILE_NS_DIRPREFIX}/DI/Opc.Ua.Di.NodeSet2.xml"
        INTERNAL
    )

    # generate PLCopen namespace which is using DI
    ua_generate_nodeset_and_datatypes(
        NAME "plc"
        # PLCopen does not define custom types. Only generate the nodeset
        FILE_NS "${FILE_NS_DIRPREFIX}/PLCopen/Opc.Ua.Plc.NodeSet2.xml"
        # PLCopen depends on the di nodeset, which must be generated before
        DEPENDS "di"
        INTERNAL
    )

    add_example(server_nodeset_plcopen server_nodeset_plcopen.c
                ${UA_NODESET_DI_SOURCES}
                ${UA_NODESET_PLC_SOURCES}
                ${UA_TYPES_DI_SOURCES}
                ${UA_TYPES_PLC_SOURCES})
    add_dependencies(server_nodeset_plcopen open62541-generator-ns-plc)


endif()

# POWERLINK requires the full ns0 as basis

if(UA_NAMESPACE_ZERO STREQUAL "FULL")

    # generate powerlink namespace which is using DI
    ua_generate_nodeset_and_datatypes(
        NAME "powerlink"
        FILE_CSV "${FILE_CSV_DIRPREFIX}/POWERLINK/Opc.Ua.POWERLINK.NodeIds.csv"
        FILE_BSD "${FILE_BSD_POWERLINK_DIRPREFIX}/Opc.Ua.POWERLINK.NodeSet2.bsd"
        NAMESPACE_IDX 3
        FILE_NS "${FILE_NS_DIRPREFIX}/POWERLINK/Opc.Ua.POWERLINK.NodeSet2.xml"
        # POWERLINK depends on the di nodeset, which must be generated before
        DEPENDS "di"
        INTERNAL
    )

    add_example(server_nodeset_powerlink server_nodeset_powerlink.c
                ${UA_NODESET_DI_SOURCES}
                ${UA_NODESET_POWERLINK_SOURCES}
                ${UA_TYPES_DI_SOURCES}
                ${UA_TYPES_POWERLINK_SOURCES})
    add_dependencies(server_nodeset_powerlink open62541-generator-ns-powerlink)


endif()


# SIMATIC requires the full ns0 as basis
if(UA_NAMESPACE_ZERO STREQUAL "FULL")
  # Generate types for simatic namespace

  # Generate types and namespace for DI
  #ua_generate_nodeset_and_datatypes(
  #    NAME "di"
  #    FILE_CSV "${PROJECT_SOURCE_DIR}/deps/ua-nodeset/DI/OpcUaDiModel.csv"
  #    FILE_BSD "${PROJECT_SOURCE_DIR}/deps/ua-nodeset/DI/Opc.Ua.Di.Types.bsd"
  #    NAMESPACE_IDX 2
  #    FILE_NS "${PROJECT_SOURCE_DIR}/deps/ua-nodeset/DI/Opc.Ua.Di.NodeSet2.xml"
  #    INTERNAL
  #)

  # prepare/fix the original unmodified Siemens OPC-UA XML file (see fix_opc_xml.sh):
  add_custom_command(
    OUTPUT
    ${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml
    COMMAND cp ${CMAKE_CURRENT_SOURCE_DIR}/testPLC.OPCUA.xml ${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml
    COMMAND ${CMAKE_CURRENT_SOURCE_DIR}/fix_opc_xml.sh ${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml
    DEPENDS
    ${CMAKE_CURRENT_SOURCE_DIR}/testPLC.OPCUA.xml
    ${CMAKE_CURRENT_SOURCE_DIR}/fix_opc_xml.sh
  )
  add_custom_target(prepare-xml-simatic DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml)

  # extract the TypeDictionary for namespace 3 from the XML file:
  add_subdirectory(${CMAKE_SOURCE_DIR}/tools/extract_typedictionary ${CMAKE_CURRENT_BINARY_DIR}/extract_typedictionary)
  add_custom_command(
    OUTPUT
    ${CMAKE_CURRENT_BINARY_DIR}/simatic.types.bsd
    COMMAND ${CMAKE_CURRENT_BINARY_DIR}/extract_typedictionary/extract_typedictionary
      ${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml
      "ns=3;i=6001"
      ${CMAKE_CURRENT_BINARY_DIR}/simatic.types.bsd
    DEPENDS
    prepare-xml-simatic
    typedictionary-extractor
    VERBATIM
  )
  add_custom_target(create-bsd-simatic DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/simatic.types.bsd)

  # generate the csv-file:
  add_subdirectory(${CMAKE_SOURCE_DIR}/tools/extract_typedefs ${CMAKE_CURRENT_BINARY_DIR}/extract_typedefs)
  add_custom_command(
    OUTPUT
    ${CMAKE_CURRENT_BINARY_DIR}/OpcUaSimaticModel.csv
    COMMAND ${CMAKE_CURRENT_BINARY_DIR}/extract_typedefs/extract_typedefs
      ${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml
      ${CMAKE_CURRENT_BINARY_DIR}/OpcUaSimaticModel.csv
    DEPENDS
      prepare-xml-simatic
      typedef-extractor
  )
  add_custom_target(create-csv-simatic DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/OpcUaSimaticModel.csv)

  ua_generate_nodeset_and_datatypes(
    NAME "simatic"
    FILE_BSD "${CMAKE_CURRENT_BINARY_DIR}/simatic.types.bsd"
    FILE_CSV "${CMAKE_CURRENT_BINARY_DIR}/OpcUaSimaticModel.csv"
    NAMESPACE_IDX 3
    INTERNAL
    FILE_NS "${CMAKE_CURRENT_BINARY_DIR}/testPLC.OPCUA.xml"
    DEPENDS "di"
  )
  add_dependencies(open62541-generator-ns-simatic create-bsd-simatic create-csv-simatic prepare-xml-simatic)

  add_example(server_nodeset_simatic server_nodeset_simatic.c
    ${UA_NODESET_DI_SOURCES}
    ${UA_NODESET_SIMATIC_SOURCES}
    ${UA_TYPES_DI_SOURCES}
    ${UA_TYPES_SIMATIC_SOURCES}
  )
  add_dependencies(server_nodeset_simatic open62541-generator-ns-simatic)

endif()
