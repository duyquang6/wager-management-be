## Dependencies
- Docker
- Go 1.17
- MySQL 8.0.25

## Bootstrap
- Run `chmod +x start.sh` if start.sh script does not have privileged to run
- Run `./start.sh --bootstrap` quick bootstrap app (include build, start docker, migrate schema and start app), it will ready to accept connection to :8080 local
- Run `make docker.local.stop` to cleanup

## For developing
- Get tools for developing: `make install-go-tools`
- Build app docker image: run `make build.docker.image`
- Startup local docker compose: `make docker.local.start`
- Stop local docker compose: `make docker.local.stop`
- Migrate schema database: run `./start.sh --migrate`

## Automate CI CD local
Or you can use skaffold to automate that pipeline
- Install skaffold, helm latest, minikube latest version
- Run `skaffold dev --port-forward`
- Every change to source code, will trigger build, unit-test and deploy locally

## Testing & Coverage
- Run integration test(docker, go1.17 required): `./start.sh --integration`
- Run unittest: `make test.unit`
- Check coverage: `make coverage`
- Clean up report files: `make clean`

## cURL request & response format case by case

### Create Wager endpoint
**total_wager_value violate**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 0,
"odds": 120,
"selling_percentage": 40,
"selling_price": 200
}'
```
Response:
```
{
    "error": "Field:TotalWagerValue Error:This field must be larger than 0"
}
```
**odds violate**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 0,
"selling_percentage": 40,
"selling_price": 200
}'
```
Response:
```
{
    "error": "Field:Odds Error:This field must be larger than 0"
}
```
**selling_percentage violate lower bound**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 120,
"selling_percentage": 0,
"selling_price": 200
}'
```
Response:
```
{
    "error": "Field:SellingPercentage Error:This field must be larger or equal 1"
}
```
**selling_percentage violate upper bound**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 120,
"selling_percentage": 101,
"selling_price": 200
}'
```
```
{
    "error": "Field:SellingPercentage Error:This field must be lesser or equal 100"
}
```
**selling_price violate monetary format**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 120,
"selling_percentage": 101,
"selling_price": 200.255
}'
```
Response
```
{
    "error": "Field:SellingPrice Error:Invalid monetary format, only accept 2 decimal point"
}
```
**selling_price violate total_wager_value * total_wager_value**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 120,
"selling_percentage": 40,
"selling_price": 40
}'
```
Response
```
{
    "error": "field SellingPrice must be larger than TotalWagerValue * SellingPercentage"
}
```
**Success**
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 120,
"selling_percentage": 40,
"selling_price": 200
}'
```
Response
```
{
    "id": 6,
    "total_wager_value": 100,
    "odds": 120,
    "selling_percentage": 40,
    "selling_price": 200,
    "current_selling_price": 200,
    "percentage_sold": null,
    "amount_sold": null,
    "placed_at": 1638097000
}
```
### List wager endpoint:
Default page 1, limit 10
```
curl --location --request GET 'http://localhost:8080/wagers'
```
Response
```
[
    {
        "id": 1,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 40.01,
        "current_selling_price": 40.01,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638096517
    },
    {
        "id": 2,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 200,
        "current_selling_price": 200,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638096519
    },
]
```
```
curl --location --request GET 'http://localhost:8080/wagers?page=2&limit=10’
```
Response
```[
    {
        "id": 11,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 200,
        "current_selling_price": 200,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638097075
    },
    {
        "id": 12,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 200,
        "current_selling_price": 200,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638097075
    },
    {
        "id": 13,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 200,
        "current_selling_price": 200,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638097075
    },
    {
        "id": 14,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 200,
        "current_selling_price": 200,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638097076
    },
    {
        "id": 15,
        "total_wager_value": 100,
        "odds": 120,
        "selling_percentage": 40,
        "selling_price": 200,
        "current_selling_price": 200,
        "percentage_sold": null,
        "amount_sold": null,
        "placed_at": 1638097076
    }
]
```
**Violate limit**
```
curl --location --request GET 'http://localhost:8080/wagers?page=2&limit=0'
```
Response
```
{
    "error": "Field:Limit Error:This field must be larger or equal 1"
}
```
**Violate page**
```
curl --location --request GET 'http://localhost:8080/wagers?page=2&limit=0'
```
Response
```
{
    "error": "Field:Page Error:This field must be larger or equal 1"
}
```
### Create purchase endpoint
**Violate buying price required**
```
curl --location --request POST 'http://localhost:8080/buy/1' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "buying_price": 0
    }'
```
Response
```
{
    "error": "Field:BuyingPrice Error:This field must be larger than 0"
}
```
**Violate buying price larger than current_selling_price**
```
curl --location --request POST 'http://localhost:8080/buy/1' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "buying_price": 500
    }'
```
Response
```
{
    "error": "buying price must be smaller or equal current selling price"
}
```

**Violate Wager ID not found**
```
curl --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "buying_price": 500
}'
```
Response
```
{
    "error": "related wager id 1 not found"
}
```

**Success**
```
curl --location --request POST 'http://localhost:8080/buy/1' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "buying_price": 49.98
    }'
```
Response
```
{
    "id": 1,
    "buying_price": 49.98,
    "wager_id": 1,
    "bought_at": 1638099523
}
```
