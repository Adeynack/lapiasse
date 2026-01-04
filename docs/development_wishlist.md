# Development Wishlist

A parking for ideas and TODOs for the development.

## Roadmap for MVP

> `-` N/A `?` To do

| Feature                   | Model | API  | MD Import   | TUI |
| ------------------------- | ----- | ---- | ----------- | --- |
| Books                     | Done  | Done | In progress | -   |
| Categories                | ?     | ?    | ?           | -   |
| TUI experimentations (\*) | -     | -    | -           | ?   |
| Books                     | -     | -    | -           | ?   |
| Categories                | -     | -    | -           | ?   |
| Accounts                  | ?     | ?    | ?           | ?   |
| Exchanges                 | ?     | ?    | ?           | ?   |
| Reminders                 | ?     | ?    | ?           | ?   |

### (\*) TUI experimentations

I will ask Claude (AI) to generate 4 PoV TUIs using different lib/frameworks (BubbleTea, tview, and 2 other recommendations). To have a decent idea of how the code looks in all of them, Claude will be asked to develop the following screens and controls:

- [ ] Books
  - [ ] List of books
    - [ ] A table with the books and their default currency
    - [ ] Controls to add, remove, and edit books
  - [ ] Book form (create, edit)
- [ ] Categories
  - [ ] Tree view listing categories and their parents
  - [ ] Controls to add, remove, and edit books
- [ ] Category form
  - [ ] Chosing a parent category. Could be one of...
    - Dropdown?
    - Drop-TreeView ?
    - Popup/Modal with a Treeview?

## TODOs for product

- [x] Minimal API security
  - [x] Starting the API logs a key (and/or display it on the screen)
  - [ ] MoneyDance Import needs that key as a parameter to provide to the API calls
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
