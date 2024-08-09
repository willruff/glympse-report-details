# Report Details

## Internship Expeirence (Glympse)
During my internship at Glympse, I worked on enhancing this weekly report for one of their customers. Initially, this task was handled by a bash script (`app.sh`) that used the `xsv` tool to manipulate CSV files. As part of my project, I converted this shell script into a more robust and maintainable solution using Python and Golang. The original bash script has been preserved in the `original_shell_script` directory for reference.

## Not Included in This Repository
For confidentiality reasons, several components from my work at Glympse cannot be open-sourced. These include:

* The `.gitlab-ci.yml` file, which was used to automate the CI/CD pipeline.
* The Makefile, which managed the build process.
* The Nomad job configurations that handled the deployment process.

## CI/CD Pipeline
In addition to refactoring the export process, I developed a comprehensive GitLab CI/CD pipeline. This pipeline:

* Ran unit tests to ensure the reliability of the code.
* Built a Docker image for consistent deployment.
* Deployed the application to Glympse's staging environment for testing.
* Successfully pushed the code into production.

These contributions not only modernized the existing process but also ensured a smoother and more reliable deployment pipeline.

## What Does This Repository Do?
The `report_details` export is a meta report generated weekly based on the output of the customers weekly `standard_report`. This report filters the `standard_report` by the rows in the `"id"` column that contain the word `"picture"`. The `"id"` column and corresponding values are then removed, and the CSV file is exported.

## How Does This Run?
This report is scheduled to execute weekly, providing the customer with a detailed report.

## Open Source with Permission

This repository contains open-source content that has been made available with the explicit permission of Glympse, Inc. All proprietary or confidential files have been excluded to comply with company policies.