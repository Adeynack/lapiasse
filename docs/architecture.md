# Architecture

## Front-Facing Parts

This diagram explains an approach allowing all user facing interfaces to funnel down to a single interface that can easily be exposed either directly (direct calls in _Go_ code or through the _Wails_ FE/BE channel) or remotely (API exposed by a HTTP server).

```mermaid
flowchart TD

user@{"label": "End User", "shape": "manual-input"}

subgraph "Front-End"
    tui@{"label": "TUI", "shape": "win-pane"}
    gui@{"label": "GUI", "shape": "win-pane"}
    web-app@{"label": "Web App", "shape": "win-pane"}
    mobile-app@{"label": "Mobile App", "shape": "win-pane"}
end

subgraph "Back-End"
    wails@{"label": "Wails"}
    server@{"label":"HTTP Server"}
    oapi-doc@{"label": "OAPI Specification<br>(YAML)", "shape": "document"}
    oapi-interface@{"label": "OAPI<br>Generated Go interfaces and request / response structs"}
    controllers@{"label": "Controllers<br><br>*Business layer implementing the OAPI generated interfaces*"}
    downstream@{"label":"..."}
end

oapi-doc-->|generates|oapi-interface
server-->|exposes|oapi-interface

user-->|interracts with|tui
user-->|interracts with|gui
user-->|interracts with|web-app
user-->|interracts with|mobile-app

tui-->|calls directly|oapi-interface

gui-->|interacts with BE through|wails
wails-->|exposes|oapi-interface

web-app-->|calls remote|server

mobile-app-->|calls remote|server

oapi-interface-->|is implemented by|controllers
controllers-->downstream
```

## Data Layout

This describes the early stages of the development. It might (will most probably) change in the future for a more "industrial" way.

### Overview

```mermaid
flowchart TD
    db@{ shape: database}

    controllers["Controllers"] --> models["GORM Models"]
    models -- CRUD --> db["SQLite Database"]
    models -- migrate schema --> mig["GORM<br>Auto-Migrate"]
    mig --> db
```

### Access

For the MVP / prototyping phase, the data access will be done using the [GORM](https://gorm.io/) Object-Relational framework. Controllers may use those models directly ([Rails](https://rubyonrails.org/)-style) in order to simplify development at first. However, more complex logic will be defined in the `models` package, in order to make them testable and re-usable.

### Schema Migration

When the software reaches a near stable _Version 1_ state, migrations will be done in manual SQL to ensure that any needed data manipulation is perform properly and without some automatic decision from any other library.

Until then, however, data migrations will be done automatically by [GORM](#access), in order to allow rapid development in the MVP phase.
