<p align="center">
  <h1 align="center">CESI Dossier de Synthèse</h1>
  <p align="center">
    <img align="center" width="70" src="https://github.com/julien-wff/cesi-dossier-synthese/blob/main/web/static/favicons/favicon-96x96.png?raw=true" alt="Dossier de synthèse logo"/>
  </p>
  <p align="center">
    A web interface to extract and visualize grades from CESI PDF files.
  </p>
  <p align="center">
    <img alt="GitHub Actions Workflow Status" src="https://img.shields.io/github/actions/workflow/status/julien-wff/cesi-dossier-synthese/test-commits.yaml">
    <img alt="GitHub Release" src="https://img.shields.io/github/v/release/julien-wff/cesi-dossier-synthese">
    <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/cefadrom/cesi-dossier-synthese">
    <img alt="Docker Image Size" src="https://img.shields.io/docker/image-size/cefadrom/cesi-dossier-synthese">
    <img alt="GitHub License" src="https://img.shields.io/github/license/julien-wff/cesi-dossier-synthese">
  </p>
</p>

In my engineering school ([CESI](https://www.cesi.fr)), the grades are only available in PDF format.
This makes it difficult to view them properly, calculate averages, and estimate future grades.

This project provides a web application where you can upload your PDF file and extract the grades inside.
It displays the grades in a user-friendly interface, allowing you to change them to simulate future ones, and see basic
statistics.

https://github.com/user-attachments/assets/8ee91f46-0bff-4c5d-a332-f803995d0659

## Privacy

PDFs are never stored; they only stay in the server's RAM for the duration of the request.
Only basic telemetry is collected to improve the project, such as the number of files processed, their size, or the
timing of operations.

If you don't trust this, the project is open source, so you can host it yourself!

## Self-hosting

You can self-host this project using Docker. The following command will run the project on port 8080:

```bash
docker run -d \
           --name cesi-dossier-synthese \
           -p 8080:8080 \
           -v ${PWD}/data:/app/data \
           cefadrom/cesi-dossier-synthese
```

You can also use Docker Compose with the following `compose.yaml` file:

```yaml
services:
  cesi-dossier-synthese:
    image: cefadrom/cesi-dossier-synthese
    container_name: cesi-dossier-synthese
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data # To store basic telemetry logs
```

Here are the configurable environment variables:

| Variable             | Description                                                                     | Default Value    |
|----------------------|---------------------------------------------------------------------------------|------------------|
| `APP_ENV`            | Set environment to production or development                                    | `production`     |
| `APP_PORT`           | Port on which the application will run                                          | `8080`           |
| `TELEMETRY_USER`     | Username to access the telemetry dashboard (`/telemetry`)                       | None (no access) |
| `TELEMETRY_PASSWORD` | Password to access the telemetry dashboard (`/telemetry`)                       | None (no access) |
| `PROXY_HEADERS`      | If the app is behind a reverse proxy, to use headers like `X-Forwarded-For`     | `false`          |
| `RATE_LIMIT_TOKENS`  | Number of parsing requests allowed per minute per IP address (`0` for no limit) | `20`             |

## Affiliations

This project is not associated with [CESI Engineering School](https://www.cesi.fr), but it is affiliated with the
student association [BDE CESI Nancy](https://bdecesinancy.fr), which offers a
[hosted version](https://dossier.bdecesinancy.fr) that can be used for free.
