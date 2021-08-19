# cms-diy
 try to be a NAS-like CMS

Tested in:

- Windows 10

## MySQL Server Config

- Database name: `cms-diy`
- Access `cms-diy` under `localhost` or `127.0.0.1`

## Script Element Standard for `.go` files

1. package
2. import
3. type
4. const
5. var
6. init()
7. main()
8. func()

> Put **private** `func()` into `.p.go` files

## Project Structure

- clients
- model
  - init
  - file
  - role
  - user
- router
  - init
  - middleware
  - routes
- main.go

> Put **configurable** `const` and `var` into `.config.go` files

> Name **extended** `.go` files with an extra extension, e.g. `.extended.go`
