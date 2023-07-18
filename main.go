package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type JSONData struct {
	Header HeaderInfo   `json:"header"`
	Items  []MemoryData `json:"items"`
}

type HeaderInfo struct {
	Version int `json:"version"`
}

type MemoryData struct {
	Name           string           `json:"name"`
	Ext            Ext              `json:"ext"`
	IndicatorTypes []IndicatorType  `json:"indicatorTypes"`
}

type Ext struct {
	SNMP SNMP `json:"snmp"`
}

type SNMP struct {
	NameExpression       string `json:"nameExpression"`
	DescriptionExpression string `json:"descriptionExpression"`
	IndexOid             string `json:"indexOid"`
}

type IndicatorType struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Format               string `json:"format"`
	DataUnits            string `json:"dataUnits"`
	SyntheticExpression  string `json:"syntheticExpression"`
	Ext                  IndicatorTypeExt `json:"ext"`
}

type IndicatorTypeExt struct {
	SNMP IndicatorTypeSNMP `json:"snmp"`
}

type IndicatorTypeSNMP struct {
	Expression          string `json:"expression"`
	MaxValueExpression  string `json:"maxValueExpression"`
}

func convertJSONToCSV(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	var jsonData JSONData
	if err := json.NewDecoder(sourceFile).Decode(&jsonData); err != nil {
		fmt.Println(err)
		return err
	}

	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{
		"Name",
		"NameExpression",
		"DescriptionExpression",
		"IndexOid",
		"IndicatorName",
		"IndicatorDescription",
		"Format",
		"DataUnits",
		"SyntheticExpression",
		"Expression",
		"MaxValueExpression",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, data := range jsonData.Items {
		for _, indicator := range data.IndicatorTypes {
			row := []string{
				data.Name,
				data.Ext.SNMP.NameExpression,
				data.Ext.SNMP.DescriptionExpression,
				data.Ext.SNMP.IndexOid,
				indicator.Name,
				indicator.Description,
				indicator.Format,
				indicator.DataUnits,
				indicator.SyntheticExpression,
				indicator.Ext.SNMP.Expression,
				indicator.Ext.SNMP.MaxValueExpression,
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	if err := convertJSONToCSV("SNMP_OID_Certification.json", "data.csv"); err != nil {
		fmt.Println("Code doesn't work!!",err)
	}
}
