# Bookmark management

# Design choices

I chose to use no framework or ORM for this project and do reduce dependencies to the bare minimum like I usually do for microservices that I have to maintain long term, even if it's a little more work. Design choices depend of the context and I would have made different choices for fast prototyping or a monolithic application.
Microservices are small and focussed so we can rapidly have dozens of them to maintain. Dependencies add to the maintenance burden and are rarely consistent across projects (they are at a certain point in time, but not on the long run), thus it's important to contain them.
Go is great for that, since the standard library covers a lot of our needs out of the box. And the standard code is of a much higher quality than any framework code (the best example is how the context is managed is most framework, with concurrent access to a shared state, versus immutability in the standard library)

Plain SQL tend to be easier to maintain in the long run too. Remembering the DSL syntax of the query builder / ORM that we where using 4 years ago is a useless overhead, and every developer knows SQL. A microservice rarely have more than 2 or 3 tables to manage so building queries manually is not a big deal. (once again it's a different story when we build a monolithic application with dozens of entities)

Since it's a go exercice I've implemented the UI in go. This is probably not what I would do on a real project because this is not where Go shines the most. I usually build user interfaces in Elm or Angular. For fast prototyping I usually use Elixir/Phoenix.

# Limitations

This is an exercise and not a real production application so I cut a few corners:

- There is no user management. A real application would require bookmarks to be attached to a user. There is no authentication
- The bookmark's title is not refreshed if changed on the provider (once stored in DB it stays the same)
- The provider list is loaded when the app starts. Needs an app restart to refresh it
- Orphan keywords are not deleted from the DB
- There are a few unit tests but no functional tests. Of course on a real project we'd have some, but provisioning the test environment seems out of scope here, and not really go programming.
- There is room for improvement. I've spread a lot of TODOs in the code. We all know that getting the last 20% correct takes 80% of the time. If you want me to implement one of these missing pieces just let me know.
- I only focussed on the go code, not UI work, so it is very rudimentary. For instance keyword edition requires you to type comma separated values. Of course this is not what we expect from a production app.

# Installation

The application runs in a docker environment but the dev environment is local. It removes the overhead of coding in a container. This is fine with go since it accepts cross compilation. The counterpart is that it requires a local go installation.

To compile and start the containers simply run:

```bash
$ make up
```

This command will build the go binary, all the containers and get everything up and running.

## Database

On first installation you also need to provision the DB. I have not spent a lot of effort here and it requires a local mysql client. Simply run `brew install mysql` on a mac.

Then run this command:

```bash
$ make init
```

If you don't want to install a mysql client locally, just log into mysql using a GUI client like `SequelPro` and execute `database/schema.sql`

Here's the configuration:

- Host: 127.0.0.1
- Username: root
- Password: test
- Port: 3307

## Logs

Logs are available using this command:

```bash
docker-compose logs -f api
```

# API

The API documentation is available here: http://localhost:8080/docs/
(the container must be started)

The API listens on port 8080 and requires basic authentication. The username and password are test:test

Here's an example of bookmark creation for a quick start:

```http
POST /bookmarks HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Basic dGVzdDp0ZXN0
Cache-Control: no-cache

{
	"url": "https://www.flickr.com/photos/adesignstudio/39146026050/in/explore-2018-03-22/",
	"keywords": [
		"awesome",
		"wonderful"
	]
}
```

# Web app

The web application is available here: http://localhost:8080
