<p align="center">
  <a href="https://appcrons.netlify.app" target="_blank" style="background-color:blues; width:auto; height:auto; display:flex; justify-content:center; align-items:end; gap:16px;">
    <picture>
      <img alt="Appcrons" src="internal/assets/logo.png" width="64" height="60" style="max-width: 100%;">
    </picture>
    <span style="font-size:48px; color:white; font-weight:bold;">Appcrons<span>
  </a>
</p>

<div style="width:100%; display:flex; justify-content:center; align-items:center;">
<p align="center" style="width:80%; max-width:500px;">
 Appcrons optimizes the uptime of your free backend instance on Render by sending automated requests, preventing it from shutting down due to inactivity.
</p>
</div>

## Documentation

For full documentation, visit [appcrons.netlify.app](https://appcrons.netlify.app).

## Installation

To install the Appcrons backend repo locally, follow these steps:

### Prerequisites

- Ensure you have **Golang 1.23** installed. You can download it from [go.dev/doc/install](https://go.dev/doc/install).
- Ensure you have **Postgresql 15** installed. You can download it from [postgresql.org/download](https://www.postgresql.org/download/).
- Ensure you have **Redis** installed. You can download it from [redis.io](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/).

- Ensure your computer system supports **Makefiles** . If you are windows, follow this guide [Run Makefile on windows](https://medium.com/@samsorrahman/how-to-run-a-makefile-in-windows-b4d115d7c516).

### Steps

1. **Clone the repository**:

   ```sh
   https://github.com/Tibz-Dankan/appcrons.git

   cd appcrons
   ```

1. **Install packages**:

   ```sh
   make install

   ```

1. **Set up environmental variables**:

| Variable            | Type   | Description                                  |
| ------------------- | ------ | -------------------------------------------- |
| `APPCRONS_TEST_DSN` | string | DSN pointing to test postgres db             |
| `APPCRONS_DEV_DSN`  | string | DSN pointing to development postgres db      |
| `APPCRONS_STAG_DSN` | string | DSN pointing to staging or CI/CD postgres db |
| `APPCRONS_PROD_DSN` | string | DSN pointing to production postgres db       |
| `REDIS_URL`         | string | URL pointing to your redis instance          |
| `JWT_SECRET`        | string | key used to sign JWT tokens                  |

create .env file in the root project directory and add all these variables

_Example_

```sh

APPCRONS_DEV_DSN="host=localhost user=postgres password=<db password> dbname=<db name> port=<db port> sslmode=disable"

```

4. **Start the application**:

   ```sh

   make run
   ```

5. **Run tests in devlopment**:

   ```sh

   make test

   ```

6. **Run tests in staging**:

   ```sh

   make stage

   ```

> Note: The application server port is **8080**
