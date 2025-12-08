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

## Ideas / Wishlist

- [ ] Try out the [interface approach](https://blog.khanacademy.org/statically-typed-context-in-go/) to Dependency Injection
  - [ ] Beforehand, write a benchmark that integrate as much of the stack as possible
  - [ ] Note results
  - [ ] Perform changes towards _statically typed context_
  - [ ] Benchmark again and observe the gains!
- [ ] Have oapi-codegen generate `OkResponse` instead of `200Response` (e.g.)

## Upon Go 1.26 release

- [ ] Replace `error.As` by `error.AsType[T]` ([link](https://antonz.org/accepted/errors-astype/))
