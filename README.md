# Keep IT Safe - app for storing your secrets safely

## Destination

App can be used for storing your personal secret information using encryption algorithms with key rotation. Key for storing data changes every data change.

## Data types

You can store

- `Credentials` (login and password)
- Free `text` data
- `Bynary` data
- `Bank Card` data

## About app structure

Application consists of server and client side. They're both written on GO.

`Server side` uses Postgres database, which is free. Server side code have minimum external libraries that checked for vulnereabilities by huge comunity.

`Client side` uses minimalistic `TUI` (Terminal User interface) with hot key. With it you can easly and very quickly add your secret personal data into your own server.

## Running

For run server and client please follow instructions
./cmd/client/README.md
./cmd/server/README.md

## TUI

### Welcome screen

On client start you dive into welcome screen

```TUI
  ╭────────────────────────────────────╮
  │    Welcome to Keep IT Safe!        │
  │                                    │
  │     Please Register or Login       │
  │ to Start keeping your secrets Safe │
  ╰────────────────────────────────────╯
  Build Version: v1.0.1
  Build Time: 2024/05/10 20:13:50



  ? toggle help • esc/q quit
```

If you want to see additial hot keys, please follow info in the footer. Press `?`

```
  ╭────────────────────────────────────╮
  │    Welcome to Keep IT Safe!        │
  │                                    │
  │     Please Register or Login       │
  │ to Start keeping your secrets Safe │
  ╰────────────────────────────────────╯
  Build Version: v1.0.1
  Build Time: 2024/05/10 20:13:50







  l login       ?     toggle help
  r register    esc/q quit
```

As you see you can press `l` for login and `r` for registration.

### Registration / Login

Due registration or login you can enter your credentials for application access. For navigations between form inputs use arros or enter key.

```
  ╭────────────────────────────────────╮
  │    Register                        │
  ╰────────────────────────────────────╯
  > Login
  > Password

  [ Submit ]


  ← back • esc quit
```

```
  ╭────────────────────────────────────╮
  │    Register                        │
  ╰────────────────────────────────────╯
  > NewUser
  > •••••••

  [ Submit ]


  ← back • esc quit

```

After `Submit` (Press enter on `Submit` button) you can see list of inputed secrets.

If you see an error, please check server is started and configured correctly.

```
  Post "http://localhost:8080/register": dial tcp [::1]:8080: connect: connection refused
  ← back • esc quit
```

### Secret List

Before insert items your secret list is empty.

```
     Secrets

    No items

  No items.



    r reload items    q quit
    a add item        ? close help
```

For adding new secret press `a` and choose apropriate type for your new secret

```
  ╭────────────────────────────────────╮
  │    Choose Secret Type              │
  ╰────────────────────────────────────╯

  [Credentials]

  Text

  Binary

  BankCard

  ← back • esc/q quit
```

Use arrows and `enter` key for choosing.

### Add Credentials

```
  ╭────────────────────────────────────╮
  │    Add credentials                 │
  ╰────────────────────────────────────╯
  > Secret Name
  > Login
  > Password

  Submit

  ctrl+c quit • ← back
```

```
  ╭────────────────────────────────────╮
  │    Add credentials                 │
  ╰────────────────────────────────────╯
  > Netflix my mom account
  > marryp@gmail.com
  > ••••••••••••••

  [Submit]

  ctrl+c quit • ← back
```

Fill secret name, login, password and you'll see new secret in `Secret List`

```
     Secrets

    1 item

  │ Netflix my mom account
  │ Credentials


    ↑/k    up             / filter          q quit
    ↓/j    down           r reload items    ? close help
    g/home go to start    a add item
    G/end  go to end
```

### View Credentials

In secret list press `enter` or `v` and will shown view Secret form

```
  ╭────────────────────────────────────╮
  │    Secret "Netflix my mom account" │
  ╰────────────────────────────────────╯


  Login: marryp@gmail.com

  Password: StrongPassword

  ctrl+c/q quit • ← back
```

### Add free Text

```
  > Name

  ┃   1 Prepare your text here.
  ┃   ~
  ┃   ~
  ┃   ~

  Submit
(ctrl+c to quit)
```

Add secret name, type secret text and choose Submit.

```

  > My secret script

  ┃   1 1. Find big tree
  ┃   2 2. Go 5 steps to West
  ┃   3 3. Go 2 steps to North
  ┃   4 4. Forget amout your juwels.
  ┃   ~
  ┃   ~

  Submit

  (ctrl+c to quit)
```

After adding secret it appeares in Secret list

```
     Secrets

    2 items

    My secret script
    Text

  │ Netflix my mom account
  │ Credentials



    ↑/k up • ↓/j down • / filter • q quit • ? more
```

### View free Text Secret

In secret list press `enter` or `v` and will shown view Secret form

```
  ╭────────────────────────────────────╮
  │    Secret :"My secret script"      │
  ╰────────────────────────────────────╯

  ╭──────╮
  │ Text ├─────────────────────────────────────────────────────
  ╰──────╯
  1. Find big tree
  2. Go 5 steps to West
  3. Go 2 steps to North
  4. Forget amout your juwels
                                                       ╭──────╮
  ─────────────────────────────────────────────────────┤ 100% │
                                                       ╰──────╯
  ctrl+c/q quit • ← back
```

### Add free Bytes Secret

Add secret name, type secret binary data and choose Submit.

```
  ╭────────────────────────────────────╮
  │    Add Binary Data                 │
  ╰────────────────────────────────────╯


  > My home key

  ┃   1 AA BB CC DD 00 00 00 00
  ┃   2 DD BB AA 00 01 02 03 04


  Submit

  ctrl+c quit
```

After adding secret it appeares in Secret list

```
     Secrets

    3 items


 │ My home key
 │ Binary

    My secret script
    Text

    Netflix my mom account
    Credentials




    ↑/k up • ↓/j down • / filter • q quit • ? more
```

### View Binary Secret

In secret list press `enter` or `v` and will shown view Secret form

```
  ╭────────────────────────────────────╮
  │    Secret "My home key"            │
  ╰────────────────────────────────────╯

  ╭────────────╮
  │ Binary Hex ├────────────────────────────────────────────────────────────────────
  ╰────────────╯
  aa


                                                                            ╭──────╮
  ──────────────────────────────────────────────────────────────────────────┤ 100% │
                                                                            ╰──────╯
  ctrl+c/q quit • ← back
```

### Add Bank Card Secret

Add secret name, card info and choose Submit.

```
  ╭────────────────────────────────────╮
  │    Add bank card                   │
  ╰────────────────────────────────────╯
  > My Credit Card
  > John Doe
  > 3123 3234 32432 32423 3459
  > 01 30
  > 123

  [Submit]


  ctrl+c quit • ← back
```

After adding secret it appeares in Secret list

```
     Secrets

    4 items

  │ My Credit Card
  │ BankCard

    My home key
    Binary

    My secret script
    Text

    ••

    ↑/k up • ↓/j down • / filter • q quit • ? more
```

### View Bank Card Secret

In secret list press `enter` or `v` and will shown view Secret form

```
  ╭────────────────────────────────────╮
  │    Secret "My Credit Card"         │
  ╰────────────────────────────────────╯


  Number: 3123 3234 32432 32423 3459

  Holder: John Doe

  Valid: 01 30

  ValidationCode: 123



  ctrl+c/q quit • ← back
```
