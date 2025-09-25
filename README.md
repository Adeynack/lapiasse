# La Piasse

## Presentation

### The Software

_La Piasse_ is a personal finances manager, the like of _Moneydance_ or _GNU Cash_.

### Justification

This is for now in the realm of the _personal projects_ and aims at pleasing my personal needs (@adeynack) when it comes to personal finance management.

As long time user of _Moneydance_, I (@adeynack) grew tired of some of its ergonomic flaws and wanted the following features to be more to my liking:

- Quick keyboard oriented navigation.
- Easy and speedy import of diverse bank statements (usually available as CSV
  files).

### Origin of the name

The word [_piasse_](https://www.je-parle-quebecois.com/lexique/definition/piasse.html) is Québec-French slang for "dollar". The expression [_faire la piasse_](https://www.japprendslequebecois.com/lexique/piasse.html#:~:text=Faire%20la%20piasse%20%3A%20faire%20un,piasse%2C%20c'est%20s%C3%BBr!) is colloquial for _earning good money_. For this software, the name seemed good as it remains easy to type (`lapiasse`) and it refers to its core function.

## Development Approach

### Step 1: TUI/API Hybrid

The initial development will aim to be a TUI (terminal) application. This decision is made to be able to start using the software as fast as possible (but also because the developer, @adeynack, feels like having fun with TUI development).

However, from the start, the development will be lead by an [Open API](https://www.openapis.org/) definition document. That means that every action will be thought as an API endpoint, even if used only in-memory. That allows some future development to already be eased, but also some flexibility for the user, namely:

- ready to be hosted (internet or intranet) as a service
- allows multiple TUI instances to access the same data in a protected way (see database choice below)
  - the first TUI is the de-facto "server"
  - the next TUIs are connecting to this first one, to make sure the single-tenant database does not get accessed by too many processes
- enables the user to build scripts to perform whatever task they seem fit, by accessing the data through a classic web API

The idea – and those are only brainstorming examples – would be things like this:

```sh
# Start the app in default TUI mode
lapiasse

# Start the app in TUI mode, but exposing the API externally.
lapiasse --api-port=3000

# Start the app in TUI mode, but using a remote API endpoint instead of
# a local database and local business logic.
lapiasse --remote="https://lapiasse.adeynack.net/api"
```

As of for the database, in order to enable a simple-to-deploy local personal application, [SQLite](https://sqlite.org/) will be used. It is a well trusted local database engine and has even been recently [pushed more and more](https://www.youtube.com/watch?v=0rlATWBNvMw) as a [high-volume website alternative](https://www.sqlite.org/whentouse.html) to database servers. In the spirit of _let's cross the bridge when we get to the river_, this is the database this application will be based upon for the time being. If the need for it ever gets big enough, in a parallel universe where this personal project will become huge and people will want to use it in the cloud massively, then we can study a [PosqgreSQL](https://www.postgresql.org/) migration or a dual-database offering.

### Step 2: GUI/Web Hybrid

In order to have this application feel more modern, a GUI will have to be offered. This allows amongst other things better graphical reports to be generated and the mouse to be used more.

The approach will be to develop this GUI in [React](https://react.dev/). This is mainly and pragmatically because the main programmer (@adeynack) knows it already and because it would serve him well in his career to get better at it. This approach allows both local and remote use cases to be dealt together. The local usecase will be presented to the user using [Wails](https://wails.io/) in order to package & serve the GUI elements from the single [Go](https://go.dev/) binary. The same GUI can then be served as a single-page web application, since it is already developped as such.
