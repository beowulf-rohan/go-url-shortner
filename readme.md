# URL Shortener Application

This application provides a Dockerized URL Shortener service built using Go Lang, ElasticSearch, Redis, and the Go-Gin framework. 

It includes APIs for creating and retrieving shortened URLs. The instructions below guide you on setting up the application and using its APIs.

---

### Steps to Start the Application

1. **Clone the repository**
   ```bash
    git clone https://github.com/beowulf-rohan/go-url-shortner.git
   ```

2. ***Run docker instance***
    ```bash
    cd go-url-shortner
    docker-compose up
   ```

3. ***Shorten API***
    
    Contract -
    ```bash
        url:        mandatory, string   -> url you want to shorten
        short_url:  optional, string    -> short_url is optional, provide a custom short_url if needed
        expiry:     optional, integer   -> short url expiration in days
    ```
    
    Endpoint - 
    ```bash
    curl --location 'http://localhost:3000/shorten' \
        --header 'Content-Type: application/json' \
        --data '{
            "url": "www.google.com",
            "short_url": "",
            "expiry": 12
        }'
    ```

    Response - 
    ```bash
    {
        "url": "https://www.google.com",
        "short_url": "rGu2aeQO",
        "created_at": "2024-11-15T20:08:48.361070715Z",
        "expiry": "2024-11-27T20:08:48.361070756Z"
    }
    ```

4. ***Resolve API***

     Contract -
    ```bash
        short_url:  mandatory, string    -> short_url that was generated using shorten api
    ```

    Endpint - 
    ```bash
    curl --location 'http://localhost:3000/{short_url}'
    example:
    curl --location 'http://localhost:3000/rGu2aeQO'
    
    ```

    Response - 
    ```bash
    Hit the endpoint in a browser and you'll be directed to the orinal url.
    ```

---

