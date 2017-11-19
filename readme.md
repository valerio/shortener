# Shortener
A simple url-shortener service, backed by redis and using [hashIDs](http://hashids.org/) for shortening.

**It is not meant for production use**, but more as an example on how to implement a url-shortener.

## Building
Make sure you have `dep` installed and run:
```
> dep ensure
```

Fill in the configuration file (`config.json`) and simply build the docker image.

Default configuration is setup to be used locally with docker-compose, simply run:
```
> docker-compose up
```

## API

The shortener offers a simple API:

- `GET /urls/{key}`: returns the full url stored with the key provided
- `POST /urls`: shortens a new url, json body: 
    ```json
    {
        "url": "http://www.example.com"
    }
    ```
- `GET /{key}`: redirects to the full url for the specified key

## License
See the [license](./LICENSE) file for license rights and limitations (MIT).