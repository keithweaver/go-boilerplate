# Go Boilerplate

## Running

```
go run main.go
```

You will need to update `isValidPassword` in `services/user_service.go`.


## Code Structure

### Why move from handlers, services, repositories in individual folders to domain specific?

For extendability is the short answer.

In the original version of this boilerplate, I had `services/`, `handlers/`, `repositories/` and `models/` all be folders. In each folder it would be `<domain>_<folder>.go`. For example, `handlers/cars_handler.go`. 

This domain specific approach allows for new packages to be added and deleted depending on the package. I can add a package for individual service, and if you do not want it, just delete it.

### Logging & Context


## Database

TODO


## Demo Commands

### Sign up

```
curl --location --request POST 'http://localhost:8080/user/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
	"email": "foobar@demo.com",
	"password": "test123"
}'
```

Response

```
{
    "message": "Signed up",
    "token": "5f2dde6175bd1baf7e9a5806"
}
```

### Sign In

```
curl --location --request POST 'http://localhost:8080/user/signin' \
--header 'Content-Type: application/json' \
--data-raw '{
	"email": "foobar@demo.com",
	"password": "test123"
}'
```

Response

```
{
    "message": "Signed in",
    "token": "5f2ddeb075bd1baf7e9a5807"
}
```

### Log Out

```
curl --location --request POST 'http://localhost:8080/user/logout' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer 5f2e9d3001e3a273ed558c49' \
--data-raw '{}'
```

Response

```
{
    "message": "Logged out"
}
```

### Create Car

```
curl --location --request POST 'http://localhost:8080/cars/' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer 5f2e9f23a9aefb542b3a8e18' \
--data-raw '{
	"make": "Mazda",
	"model": "3"
}'
```

Response

```
{
    "message": "Error: Year is missing"
}
```

```
curl --location --request POST 'http://localhost:8080/cars/' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer 5f2e9f23a9aefb542b3a8e18' \
--data-raw '{
	"make": "Mazda",
	"model": "3",
	"year": 2013
}'
```

Response

```
{
    "message": "Created car"
}
```

### Get Car

```
curl --location --request GET 'http://localhost:8080/cars/5f2ea7d2dd45bb607bc45707' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer 5f2e9f23a9aefb542b3a8e18'
```

Response

```
{
    "car": {
        "id": "5f2ea7d2dd45bb607bc45707",
        "make": "Mazda",
        "model": "3",
        "year": 2013,
        "status": "",
        "email": "foobar@demo.com",
        "created": "2020-08-08T13:25:38.6Z"
    },
    "message": "Car retrieved"
}
```

### Delete Car

```
curl --location --request DELETE 'http://localhost:8080/cars/5f2ea7d2dd45bb607bc45707' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer 5f2e9f23a9aefb542b3a8e18'
```

Response

```
{
    "message": "Deleted car"
}
```

### Update Car

```
curl --location --request PUT 'http://localhost:8080/cars/5f2ea852dd45bb607bc45708' \
--header 'Authorization: Bearer 5f2e9f23a9aefb542b3a8e18' \
--header 'Content-Type: application/json' \
--data-raw '{
	"year": 2018
}'
```

Response

```
{
    "message": "Updated car"
}
```

### List Cars

```
curl --location --request GET 'http://localhost:8080/cars/' \
--header 'Authorization: Bearer 5f2e9f23a9aefb542b3a8e18' \
--header 'Content-Type: application/json'
```

Response

```
{
    "cars": [
        {
            "id": "5f2ea85bdd45bb607bc45709",
            "make": "Landrover",
            "model": "Defender",
            "year": 2019,
            "status": "",
            "email": "foobar@demo.com",
            "created": "2020-08-08T13:27:55.132Z"
        },
        {
            "id": "5f2ea852dd45bb607bc45708",
            "make": "Landrover",
            "model": "Rangerover",
            "year": 2018,
            "status": "",
            "email": "foobar@demo.com",
            "created": "2020-08-08T13:27:46.969Z"
        }
    ],
    "message": "Cars retrieved"
}
```
