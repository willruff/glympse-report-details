# Original code by Will Ruff, written for Glympse, Inc in 2024.
# Open Sourced with permission.
import csv
import os
import boto3
from datetime import datetime, timedelta
from io import StringIO
from report_details.logger.json_logger import get_logger, with_fields

#Replace prints with logs
def sortFile(s3_data, search_filter, column_select):
    # Check if there is data
    if s3_data is None or s3_data == "":
        print(f"ERROR: No data was provided")
        return None, None

    export_data = []
    reader = csv.DictReader(StringIO(s3_data))

    field_names = reader.fieldnames

    # Check if column_select is in the CSV header
    if column_select not in reader.fieldnames:
        print(f"ERROR: Column '{column_select}' was not found in CSV data")
        return field_names, export_data
    
    for row in reader:
        if row[column_select] == search_filter:
            row.pop(column_select)
            export_data.append(row)
    
    field_names = [f for f in field_names if f != column_select]

    if not export_data:
        print(f"ERROR: No rows in the '{column_select}' column matched with '{search_filter}'")
        return field_names, export_data
    
    return field_names, export_data

def uploadFile(field_names, export_data, upload_key, bucket_name, s3_client):
    output = StringIO()
    writer = csv.DictWriter(output, fieldnames=field_names)    

    writer.writeheader()
    for row in export_data:
        writer.writerow(row)

    s3_client.put_object(Bucket=bucket_name, Key=upload_key, Body=output.getvalue())
    
def main():
    rpt_date = os.getenv("RPT_DATE", "")
    days_back = os.getenv("DAYS_BACK", "7")
    begin_date = os.getenv("BEGIN_DATE", "")
    end_date = os.getenv("END_DATE", "")
    org_download_id = os.getenv("ORG_DOWNLOAD_ID", "input-data")
    org_upload_id = os.getenv("ORG_UPLOAD_ID", "output-data")
    bucket_name = os.getenv("BUCKET_NAME", "reporting")
    search_filter = os.getenv("SEARCH_FILTER", "picture")
    column_select = os.getenv("COLUMN_SELECT", "id")
    #Change rpt_date to datetime object

    if rpt_date == "":
        rpt_date = (datetime.now() - timedelta(days=1)).strftime("%Y-%m-%d")

    #log all variable values
    l = get_logger()
    with_fields(l, rpt_date=rpt_date, begin_date=begin_date, end_date=end_date, days_back=days_back).info("Starting Up")
    with_fields(l, org_download_id=org_download_id, org_upload_id=org_upload_id, bucket_name=bucket_name, search_filter=search_filter, column_select=column_select).info("Starting Up")

    #Change rpt_date to datetime object
    rpt_datetime = datetime.strptime(rpt_date, "%Y-%m-%d")
    if begin_date == "":
        begin_date = (rpt_datetime - timedelta(days=int(days_back))).strftime("%Y-%m-%d")
    if end_date == "":
        end_date = (rpt_datetime - timedelta(days=1)).strftime("%Y-%m-%d")

    download_key = f"reports/{org_download_id}/date={rpt_date}/standard_report_{begin_date}_{end_date}_UTC.csv"
    with_fields(l, download_path=f"s3://{bucket_name}/{download_key}").info("Attempted download path")
    s3_client = boto3.client('s3')
    # Download the file from S3
    s3_object = s3_client.get_object(Bucket=bucket_name, Key=download_key)
    s3_data = s3_object['Body'].read().decode('utf-8')
    # Get the header names and the rest of the data from sortFile Function
    field_names, export_data = sortFile(s3_data, search_filter, column_select)

    upload_key = f"reports/{org_upload_id}/date={rpt_date}/report_details_{begin_date}_{end_date}_UTC.csv"
    uploadFile(field_names, export_data, upload_key, bucket_name, s3_client)
    with_fields(l, upload_path=f"s3://{bucket_name}/{upload_key}").info("Upload was successful!")

if __name__ == '__main__':
    main()