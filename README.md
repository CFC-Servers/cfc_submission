# CFC Suggestions
Backend API for CFC Suggestions
Temporary4
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
    
