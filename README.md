# go-github
Github User API in Golang

## Usage
Run the following from terminal:
```
go get github.com/jtejido/go-github
cd <golang-path>/go-github
go build
./go-github -config=<config_path>
```

### Configuration file
The API client reads a yaml config file with the ff properties:

```
max_limit: 10 # maximum number of usernames allowed
user_lifetime: 120 # in seconds, the user's lifetime in cache
listen: 0.0.0.0:8080  # domain it listens to
token: <RANDOM_HASH> 	# session token, optional, can be empty
debug: true
```

When running the binary without **-config** flag, it will attempt to lookup from **GITHUB_API_CONFIG** environment, otherwise an error.

### Usage

Run the following in the terminal:

```
$ ./go-github -config=config.yml
```

### Example

```
$ curl "http://localhost:8080/user?name=jtejido&name=abc"
```

**Result:**

```
[{"id":3063240,"login":"abc","avatar_url":"https://avatars.githubusercontent.com/u/3063240?v=4","url":"https://api.github.com/users/abc","name":"Alastair Blake Campbell","company":"Bibliographic Data Services, Ltd","blog":"","location":"Scotland, UK","email":"","bio":"I'm a Systems Analyst/Programmer. I like C#, Java and C++ primarily. I'm also into Rust! ","public_repos":21,"followers":14,"following":3,"created_at":"2012-12-17T13:22:55Z","type":"User"},{"id":13869015,"login":"jtejido","avatar_url":"https://avatars.githubusercontent.com/u/13869015?v=4","url":"https://api.github.com/users/jtejido","name":"","company":"","blog":"","location":"","email":"","bio":"https://github.com/lucky-se7en","public_repos":60,"followers":5,"following":0,"created_at":"2015-08-19T12:01:20Z","type":"User"}]
```

Open a browser or send a CURL request to the following url:

**http://localhost:8080/user?name=[username1]&name=[username2]**

The list of names and its result will be sorted alphabetically.

### Caching
Each usernames will be cached within the given *item_lifetime* defined in the **config file** (default is 120 seconds).
