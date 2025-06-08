# WELCOME TO GATOR

## What is Gator?
Gator is cli blog aggregator project tool that I coded to gain experience with sql databases.

## What you'll need
To run the blog aggregator you'll just need the Go toolchain and Postgresql.

## How to setup
### 1. Download Gator:
pull this repository onto your computer and run `go install https://github.com/Xeninon/Gator` in your terminal from the root of the repository.

### 2. Setup a postgres database:
Create a local postgres database server on your computer and then run the up migrations from the sql/schema directory within this repository.

### 3. Setup config:
In your home directory create a .gatorconfig.json with the contents `{"db_url":"<connection string>?sslmode=disable"}` where `<connection string>` is the url to your database, it should look something like this "postgres://postgres:postgres@localhost:5432/gator".

## How do I use Gator?
### 1. Register a user:
Run `Gator register <username>` in your terminal.

If already registered you can login with `Gator login <username>`.

### 2. Add rss feeds:
To add feeds run `Gator addfeed <name> <rssURL>` you can name the feed anything but the url must be the link to a site's rss feed in order to work as intended.

If other users on the computer have already added feeds you can run `Gator feeds` to view them and run `Gator follow <url>` to follow a feed.

You can run `Gator unfollow <url>` to unfollow a feed.

### 3. AGGREGATE!!!
In a seperate terminal window, run `Gator agg <timestring>`, where `<timestring>` is a valid timestring such as "30s", in order to continuously aggregate posts from the feeds you're following into your database.

Now that we're aggregating, you can run `Gator browse <listLength>` in a terminal window that isn't aggregating, it returns a list of the most recent posts from your followed feeds along with their urls for viewing in a browser.
(`<listLength>` isn't a required argument and defaults to 2)


##### Pssst... you can use `Gator reset` to reset the state of your database without dropping the tables

Have fun with the gator!!! :)