# This log parser is a CLI tool that allows user to parse log the following format:

##Behavior
The behavior of the tool is the following
The tool will read one or several log files, given by user. The tool will extract only supported format logs and aggregate them.
The aggregations are the following:
Logs are aggregated by levels, showing number of lines for each level
Logs are aggregated by service, showing for each service:
- Number of lines
- Average duration of processes
The tool also shows number of errors encountered during parsing (non existing files, format errors, etc) but does not natively give details on those errors.
If the user is requesting slowest logs using slowestLogs flag, the tool will also give full log lines for slowest logs, depending on the number requested

## Logs Format
Currently, only one format of logs is supported :

DATE [LEVEL] service=SERVICE message=MESSAGE duration=DURATION PROPS

- **DATE** is date and time, RFC3339 format (YYYY\-MM\-DDTHH\:MM\:SSZ)
- **LEVEL** is a string in [DEBUG, INFO, WARNING, ERROR, FATAL]
- **SERVICE** is a string describing service logged
- **MESSAGE** is a message
- **DURATION** is the duration of the process, in milliseconds, in the format XXXXms
- **PROPS** is one or multiple key-value pairs, in the format key=value, separated by spaces


## Flags
When using this tool, it is possible to add flags to modify behavior. Some flags are necessary, some are optionals.
### Necessary flags
- **files**: log files to analyse
### Optional flags
- **startDate**: log date to start from
- **endDate**: log date to end to
- **service**: filter logs by service
- **level**: filter logs by level
- **top**: number of slowest logs to show
