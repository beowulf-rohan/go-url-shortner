# URL Shortener Application

This application provides a Dockerized URL Shortener service with APIs for creating and retrieving shortened URLs. 
The instructions below guide you on setting up the application and using its APIs.

---

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
    ```bash
    curl --location 'http://localhost:3000/shorten' \
        --header 'Content-Type: application/json' \
        --data '{
            "url": "www.facebook.com",
            "shortened_url": "",
            "expiry": 12
        }'
    ```

4. ***Resolve API***
    ```bash
    curl --location 'http://localhost:3000/{shortURL}'
    ```