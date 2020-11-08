# CFC Suggestions
Backend API for CFC Suggestions
## Configuration
can be configured with a suggestions_config.yaml file or environment variables 
see example_config.yaml for an example of what configuration values are available 
environment variables can be used be replacing - with _ and uppercasing the name
e.g.
`suggestions-channel` becomes `SUGGESTIONS_CHANNEL`

## Docker Example
#### building
`docker build -t cfc_suggestions .`

#### running
```bash
docker run --name cfc_suggestions \
-e DATABASE_FILE=/var/cfc_suggestions/database.db \
-e PORT=5023 \
-v /var/cfc_suggestions/test.db:/var/cfc_suggestions/database.db \
cfc_suggestions
```


## endpoints
- `POST /suggestions` 

    example request
    ```json
    {
      "owner": "179237013373845504"
    }
    ```
    example response
    ```json
    {
        "identifier": "312c251b27de46c3a84c69482ebcbd59",
        "owner": "179237013373845504"
    }
    ```
- POST /suggestions/{id}/send

    example request
    ```json
    {
        "title": "My Title",
        "description": "This is a description"
    }
    ```
    example response
    ```json
    {
        "status": "success"
    }
    ```
    
- POST /suggestions/{id}

    example response
    ```json
    {
        "identifier": "312c251b27de46c3a84c69482ebcbd59",
        "owner": "179237013373845504"
    }
    ```
    
- GET /suggestions
    
