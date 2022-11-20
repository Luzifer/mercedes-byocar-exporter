# Luzifer / mercedes-byocar-exporter

This repository contains an Prometheus & InfluxDB exporter for the Mercedes Benz "Bring Your Own Car" (BYOCAR) API products.

Features:

- Store credentials either in Vault or in a local JSON file
- Fetch data for all cars in your MercedesME account
- Prometheus exporter for the metrics
- InfluxDB exporter avoiding spamming entries to the database by using reported dates from Mercedes API

## Usage

```console
# mercedes-byocar-exporter
Usage of mercedes-byocar-exporter:
      --client-id string          Client-ID of Mercedes Developers Console App
      --client-secret string      Client-Secret of Mercedes Developers Console App
      --credential-file string    Where to store tokens when using client-id from CLI parameters (default "credentials.json")
      --fetch-interval duration   How often to ask the Mercedes API for updates (default 15m0s)
      --influx-export string      Set to url (http[s]://user:pass@host[:port]/database) to enable Influx exporter
      --listen string             Port/IP to listen on (default ":3000")
      --log-level string          Log level (debug, info, warn, error, fatal) (default "info")
      --redirect-url string       Redirect URL registered in Mercedes Developers Console (default "http://127.0.0.1:3000/store-token")
      --vault-key string          Use credentials from and update in Vault
      --vehicle-id strings        Vehicle identification number (e.g. WDB111111ZZZ22222)
      --version                   Prints current version and exits
```

## Setup: Create the Mercedes Developer App

- Go to the [Mercedes Benz Developer Portal](https://developer.mercedes-benz.com/) and log in with your Mercedes ID (the same you've registered your car to in Mercedes ME)
- Create a new project in the console section
- Add these products ("Get free" -> BYOCAR -> Select your Project)
  - [Vehicle Status BYOCAR](https://developer.mercedes-benz.com/products/vehicle_status)
  - [Pay As You Drive Insurance BYOCAR](https://developer.mercedes-benz.com/products/pay_as_you_drive_insurance)
  - [Electric Vehicle Status BYOCAR](https://developer.mercedes-benz.com/products/electric_vehicle_status)
  - [Vehicle Lock Status BYOCAR](https://developer.mercedes-benz.com/products/vehicle_lock_status)
  - [Fuel Status BYOCAR](https://developer.mercedes-benz.com/products/fuel_status)
- Note down **Client ID** and **Client Secret** of your project
- Add the redirect URL you will deploy this exporter to (`https://exporter.example.com/store-token`)

## Setup: Deploy the exporter

You can

- build the Go application by running `go build` in the checkout
- build the Docker container by running `docker build .` in the checkout
- get a pre-built image

When running with local JSON-file as storage you need to specify the `client-id`, `client-secret` and `credential-file` flags or corresponding environment variables (`CLIENT_ID`, `CLIENT_SECRET`, `CREDENTIAL_FILE`).

When running with Vault as storage backend specify the `vault-key` (`VAULT_KEY`), `VAULT_ADDR` and `VAULT_TOKEN` or `VAULT_ROLE_ID` / `VAULT_SECRET_ID` for access to Vault. Inside Vault KV v1 backend store this JSON (set your client-id and secret): `{"client-id": "", "client-secret": ""}` and make sure the process can **write** to that key to store user tokens.

In all cases specify one or more `--vehicle-id` (`VEHICLE_ID=WDB111111ZZZ22222,WDB111111ZZZ22223`) to fetch data for. All of those cars **must** be associated to your Mercedes ID.

## Setup: Authorize exporter

When everything is running you should be able to access the exporter:

- `https://exporter.example.com/auth` - Redirect to authorize your project to access your car(s)
- `https://exporter.example.com/healthz` - Health-Check endpoint
- `https://exporter.example.com/metrics` - Text-version of exported metrics

You need to access the `/auth` route once to fetch access- and refresh-keys. If something wents wrong with those keys you can re-authorize the app using this route.

## Setup: Security

⚠️ This exporter does **not** have any security measures like access control and will never have them!

I strongly advice to put the exporter behind auth or any non-public network and ensure no unauthorized user can access any of the endpoints:

- The `/auth` endpoint can be used to mess with the authorization (even though this makes no sense as it will just replace the credentials)
- The `/metrics` endpoint will expose your VIN/FIN to anyone accessing it
