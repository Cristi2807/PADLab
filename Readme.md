# In order to run my image, type the command in terminal:
> docker compose up

All the images will be pulled from DockerHub, unless they are present locally. 

My docker-compose.yml is mapping local PORT 5000 to the container port 5000. My Gateway can be reached at 
http://localhost:5000.

List of endpoints:

##### GET /status -> returns HTTP 200

##### GET /shoes -> returns all shoes available in catalog

##### GET /shoes/:Id -> returns shoes with given Id (if found)

##### POST /shoes -> payload ex:
> {"color":"red",
"size":"38",
"price": "123.5", 
"brand": "Gucci",
"category": "casual", 
"model": "ab-46"}

and returns 
> {
    "id": "8d103d46-5905-43d7-a7e3-dfedd8054546",
    "color": "red",
    "size": "38",
    "price": "123.5",
    "brand": "Gucci",
    "category": "casual",
    "model": "ab-46"
}

##### PUT /shoes/:Id -> payload ex. to modify shoes 
> {"color":"green",
"size":"40",
"price": "124.5", 
"brand": "Nike",
"category": "sport", 
"model": "c-33"}

and returns HTTP 200, and an integer in the body, showing nr. of changed rows in the database.

##### GET /transaction/:Id -> returns all transactions for shoes with given Id

##### GET /stock/:Id -> returns how many shoes with given Id are in stock right now

##### GET /turnaround/:Id/:opType -> returns turnaround for shoes with given Id and for given opType (1 or -1)

##### GET /turnaround/:Id/:opType/:since/:until -> returns turnaround for shoes with given Id and for given opType (1 or -1) in period from since and until (since/until of format 2023-10-10 12:34:56+00)

##### POST /transaction -> creates a transaction (if possible, for example retrieving more than available in stock will not be possible) payload ex:
> {
    "shoesId": "30052dda-bf72-48db-a51a-d148e106e473",
    "quantity": "2",
    "operationType": 1
}

and returns transaction with Id and CreationDate
> {
    "id": "a9dddf34-c161-4904-9b8c-20673fe1d347",
    "shoesId": "30052dda-bf72-48db-a51a-d148e106e473",
    "creationDate": "2023-10-18T15:09:53.4287750Z",
    "quantity": "2",
    "operationType": 1
}
