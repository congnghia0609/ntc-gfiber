# ntc-gfiber
ntc-gfiber is an example golang http server using Fiber  

## Quick start
- Install and start [MongoDB](https://docs.mongodb.com/manual/installation/)  
```bash
# install library dependencies
#make deps
export GO111MODULE=on
go mod tidy

# build
make build

# start mode development
make run

# clean build
make clean
```


## Call API Post
### Add New Post
```bash
curl -X POST -i 'http://127.0.0.1:8080/post' \
  -H "Content-Type: application/json" \
  --data '{
    "title": "title1",
    "body": "body1"
  }'
```

### Update Post
```bash
curl -X PUT -i 'http://127.0.0.1:8080/post' \
  -H "Content-Type: application/json" \
  --data '{
    "id": 1,
    "title": "title1 update",
    "body": "body1 update"
  }'
```

### Get Post
```bash
# Get a post
curl -X GET -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/post/1'

# Get slide posts
curl -X GET -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/posts?page=1'
```

### Delete Post
```bash
curl -X DELETE -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/post/1'
```


## Call API Tag
### Add New Tag
```bash
curl -X POST -i 'http://127.0.0.1:8080/tag' \
  -H "Content-Type: application/json" \
  --data '{
    "name": "tag1"
  }'
```

### Update Tag
```bash
curl -X PUT -i 'http://127.0.0.1:8080/tag' \
  -H "Content-Type: application/json" \
  --data '{
    "id": "5ff379a2669ad8ac6d1addc1",
    "name": "tag1 update"
  }'
```

### Get Tag
```bash
# Get a tag
curl -X GET -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/tag/5ff379a2669ad8ac6d1addc1'

# Get slide tags
curl -X GET -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/tags?page=1'
```

### Delete Tag
```bash
curl -X DELETE -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/tag/5ff37b2a669ad8ac6d1addda'
```


## License
This code is under the [Apache License v2](https://www.apache.org/licenses/LICENSE-2.0).  
