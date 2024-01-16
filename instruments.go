package icici_breeze_sdk

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"github.com/algotuners-v2/icici-breeze-sdk/breeze_models"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

func downloadIciciMasterScrips() {
	url := "https://directlink.icicidirect.com/NewSecurityMaster/SecurityMaster.zip"
	err := downloadFile(url, "SecurityMaster.zip")
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	err = extractZip("SecurityMaster.zip", "extracted")
	if err != nil {
		fmt.Println("Error extracting zip file:", err)
		return
	}
}

func getNseScrips() []breeze_models.NseScrip {
	currentWd, _ := os.Getwd()
	originalFilePath := path.Join(currentWd, "extracted", "NSEScripMaster.txt")
	csvFilePath, err := convertTextToCSV(originalFilePath, ',') // Specify the delimiter
	if err != nil {
		fmt.Println("Error converting text to CSV:", err)
		return nil
	}
	result := readNseScripCSVFile(csvFilePath)
	deleteDir()
	return result
}

func GetIciciInstruments() ([]breeze_models.NseScrip, []breeze_models.NseFnoScript) {
	downloadIciciMasterScrips()
	return getNseScrips(), getNseFnoScrips()
}

func getNseFnoScrips() []breeze_models.NseFnoScript {
	currentWd, _ := os.Getwd()
	originalFilePath := path.Join(currentWd, "extracted", "FONSEScripMaster.txt")
	csvFilePath, err := convertTextToCSV(originalFilePath, ',') // Specify the delimiter
	if err != nil {
		fmt.Println("Error converting text to CSV:", err)
		return nil
	}
	result := readNseFnoScripCSVFile(csvFilePath)
	deleteDir()
	return result
}

func downloadFile(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

func extractZip(zipFile, destFolder string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	os.MkdirAll(destFolder, os.ModePerm)
	for _, file := range reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()
		destFilePath := filepath.Join(destFolder, file.Name)
		extractedFile, err := os.Create(destFilePath)
		if err != nil {
			return err
		}
		defer extractedFile.Close()
		_, err = io.Copy(extractedFile, zippedFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertTextToCSV(textFilePath string, delimiter rune) (string, error) {
	file, err := os.Open(textFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the content of the text file
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Convert the content to CSV format by replacing spaces with the specified delimiter

	// Change the file extension to .csv
	csvFilePath := changeFileExtension(textFilePath, "csv")

	// Write the CSV content to the new file
	err = os.WriteFile(csvFilePath, []byte(content), os.ModePerm)
	if err != nil {
		return "", err
	}

	return csvFilePath, nil
}

func changeFileExtension(filePath, newExtension string) string {
	dir, file := filepath.Split(filePath)
	newFileName := fmt.Sprintf("%s.%s", file[:len(file)-len(filepath.Ext(file))], newExtension)
	return filepath.Join(dir, newFileName)
}

func extractNseScripsFromCseRecords(records [][]string) []breeze_models.NseScrip {
	lst := []breeze_models.NseScrip{}
	for _, record := range records[1:] {
		var nseScrip breeze_models.NseScrip
		err := nseScrip.MapCSVToStruct(record, &nseScrip, records[0])
		if err != nil {
			log.Fatal(err)
		}
		lst = append(lst, nseScrip)
	}
	return lst
}

func readNseScripCSVFile(filename string) []breeze_models.NseScrip {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	records[0][0] = ` "` + records[0][0] + `"`
	return extractNseScripsFromCseRecords(records)
}

func extractNseFnoScripsFromCseRecords(records [][]string) []breeze_models.NseFnoScript {
	lst := []breeze_models.NseFnoScript{}
	for _, record := range records[1:] {
		var nseFnoScrip breeze_models.NseFnoScript
		err := MapCSVToStruct(record, &nseFnoScrip, records[0])
		if err != nil {
			log.Fatal(err)
		}
		lst = append(lst, nseFnoScrip)
	}
	return lst
}

func readNseFnoScripCSVFile(filename string) []breeze_models.NseFnoScript {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	//records[0][0] = ` "` + records[0][0] + `"`
	return extractNseFnoScripsFromCseRecords(records)
}

func MapCSVToStruct(record []string, target interface{}, header []string) error {
	//r := csv.NewReader(strings.NewReader(strings.Join(header, ",")))
	//r.Comma = ','
	//r.LazyQuotes = true
	//// Read the CSV record
	//recordFields, err := r.Read()
	//if err != nil {
	//	return err
	//}

	// Create a map to store the mapping of struct field names to their index in the CSV record
	fieldMap := make(map[string]int)
	for i, field := range header {
		fieldMap[field] = i
	}

	// Use reflection to set struct fields based on the CSV record
	structValue := reflect.ValueOf(target).Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		csvTag := field.Tag.Get("csv")
		csvIndex, ok := fieldMap[csvTag]
		if ok {
			fieldValue := structValue.Field(i)
			fieldType := fieldValue.Type()

			// Convert the CSV field value to the appropriate struct field type
			switch fieldType.Kind() {
			case reflect.Int:
				intValue, err := strconv.Atoi(record[csvIndex])
				if err != nil {
					fieldValue.SetInt(int64(0))
				} else {
					fieldValue.SetInt(int64(intValue))
				}
			case reflect.Float64:
				floatValue, err := strconv.ParseFloat(record[csvIndex], 64)
				if err != nil {
					fieldValue.SetFloat(0)
				} else {
					fieldValue.SetFloat(floatValue)
				}
			case reflect.String:
				fieldValue.SetString(record[csvIndex])
			case reflect.Struct:
				timeValue, err := time.Parse("2006/01/02 15:04:05", record[csvIndex])
				if err != nil {
					fieldValue.Set(reflect.ValueOf(time.Time{}))
				} else {
					fieldValue.Set(reflect.ValueOf(timeValue))
				}
			}
		}
	}
	return nil
}

func deleteDir() {
	currentWd, _ := os.Getwd()
	dirPath := path.Join(currentWd, "extracted")
	zipPath := path.Join(currentWd, "SecurityMaster.zip")
	err := os.Remove(zipPath)
	err = os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Directory deleted successfully!")
}
