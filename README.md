# Shopify Production Engineer Intern Challenge - Summer 2022
 
The [inventory tracking web application](https://docs.google.com/document/d/1wir0XQuviR6p-uNEUPzsGvMFwqgMsY8sEjGUx74lNrg/edit)
for the Shopify Production Engineer Intern Challenge - Summer 2022. I have
implemented this application in [Go programming language](https://go.dev/) and
have used server side rendering using [Go templates](https://pkg.go.dev/html/template). 
For data store, I used a [SQLite](https://www.sqlite.org/index.html) database 
and [GORM](https://gorm.io/) object relational mapper on top of that.

You can see a running version of application at [TODO]().

## Install

Once you have installed [Git](https://git-scm.com/downloads),
[Go](https://go.dev/doc/install#releases), and [GNU
Make](https://www.gnu.org/software/make/) continue with the
following instructions.

### Build
```shell
git clone https://github.com/shayanh/shopify-challenge-2022
cd shopify-challenge-2022
make build
```

### Run
After building the project, execute the following command to run the app.

```shell
./shopify-challenge-2022
```

Then open [http://127.0.0.1:8000/items](http://127.0.0.1:8000/items) in your 
browser to see the running web app.

## Testing
Run the command below to execute the tests.
```shell
make test
```