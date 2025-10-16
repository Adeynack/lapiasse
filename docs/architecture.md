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
    server@{"label":"HTTP Server", "shape": "cloud"}
    oapi-doc@{"label": "OAPI Specification<br>(YAML)", "shape": "document"}
    oapi-interface@{"label": "OAPI<br>Generated Go interfaces and request / response structs"}
    controllers@{"label": "Controllers<br><br>*Business layer implementing the OAPI generated interfaces*"}
    data-layer@{"label": "Data Layer<br><br>*Implementation details to be decided*"}
    db[("Database")]
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

controllers-->|reads & writes data through|data-layer
data-layer-->|queries|db
```
