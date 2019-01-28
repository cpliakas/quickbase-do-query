# Quick Base Do Query

A command line tool that gets records from a Quick Base table.

## Usage

First, familiarize yourself with the [Quick Base query syntax](https://help.quickbase.com/api-guide/do_query.html#queryOperators).

Next, set the following configuration with environment variables, replacing
`[USER_TOKEN]`, `[REALM_NAME]`, and `[APP_ID]` according to your environment.

```sh
export QUICKBASE_USER_TOKEN="[USER_TOKEN]"
export QUICKBASE_REALM_HOST="https://[REALM_NAME].quickbase.com"
export QUICKBASE_APP_ID="[APP_ID]"
```

The example below returns records where field `7` exactly matches the value "Find Me".
Replace `[TABLE_ID]` according to your environment.

```sh
quickbase-do-query --table-id="[TABLE_ID]" --query="{7.EX.'Find me'}"
```
