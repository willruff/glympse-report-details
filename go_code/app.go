// Original code by Will Ruff, written for Glympse, Inc in 2024.
// Open Sourced with permission.
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func fileDownload(sess *session.Session, bucket string, key string) [][]string {
	svc := s3.New(sess)
	// Get the object from S3
	file, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		panic(err)
	}
	defer file.Body.Close()

	reader := csv.NewReader(file.Body)
	input_data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	//fmt.Print(input_data)
	return input_data
}

func fileSort(search_filter string, column_select string, input_data [][]string) (filtered_data [][]string, target int) {
	target = -1
	if len(input_data) == 0 {
		filtered_data = nil
		return filtered_data, target
	}

	// Add header row
	filtered_data = append(filtered_data, input_data[0])

	for i := 0; i < len(input_data[0]); i++ {
		if input_data[0][i] == column_select {
			target = i
		}
	}
	if target == -1 {
		return filtered_data, target
	}

	for i := 1; i < len(input_data); i++ {
		if input_data[i][target] == search_filter {
			filtered_data = append(filtered_data, input_data[i])
		}
	}

	return filtered_data, target
}

func columnRemoval(filtered_data [][]string, target int) (export_data [][]string) {
	if len(filtered_data) == 0 || target < 0 || target >= len(filtered_data[0]) {
		return filtered_data // Return the original data if it's empty or the column index is out of range
	}
	// Create a new 2D slice to store the result
	export_data = make([][]string, len(filtered_data))

	// Iterate over each row
	for i := 0; i < len(filtered_data); i++ {
		for j := 0; j < len(filtered_data[i]); j++ {
			if j != target {
				export_data[i] = append(export_data[i], filtered_data[i][j])
			}
		}
	}
	return export_data
}

func fileUpload(sess *session.Session, bucket string, key string, export_data [][]string) {
	svc := s3.New(sess)

	var buffer bytes.Buffer

	// Create a CSV writer
	writer := csv.NewWriter(&buffer)
	err := writer.WriteAll(export_data)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success!")
}

func main() {
	rpt_date := os.Getenv("RPT_DATE")
	if rpt_date == "" {
		rpt_date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	rpt_date_time, err := time.Parse("2006-01-02", rpt_date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	days_back := os.Getenv("DAYS_BACK")
	if days_back == "" {
		days_back = "7"
	}
	// Convert the string to an integer
	days_back_int, err := strconv.Atoi(days_back)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}

	begin_date := os.Getenv("BEGIN_DATE")
	if begin_date == "" {
		begin_date = rpt_date_time.AddDate(0, 0, -days_back_int).Format("2006-01-02")
	}

	end_date := os.Getenv("END_DATE")
	if end_date == "" {
		end_date = rpt_date_time.AddDate(0, 0, -1).Format("2006-01-02")
	}

	org_download_id := os.Getenv("ORG_DOWNLOAD_ID")
	if org_download_id == "" {
		org_download_id = "input-data"
	}

	org_upload_id := os.Getenv("ORG_UPLOAD_ID")
	if org_upload_id == "" {
		org_upload_id = "output-data"
	}

	bucket_name := os.Getenv("BUCKET_NAME")
	if bucket_name == "" {
		bucket_name = "reporting"
	}
	search_filter := os.Getenv("SEARCH_FILTER")
	if search_filter == "" {
		search_filter = "picture"
	}
	column_select := os.Getenv("COLUMN_SELECT")
	if column_select == "" {
		column_select = "id"
	}

	download_key := fmt.Sprintf("reports/%s/date=%s/standard_report_%s_%s_UTC.csv", org_download_id, rpt_date, begin_date, end_date)
	upload_key := fmt.Sprintf("reports/%s/date=%s/report_details_%s_%s_UTC.csv", org_upload_id, rpt_date, begin_date, end_date)

	// Initialize a session using the default credentials and region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	input_data := fileDownload(sess, bucket_name, download_key)
	filtered_data, target := fileSort(search_filter, column_select, input_data)
	export_data := columnRemoval(filtered_data, target)
	fileUpload(sess, bucket_name, upload_key, export_data)
}
