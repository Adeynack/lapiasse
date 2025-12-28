# Development Wishlist

A parking for ideas and TODOs for the development.

## Roadmap

| Feature    | Model | API | MD Import | TUI |
| ---------- | ----- | --- | --------- | --- |
| Books      | ✅    | 🛠️  | 🛠️        | 🛑  |
| Categories | 🛑    | 🛑  | 🛑        | 🛑  |
| Accounts   | 🛑    | 🛑  | 🛑        | 🛑  |
| Exchanges  | 🛑    | 🛑  | 🛑        | 🛑  |
| Reminders  | 🛑    | 🛑  | 🛑        | 🛑  |

## TODOs for product

- [ ] Implement pagination
  - [ ] Change to only `?after={pseudo-cursor}` (remove `page` and `page-size`).
  - [ ] Returns `{"pagination": {"has_next": true, "next": "pseudo-cursor"}}` when there is a next page
  - [ ] Returns `{"pagination": {"has_next": false}}` when there is not next page
  - [ ] The server decides the size of pages.
  - [ ] The server does NOT communicate the total number of pages.

## Ideas / Wishlist

- [ ] Try out the [interface approach](https://blog.khanacademy.org/statically-typed-context-in-go/) to Dependency Injection
  - [ ] Beforehand, write a benchmark that integrate as much of the stack as possible
  - [ ] Note results
  - [ ] Perform changes towards _statically typed context_
  - [ ] Benchmark again and observe the gains!
- [ ] Have oapi-codegen generate `OkResponse` instead of `200Response` (e.g.)
- [ ] Use [xdg](https://pkg.go.dev/github.com/adrg/xdg#section-readme) to simplify user config/data directory management.

## Upon Go 1.26 release

- [ ] Replace `error.As` by `error.AsType[T]` ([link](https://antonz.org/accepted/errors-astype/))
