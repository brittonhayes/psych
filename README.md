# 👩‍⚕️ Psych - Find a mental health professional

Meet **Psych**, a Go application that allows you to find therapists from psychologytoday.com using a more powerful search engine. This tool provides various functionalities, including pulling therapists from the web, browsing therapists in the terminal, and running a GraphQL playground to more effectively search for a therapist that meets your needs.

<img src="https://raw.githubusercontent.com/ashleymcnamara/gophers/master/GOPHER_SHARE.png" alt="drawing" width="300"/>

## Features

- **Fetching:** Retrieve therapist information from psychologytoday.com based on various criteria such as state, county, city, or zip code.
- **Browsing:** View therapist information in the terminal in a user-friendly interface.
- **GraphQL Playground:** Run a GraphQL server to query therapist data programmatically.

## Installation

### Install with Go

```bash
go install github.com/brittonhayes/psych@latest
```

### Run with Docker

```bash
docker run -p 8080:8080 ghcr.io/brittonhayes/psych -- fetch --state <state> --county <county> --zip <zip> --view
```

## Usage

Psych provides a set of commands to perform various tasks. Here's a brief overview:

### Fetch 

Use the `fetch` command to retrieve therapist from the web

```bash
# Retrieve all therapists in the United States in your county
psych fetch --state <state> --county <county>

# Retrieve all therapists in your zip code
psych fetch --zip <zip>

# Retrieve all therapists in your city
psych fetch --city <city> --state <state>
```

Replace `<state>`, `<county>`, `<city>`, and `<zip>` with the desired criteria for searching therapists.

### Browse

Browse therapists in the terminal using the `view` command.

```bash
psych view
```

### GraphQL Playground 

Run a GraphQL playground to query therapist data using the `view -w` command.

```bash
psych view --port <port> -w
```

Example GraphQL query

```graphql
{
  therapists(filter: { credentials: "LMFT" }) {
    title
    accepting_appointments
    credentials
    statement
    link
  }
}
```

Replace `<port>` with the desired port number for the GraphQL server.

### Additional Flags

- Use `-v` or `--verbose` to enable verbose logging.
- Use `-c` or `--config` to specify the configuration directory path.
- Use `--db` to specify the path to the SQLite DB file.

## Configuration

Psych allows you to customize its behavior using command-line flags. You can also modify the application's source code to further customize its behavior according to your needs.

## License

This project is licensed under the [MIT License](LICENSE).

---

Note:** This README provides a high-level overview of the `Psych` tool. For detailed usage instructions and examples, refer to the application's help documentation by running `psych --help` and `psych <command> --help`.

*This project is not affiliated with or endorsed by psychologytoday.com.*