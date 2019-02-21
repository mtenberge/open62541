#!/bin/bash

# Replace all node types 62 by node type 63:
sed -i -e 's|<Reference ReferenceType="HasTypeDefinition">i=62<\/Reference>|<Reference ReferenceType="HasTypeDefinition">i=63</Reference>|g' $1
