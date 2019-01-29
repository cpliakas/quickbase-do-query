# Quick Base Do Query

A command line tool that gets records from a Quick Base table.

## Usage

First, familiarize yourself with the [Quick Base query syntax](https://help.quickbase.com/api-guide/do_query.html#queryOperators).

Next, set the following configuration through environment variables, replacing
`[USER_TOKEN]`, `[REALM_NAME]`, and `[APP_ID]` according to your environment.

```sh
export QUICKBASE_USER_TOKEN="[USER_TOKEN]"
export QUICKBASE_REALM_HOST="https://[REALM_NAME].quickbase.com"
export QUICKBASE_APP_ID="[APP_ID]"
```

The example below returns records where field `7` exactly matches the value `Find Me`.
Replace `[TABLE_ID]` according to your environment.

```sh
quickbase-do-query --table-id="[TABLE_ID]" --query="{7.EX.'Find me'}"
```

You should see output similar to the example below:

```json
{
    "records": [
        {
            "record-id": 1,
            "update-id": 1548209252934,
            "fields": {
                "Match Field": "Find me",
                "Another Field": "Some value 1"
            }
        },
        {
            "record-id": 2,
            "update-id": 1548194623663,
            "fields": {
                "Match Field": "Find me",
                "Another Field": "Some value 1"
            }
        }
    ]
}
```

Can't remember what the numeric field IDs are? Run the following command:

```sh
quickbase-do-query list-fields --table-id="[TABLE_ID]"
```

```json
{
    "fields": {
        "7": "Match Field",
        "8": "Another Field"
    }
}
```
