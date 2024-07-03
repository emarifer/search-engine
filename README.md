<div align="center">
  
<h1 align="center">Search Engine</h1>

<img src="./assets/img/logo-doc.png" width="55%">

<hr />

<p style="margin-bottom: 16px;">
    A mini Google. Custom web crawler & indexer written in Golang.
</p>

> 🚧 This is a work in progress and therefore you should expect that the
> application may not have all the features at this moment.

<br />
  
![GitHub License](https://img.shields.io/github/license/emarifer/go-echo-templ-htmx) ![Static Badge](https://img.shields.io/badge/Go-%3E=1.18-blue)

</div>

<hr />

## Features 🚀

- [x] **Golang-Powered:** Leverage the performance and safety of one of the
  best languages in the market for backend development.
- [x] **Search engine based on the `Depth-first search (DFS)` algorithm:** [Depth-first search](https://www.geeksforgeeks.org/depth-first-search-or-dfs-for-a-graph/) is an algorithm for traversing or searching tree or graph data structures, as is the case with HTML documents. To avoid processing the same link more than once, a `unique` constraint is used when storing the urls that will be crawled in subsequent cycles.
- [x] **Indexing `full text search`:** It is carried out using a parser/tokenizer that uses the [Snowball](https://github.com/kljensen/snowball) library and implemented an [inverted index](https://www.geeksforgeeks.org/inverted-index/), which is stored in the database, allowing an efficient query of the terms search.
- [x] **SQL Database Integration:** Storing crawled urls and indexing results in a `Postgres` DB, which allows greater scalability and efficiency in searches.
- [x] **Caching of the responses (in `JSON` format) of the searches performed:** The `Fiber` framework provides [middleware](https://docs.gofiber.io/api/middleware/cache) for easy caching of server responses.
- [x] **Using the `Fiber` framework, `A-H/Templ` and `Htmx` libraries::** The use of [Fiber](https://gofiber.io/), [Templ](https://templ.guide/) and [Htmx](https://htmx.org/) greatly speeds up the creation of a simple user interface for minimal search engine administration. Check out some of my other [repositories](https://github.com/emarifer/gofiber-templ-htmx) for more explanations.
- [ ] **Using interfaces in the `services` package:** The architecture follows a typical "onion model" where each layer doesn't know about the layer above it, and each layer is responsible for a specific thing, in this case, the `services` (package) layer, which allows for better separation of responsibilities and `dependency injection`.
- [ ] **Using concurrency in engine-built crawling functions:** Use is made of one of the features in which the Go language shines most: concurrency, to try to speed up the always heavy link crawling tasks. 🚧 This is a work in progress.

<br />

<hr />

## 🖼️ Screenshots:

<div align="center">

###### Admin login screen and dashboard:

<img src="assets/img/screenshot-01.png" width="26%" align="top">&nbsp;&nbsp;<img src="assets/img/screenshot-02.png" width="26%" align="top">


###### Response to a search performed with cURL:

<img src="assets/img/screenshot-03.png" width="55%" align="top">


</div>


---

## 👨‍🚀 Installation and Usage

Before compiling the view templates, you'll need to regenerate the CSS. First, you need to install the dependencies required by `Tailwind CSS` and `daisyUI` (you must have `Node.js` installed on your system) and then run the regeneration of the `main.css` file. To do this, apply the following commands:

```
$ cd tailwind && npm i
$ npm run build-css-prod # `npm run watch-css` regenerate the css in watch mode for development
```

Since we use the PostgreSQL database from a Docker container, it is necessary to have the latter also installed and execute this command in the project folder:

```
$ docker compose up -d
```

These other commands will also be useful to manage the database from its container:

```
$ docker start search-engine # start container
$ docker stop search-engine # stop container
$ docker exec -it search-engine psql -U postgres # (user: postgres, without password)
```

Besides the obvious prerequisite of having Go! on your machine, you must have [Air](https://github.com/air-verse/air) installed for hot reloading when editing code.

>[!TIP]
>***In order to have autocompletion and syntax highlighting in VS Code for the `Templ templating language`, you will have to install the [templ-vscode](https://marketplace.visualstudio.com/items?itemName=a-h.templ) extension (for vim/nvim install this [plugin](https://github.com/joerdav/templ.vim)). To generate the Go code corresponding to these templates you will have to download this [executable binary](https://github.com/a-h/templ/releases/tag/v0.2.476) from Github and place it in the PATH of your system. The command:***

```
$ templ generate # `templ generate --watch` to enable watch mode
```

>[!TIP]
>***This command allows us to regenerate the `.templ` templates and, therefore, is necessary to start the application. This will also allow us to monitor changes to the `.templ` files (if we have the `--watch` flag activated) and compile them as we save them if we make changes to them. Review the documentation on Templ [installation](https://templ.guide/quick-start/installation) and [support](https://templ.guide/commands-and-tools/ide-support/) for your IDE .***

Build for production:

```
$ go build -ldflags="-s -w" -o ./bin/search-engine ./cmd/search-engine/main.go # ./bin/search-engine to run the application / Ctrl + C to stop the application
```

Start the app in development mode:

```
$ air # This compiles the view templates automatically / Ctrl + C to stop the application
```

---

### Happy coding 😀!!