This is an example of how to integrate [Twirp](https://github.com/twitchtv/twirp) with SuperTokens.

To start this server, run:
```
go run cmd/server/main.go
```

## Few important points
- The code [here](https://github.com/supertokens/supertokens-golang/blob/master/examples/with-twirp/cmd/server/main.go#L58) does `CORS(supertokens.Middleware(verifySession(twirpServer)))`.
    - The middleware exposes all the APIs for the frontend to work. Like sign up / sign in etc.
    - We use `verifySession` in optional mode (`SessionRequired` is `false`). This means that all API requests will go through session verification, but sessions won't be enforced. If a session exists, the `SessionContainer` object will be available to the API, else not.
- We can fetch the `SessionContainer` object in an API like shown [here](https://github.com/supertokens/supertokens-golang/blob/master/examples/with-twirp/internal/haberdasherserver/random.go#L38). The API must check if this object is `nil` or not. If a session exists, this object will not be `nil`, and then one can retrieve session info like the userId, or any session data. They can even revoke the session if needed.
- We need to add a request interceptor defined [here](https://github.com/supertokens/supertokens-golang/blob/master/examples/with-twirp/internal/interceptor/errorHandler.go). Some functions provided by SuperTokens (like getting session information) might throw an error which will be caught by this interceptor and a `401` reply will be sent to the frontend client.
