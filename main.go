// This file is in the public domain.

package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// HealthData is the top-level element of export.xml.
type HealthData struct {
	Records []Record `xml:"Record"`
}

// Record contains an individual record entered into HealthKit.
type Record struct {
	Type  string `xml:"type,attr"`
	Unit  string `xml:"unit,attr"`
	Date  string `xml:"creationDate,attr"`
	Value string `xml:"value,attr"`

	// other dates, etc. ignored
}

func fatalf(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", a...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		fatalf("Usage: %s [filename]", os.Args[0])
	}

	filename := os.Args[1]

	var data []byte

	extension := path.Ext(filename)
	switch extension {
	case ".zip":
		rc, err := zip.OpenReader(filename)
		if err != nil {
			fatalf("Couldn’t open zipfile named %s: %s", filename, err)
		}
		defer rc.Close()

		for _, f := range rc.File {
			if f.Name != "apple_health_export/export.xml" {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				fatalf("Couldn’t open export.xml inside %s: %s", filename, err)
			}
			defer rc.Close()

			data, err = ioutil.ReadAll(rc)
			if err != nil {
				fatalf("Couldn’t read all of export.xml inside %s: %s", filename, err)
			}
		}
	default:
		if extension != ".xml" {
			fmt.Fprintf(os.Stderr, "Warning: Could not tell whether %s is an XML file or a zipfile. Assuming it’s XML.\n", filename)
		}

		rc, err := os.Open(filename)
		if err != nil {
			fatalf("Couldn’t open XML file named %s: %s", filename, err)
		}
		defer rc.Close()

		data, err = ioutil.ReadAll(rc)
		if err != nil {
			fatalf("Couldn’t read all of file named %s: %s", filename, err)
		}

	}

	hd := HealthData{}
	err := xml.Unmarshal(data, &hd)
	if err != nil {
		fatalf("Couldn’t unmarshal health data into a buffer: %s", err)
	}

	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"Date", "Weight"})

	for _, record := range hd.Records {
		if record.Type != "HKQuantityTypeIdentifierBodyMass" {
			continue
		}

		var pounds string

		if record.Unit == "kg" {
			kg, err := strconv.ParseFloat(record.Value, 64)
			if err != nil {
				fatalf("Couldn’t parse mass in kilograms: %s", err)
			}
			pounds = strconv.FormatFloat(kg*2.204623, 'f', -1, 64)
		} else if record.Unit == "lb" {
			pounds = record.Value
		} else {
			fatalf("Couldn’t figure out how to convert “%s” to lb (pounds)", record.Unit)
		}

		entry := []string{record.Date, pounds}
		err = writer.Write(entry)
		if err != nil {
			fatalf("Couldn’t write line: %s", entry)
		}
	}

	writer.Flush()
	err = os.Stdout.Sync()
	if err != nil {
		fatalf("Couldn’t sync stdout")
	}
}
