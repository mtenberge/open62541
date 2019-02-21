/*
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

###
### Author:
### - Matthijs H. ten Berge (m.tenberge@awl.nl)
###
### This program was created for educational purposes and has been
### contributed to the open62541 project by the author. All licensing
### terms for this source is inherited by the terms and conditions
### specified for by the open62541 project (see the projects readme
### file for more information on the MPLv2 terms and restrictions).
*/

/*
This stand-alone executable is used to extract the type definitions out of a UANodeSet XML.

The XML is read and parsed in streaming mode, so performance with huge XML files should be ok.
*/
package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func printUsage() {
	log.Println("Usage:")
	log.Println("  extract_typedictionary <source file> <output file>")
	log.Println("    source file: the source filename containing a UANodeSet in XML format")
	log.Println("    output file: the CSV-file to which the data types will be *appended*")
	log.Println()
}

type UADataType struct {
	NodeId      string `xml:",attr"`
	DisplayName string
}

var reNamespace *regexp.Regexp = regexp.MustCompile("^ns=[0-9]+;")
var reNamespaceAndI *regexp.Regexp = regexp.MustCompile("^ns=[0-9]+;i=")
var containsSpecialChars *regexp.Regexp = regexp.MustCompile("[^\\w]")

func processOutputLine(outfile *os.File, node UADataType) {
	displayName := node.DisplayName
	dataType := node.NodeId

	// Remove the namespace part from the data type name:
	dataType = reNamespaceAndI.ReplaceAllString(dataType, "") // for integer IDs, remove everything except the number
	dataType = reNamespace.ReplaceAllString(dataType, "")     // for string IDs, keep the "s=" part

	// ugly hack: if the type name contains quotes, then make sure the display name also has quotes:
	if strings.ContainsRune(dataType, '"') && !strings.ContainsRune(displayName, '"') {
		displayName = fmt.Sprintf("\"%s\"", displayName)
	}

	// double the quotes:
	displayName = strings.Replace(displayName, "\"", "\"\"", -1)
	dataType = strings.Replace(dataType, "\"", "\"\"", -1)

	// if the identifiers contain non-alphanumerical characters, add surrounding quotes:
	if containsSpecialChars.MatchString(displayName) {
		displayName = fmt.Sprintf("\"%s\"", displayName)
	}
	if containsSpecialChars.MatchString(dataType) {
		dataType = fmt.Sprintf("\"%s\"", dataType)
	}

	_, err := outfile.WriteString(fmt.Sprintf("%s,%s,DataType\n", displayName, dataType))
	if err != nil {
		log.Fatalf("Cannot write output file: %s", err.Error())
	}
}

func main() {
	cmdlineArgs := os.Args[1:]
	var err error

	if len(cmdlineArgs) != 2 {
		log.Println("Invalid number of command line arguments specified")
		printUsage()
		return
	}

	log.Printf("Opening input file %s\n", cmdlineArgs[0])
	infile, err := os.Open(cmdlineArgs[0])
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	defer infile.Close()

	decoder := xml.NewDecoder(infile)

	log.Printf("Opening output file %s\n", cmdlineArgs[1])
	outfile, err := os.OpenFile(cmdlineArgs[1], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	defer outfile.Close()

	// first search for the top-level UANodeSet:
	log.Printf("Searching for UADataType nodes\n")
	for {
		token, err := decoder.Token()
		if token == nil {
			log.Printf("End of file encountered, done!\n")
			return
		}
		if err != nil {
			log.Fatalf("XML decoder error: %s\n", err.Error())
		}
		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "UADataType" {
				var node UADataType
				err = decoder.DecodeElement(&node, &se)
				if err != nil {
					log.Fatalf("DecodeElement failed: %s\n", err.Error())
				}
				if node.NodeId == "" {
					log.Printf("Found UADataType without NodeId, skipping\n")
					decoder.Skip()
				}
				if node.DisplayName == "" {
					log.Printf("Found UADataType (NodeId %s) without DisplayName, skipping\n", node.NodeId)
					decoder.Skip()
				}
				processOutputLine(outfile, node)
			}
		}
	}
}
