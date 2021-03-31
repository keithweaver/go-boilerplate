# User

## Sign In

Features:
* Email on sign in
* Verify
* Tracks user agent
* Tracks general location

### Verification

TODO - Approach (One table or two)
TODO - User Agent can be altered so it's a nice to have

## Walk through of Sign In & Sign Up Flow

This will be a walk through of how the sign in and sign up flow works. The simplest approach to accessing an account is using SSO.

The basic sign up, you will need an API key from Sendgrid. Run your API with:

```
SENDGRID_API_KEY=<INSERT> SENDER_EMAIL=<YOUR_EMAIL> SENDER_NAME="Your Name" go run main.go
```

When API is running, let's start with the sign up.

```bash
curl --location --request POST 'http://localhost:8080/user/signup/' \
--header 'Content-Type: application/json' \
--data-raw '{
 "email": "me@keithweaver.ca",
 "password": "demodemo1"
}'
```

The response is:

```
{
    "message": "error: Your password does not meet requirements."
}
```

That's a pretty generic response. Yes, it is, but you would have the password checking on the frontend. You can see the log:

```
WARNING: 2021/03/15 07:24:59 logging.go:67: {"message": "password is has less than 5 special characters", "error": "error: invalid password", "requestId": "a09bd215-15d0-46fc-8d83-ad3579612f25", "domain": "user", "handlerMethod": "SignUp", "serviceMethod": "SignUp", "clientIP": "::1"}
```