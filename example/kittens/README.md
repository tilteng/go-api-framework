# example api

```bash
$ curl -X POST -H 'Content-Type: application/json' http://localhost:31337/kittens -d '{ "data": { "attributes": { "name": "Sparky" } } }'
$ curl http://localhost:31337/kittens/<uuid>
```
