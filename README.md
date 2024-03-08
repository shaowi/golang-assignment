<img src="./GIB-logo.png" alt="Group-IB logo" width="300"/>

# Golang Developer Test Assignment: Queue Broker Web Service

## Objective

Your task is to develop a **queue broker** as a web service, adhering to the guidelines and restrictions outlined below.

## Service Endpoints

Your service needs to implement two main functions:

### 1. Adding Messages to a Queue

- **Endpoint**: `PUT /<queue>?v=<message>`
- **Functionality**: Stores a message in the specified queue.
- **Response Codes**:
    - `200 OK` if the operation is successful.
    - `400 Bad Request` if the queue name or the `v` parameter is missing.

**Examples**:

```shell
curl -XPUT http://127.0.0.1
# Response: HTTP/1.1 400 Bad Request

curl -XPUT http://127.0.0.1/color
# Response: HTTP/1.1 400 Bad Request

curl -XPUT http://127.0.0.1/color?v=red 
# Response: HTTP/1.1 200 OK

curl -XPUT http://127.0.0.1/color?v=green 
# Response: HTTP/1.1 200 OK

curl -XPUT http://127.0.0.1/name?v=alex  
# Response: HTTP/1.1 200 OK

curl -XPUT http://127.0.0.1/name?v=anna  
# Response: HTTP/1.1 200 OK
```

### 2. Retrieving Messages from a Queue
- **Endpoint**: `GET /<queue>?timeout=<N>`
- **Functionality**: Fetches a message from the queue, adhering to the FIFO principle, where `N` is the timeout in seconds.
If no message is immediately available in the queue, the service should make the requesting client wait for a message
to arrive or until a specified timeout elapses. Should the timeout period expire without a message becoming available, 
the service must respond with a 404 status code. It's crucial to ensure that messages are delivered to recipients 
**in the same sequence as their requests were received**. In scenarios where multiple recipients are awaiting 
messages with a specified timeout, the message must be dispatched to the recipient whose request was registered first.
- **Response Codes**:
   - `200 OK` with the message in the response body if successful.
   - `404 Not Found` if no message is available within the specified timeout.

**Examples**:

```shell
curl http://127.0.0.1/color -i    
# Response: HTTP/1.1 200 OK
# Body: red

curl http://127.0.0.1/color -i    
# Response: HTTP/1.1 200 OK
# Body: green

curl http://127.0.0.1/color -i    
# Response: HTTP/1.1 404 Not Found

curl http://127.0.0.1/name -i     
# Response: HTTP/1.1 200 OK
# Body: alex

curl http://127.0.0.1/name -i     
# Response: HTTP/1.1 200 OK
# Body: anna

curl http://127.0.0.1/color -i    
# Response: HTTP/1.1 404 Not Found
```

## Restrictions
- **Third-Party Packages**: Use of third-party packages is prohibited, except for standard libraries.
- **Port Configuration**: The port for the service should be configurable via command line arguments.
- **Code Organization**: While not mandatory, it's preferred to consolidate your code into a single .go file 
without additional README files or documentation.
- **Code Conciseness**: Strive for conciseness. Avoid unnecessary "flexibility," logging (except for error handling), 
or debugging features. **The less code, the better**.

## Additional Notes

If any part of the task is not clear - please do not hesitate to contact us with your questions,
we are always happy to help and explain!
