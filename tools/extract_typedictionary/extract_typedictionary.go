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
This stand-alone executable is used to extract the TypeDictionary section out of a UANodeSet XML.

The XML is read and parsed in streaming mode, so performance with huge XML files should be ok.
*/

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
	"log"
	"os"
)

func printUsage() {
	log.Println("Usage:")
	log.Println("  extract_typedictionary <source> <node-id> <output file>")
	log.Println("    source: the source filename containing a UANodeSet in XML format")
	log.Println("    node-id: the literal node-ID of the node containing the TypeDictionary (as it occurs in the XML), for example: ns=3;s=&quot;demoNodeName&quot;")
	log.Println()
}

type UAVariable struct {
	DisplayName string
	Value       struct {
		ByteString []byte
	}
}

func main() {
	cmdlineArgs := os.Args[1:]
	var err error

	if len(cmdlineArgs) != 3 {
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

	// first search for the top-level UANodeSet:
	log.Printf("Searching for the UAVariable with NodeID %s\n", cmdlineArgs[1])
	for {
		token, err := decoder.Token()
		if token == nil {
			log.Fatalf("End of file encountered before finding it\n")
		}
		if err != nil {
			log.Fatalf("XML decoder error: %s\n", err.Error())
		}
		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "UAVariable" {
				for _, attr := range se.Attr {
					if attr.Name.Local == "NodeId" {
						if attr.Value != cmdlineArgs[1] {
							err = decoder.Skip()
							break
						} else {
							// found!
							log.Printf("Found! Now unmarshalling this node\n")
							var node UAVariable
							err = decoder.DecodeElement(&node, &se)
							if err != nil {
								log.Fatalf("DecodeElement failed: %s\n", err.Error())
							}
							log.Printf("DisplayName: %s\n", node.DisplayName)

							log.Printf("Writing output file %s\n", cmdlineArgs[2])
							outfile, err := os.Create(cmdlineArgs[2])
							if err != nil {
								log.Fatalf("Error while creating output file: %s\n", err.Error())
							}
							defer outfile.Close()

							base64Decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(node.Value.ByteString))
							io.Copy(outfile, base64Decoder)

							return
						}
					}
				}
			}
		}
	}

}
