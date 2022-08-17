# Authentication in Voogle
## Context
For now, Voogle uses a static HTTP basic auth.

A natural evolution is to support user management (CRUD) and password based authentication. Notice that for performance purpose (i.e. avoiding to verify password for each API request), a cookie based session protocol will also be used aside the authentication protocol.

## API Gateway perimeter
### Description
Voogle authentication will be handle by an API gateway, which is a frontal service dedicated to authentication that acts as a proxy between the world (WAN requests) and the internal Voogle API, that serves the video streaming application.

The choice of implementing authentication only at the API gateway level serves the following purposes:

1. Have a low coupling between business service implementation and authentication. By instance, we can change the authentication protocol without modifying business services.
2. Restrict the authentication security perimeter to the API gateway. The API gateway is the unique service to be allowed to query password verification, and is the unique service that has access to authentication secrets like cookie signing symmetric key and user id derivation key.

### API gateway authentication parameters
The API gateway will be provided two secret keys, via some request to the Vault secret manager, that will be used to handle the authentication protocol:

- `sessionPrivateKey` - the private symmetric key that will be used to crypt and decrypt the session cookies (see [gorilla session](https://pkg.go.dev/github.com/gorilla/sessions))
- `userIdDerivationKey` - a private symmetric key used for user identifier derivation (i.e. derivate IDs used in Vault from the secret key + ID of users within the database). The usage of such derivation mechanism will ensure that the API gateway will be the only actors within the system to be able to associate a password to a publicly identified user.

## Authentication protocol
### The choice of HashiCorp Vault
The chosen solution will rely on [HashiCorp Vault](https://learn.hashicorp.com/vault) for the following reasons : 

1. Explore password management outside of database (i.e. delegate hashing/salting problematics)
3. Use the Vault secret manager to handle the session symmetric key
4. Use the Vault secret manager to handle a derivation key, so as to anonymise user credentials (i.e. derivate password index from a derivation from a secret key and the user ID within the database) 
5. Implement a use case for Squarescale integrated Hashicorp Vault

### Using session cookies
So as not to query password verification for each request, the authentication protocol will use signed sessions,  which will be  store within client cookies.

Consequently, two security parameters be adjusted to customise the protocol: 

1. The session expiration date
2. The cookie signing key rotation frequency

The session will be implemented using the [Gorilla sessions library](https://pkg.go.dev/github.com/gorilla/sessions) so as to ensure consistency with the HTTP router (also Gorilla).

### User registration
![Voogle-authent-register drawio](https://user-images.githubusercontent.com/4182953/185108135-fa9b383f-9677-480c-ba7a-0a6199baf441.png)

1. On the non-authenticated endpoint `/api/register`, a client can request a user registration, by providing a user defined attribute set. These attribute set must contain, at least, credentials information. The credentials informations are composed of a public identifier (let use the client email)  and a password. Password policy will only rely on the length of the password (see [https://auth0.com/blog/dont-pass-on-the-new-nist-password-guidelines/](https://auth0.com/blog/dont-pass-on-the-new-nist-password-guidelines/) for explanation).
2. The user will be created within the database first. **But the password will never be provided to nor stored within the database**. 
3. The password will be stored within the [HashiCorp Vault secret manager](https://www.vaultproject.io/), in a [key/value store secret engine](https://www.vaultproject.io/docs/secrets/kv). The key used to index the password in the store will be derivate from the public user ID (say `email`) + the derivation key `userDerivationKey` hold by the API gateway.
4. The API gateway will return a success code (typically `HTTP 200`) to the web client.

The new user is registered and can know authenticate and generate an authentication session cookie.

### Generate an authentication session
When a user is registered, but has not an authentication session yet, or has an expired authentication session, she will need to authenticate with her password so as to receive an authentication session cookie.

![Voogle-authent-no-session drawio](https://user-images.githubusercontent.com/4182953/185108129-ce019d6a-17ea-42be-8791-002877b88fa4.png)

1. The user requests the endpoint `api/authentication` along with credentials attributes (email + password).
2. Using the derivation technic described above, the API gateway verifies that the password matches the user's email. If it matches, then the protocol continues.
3. The API gateway checks that the user still exists within the database (later we could put more access information within the database and perform some further checking on it). If the user exists, then the protocol continues.
4. A session cookie is generated (let put the user email in it for now) using the [gorilla session library](https://pkg.go.dev/github.com/gorilla/sessions), along with the `sessionPrivateKey` hold by the API Gateway. The generated cookie is sent back to the Web client.

### Authenticating with session
The authentication with cookies will be implemented to most of the API endpoint. It will perform as follows:

![Voogle-authent-session drawio](https://user-images.githubusercontent.com/4182953/185108143-e4516f60-9ad6-465c-a832-96eb1e8c3cfc.png)

1. The cookie is attached to the API request.
2. The session cookie (containing the user email) is checked using the [gorilla session library](https://pkg.go.dev/github.com/gorilla/sessions). The authentication succeeds if we can retrieve some user email within the cookie. ⚠️ Is the request payload or path contains some user ID, we will need to check that the provided identifier (user email) matches.
3. Request is forwarded to Voogle API (or any other internal service).
4. Reply is forwarded to the web client
