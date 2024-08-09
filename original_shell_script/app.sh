#!/bin/bash
# Copied directly with permission from Glympse, Inc.

set -e

function log() {
  local code=$1
  local data_key=$2
  local data_val=$3

  jq -nc \
    --arg code "${code}" \
    --arg data_key "${data_key}" \
    --arg data_val "${data_val}" \
    '
      {
        event_type: "event",
        code: $code,
        data: {
          ($data_key): $data_val
        }
      }
    '
}

function dd_event() {
  evt_text=$1

  local json
  json=$(
    jq -nc \
      --arg t "${evt_text}" \
      '
        ("Error: "+$t) as $text |
        {
          title: "Failed to upload Customer data from report_details",
          text: $text,
          alert_type: "warning",
          aggregation_key: "report_details",
        }
      '
  )

  log "notification" "payload" "$json"
}

RPT_DATE=$1
BEGIN_DATE=$2
END_DATE=$3
[ -n "${RPT_DATE}" ] || RPT_DATE=$(date -d "yesterday" -I)
[ -n "${BEGIN_DATE}" ] || BEGIN_DATE=$(date -d "${RPT_DATE} - 7 days" -I)
[ -n "${END_DATE}" ] || END_DATE=$(date -d "${RPT_DATE} - 1 day" -I)

dlfile=$(mktemp -t report-dl.XXX)
tmpfile=$(mktemp -t report-tmp.XXX)
ulfile=$(mktemp -t report-ul.XXX)

# Make sure the begin and end dates are correct
rpt_date=$(date -I -d "$RPT_DATE") || (echo "Invalid report date!" && exit 1)
begin_date=$(date -I -d "$BEGIN_DATE") || (echo "Invalid begin date!" && exit 1)
end_date=$(date -I -d "$END_DATE")     || (echo "Invalid end date!" && exit 1)

do_download() {
  s3path="s3://reporting/reports/input-data/date=${rpt_date}/standard_report_${begin_date}_${end_date}_UTC.csv"

  aws s3 cp \
    "${s3path}" \
    "$dlfile"
}

do_upload() {
  s3path="s3://reporting/reports/output-data/date=${rpt_date}/report_details_${begin_date}_${end_date}_UTC.csv"

  aws s3 cp \
    "${ulfile}" \
    "${s3path}"
}

log "startup" "rpt_date" "${RPT_DATE}"
log "startup" "begin_date" "${BEGIN_DATE}"
log "startup" "end_date" "${END_DATE}"

# We don't want to exit out on errors from here down.
set +e

dl_debug=$(do_download 2>&1)
log "download" "output" "$dl_debug"

xsv search -s id 'picture' "$dlfile" > "$tmpfile"
xsv select '!id' "$tmpfile" > "$ulfile"

if [ -s "$ulfile" ]; then
  ul_debug=$(do_upload 2>&1)
  log "upload" "output" "$ul_debug"
else
  log "upload" "error" "Not uploading, zero byte file."
  dd_event "Not uploading, zero byte file."
fi
