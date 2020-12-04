# Airlines

Airlines is a rest api to manage airlines, ticketing and schedules.

## Configuration
Change configuration on config.json for databae and credentials.
```
{
    "application_name": "airlines",
    "debug": true,
    "server": {
        "address": ":9090"
    },
    "context": {
        "timeout": 2
    },
    "database": {
        "host": "localhost",
        "port": "3306",
        "user": "root",
        "pass": "",
        "name": "airlines"
    },
    "credentials": {
        "jwt_signature_key": "your key"
    }
}
```
## Installation

Make sure go installed on your computer.

```bash
go run app/main.go
```

## Endpoint

|Method | Endpoint      | Description|
|-------| ----------- | ----------- |
| GET |/planes      | Get all planes|
| GET |/planes/:id| Get plane by id|
| POST| /planes | Create new plane|
| DELETE|/planes/:id| Delete plane by id|
|.....|....|....|

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)