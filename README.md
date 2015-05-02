Cassandra Log to SQLite DB
===================================

This script will consume a cassandra system.log and attempt to parse its rows into
a SQLite database for analysis.

Note: Multiline messages (e.g. exception trace information) are discarded after the first line.

### Usage

    go run main.go /var/log/cassandra/system.log


Once the file has been generated any sqlite client will be able to open it.
The Sqlite admin extension for Firefox is good.
