package catalystx3ext

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

func Load_Sitetype_Data(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}
	//Open received file
	fileToImport, err := fileHeader.Open()
	if err != nil {
		ctx.Error(err)
		return
	}
	defer fileToImport.Close()

	//Reading the name of received file and creating a new file with the same name
	filenamenew := fileHeader.Filename
	os.Mkdir("/app/downloads", 0700)
	newFile, err := os.Create("/app/downloads/" + filenamenew)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer newFile.Close()

	// //Delete temp file after importing
	// defer os.Remove("/app/downloads/"+filenamenew)

	//Now Write data from received file to the newly created file
	fileBytes, err := io.ReadAll(fileToImport)
	if err != nil {
		ctx.Error(err)
		return
	}
	_, err = newFile.Write(fileBytes)
	if err != nil {
		ctx.Error(err)
		return
	}
	newFile.Close()
	//////////////////////////

	var rowExcel []string
	excelfile, err := excelize.OpenFile("/app/downloads/" + filenamenew)
	if err != nil {
		println(err.Error())
	}
	sheetName := excelfile.GetSheetName(1)
	allRowsExcel, err := excelfile.Rows(sheetName)
	if err != nil {
		log.Fatal(err)
		return
	}
	logfile, err := os.OpenFile("/home/ubuntu/catalyst/log/sheet2processed.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logfile.Close()
	logger := log.New(logfile, "prefix: ", log.LstdFlags)
	rowsProcessed := 0
	header := true
	for allRowsExcel.Next() {
		rowExcel = allRowsExcel.Columns()
		if header {
			if strings.ToLower(rowExcel[0]) == "cell ids_optima" && strings.ToLower(rowExcel[54]) == "monthly sites" {
				rowsProcessed++
				continue
			} else {
				break
			}
		}
		if !header {
			if len(rowExcel[46]) != 5 {
				logger.Print("Erroneous Property_ID", rowExcel, "\r\n")
				continue
			}
			rowsProcessed++
		}
		header = false
	}
	//////////////////////////
	ctx.JSON(http.StatusOK, "File "+filenamenew+" uploaded successfully")
	// ctx.JSON(http.StatusOK, "File uploaded successfully")
	time.Sleep(1 * time.Second)
}
