# Access Console

A small React app that talks to the Go `User_auth_service` API: register, log in, view profile, and complete profile — with a live request log and a JWT "access badge" display.

## Run it

```bash
npm install
npm run dev
```

Then open the URL Vite prints (usually `http://localhost:5173`).

## Before you start clicking

Your Go backend needs to be running first. From the `User_auth_service` folder:

```bash
go run main.go
```

It listens on `:8080` by default, which is what this app points at out of the box (editable in the left panel under "Service address").

Make sure your `.env` file in the Go project has `MYSQL_DSN` and `Secret` set, and that your MySQL `user` table exists — otherwise the server will fail on startup or crash on the first request (see note below).

## Known backend gotcha

`database.Login`, `View_p`, and `Complete_p` in `mysql.go` call `log.Fatal(err)` on any error — including things like a wrong password or "no such user." That kills the whole Go process, not just that request. Worth changing those to return an error instead of fataling, so one bad login doesn't take the server down.

## Build for production

```bash
npm run build
npm run preview
```
