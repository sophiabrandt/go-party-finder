<!-- PROJECT SHIELDS -->

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.][status-shield]][status-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/sophiabrandt/go-party-finder">
    <img src="logo.png" alt="Logo">
  </a>

  <h3 align="center">Go Party Finder</h3>

  <p align="center">
    Find other players for your D&D game (toy app)
    <br />
    <a href="https://github.com/sophiabrandt/go-party-finder"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/sophiabrandt/go-party-finder">View Demo</a>
    ·
    <a href="https://github.com/sophiabrandt/go-party-finder/issues">Report Bug</a>
    ·
    <a href="https://github.com/sophiabrandt/go-party-finder/issues">Request Feature</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [About the Project](#about-the-project)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgements](#acknowledgements)

<!-- ABOUT THE PROJECT -->

## About The Project

I've started learning Go, and I'm using this project to solidify my learnings.

**Go Party Finder** is a toy application to find a group of folks who want to play D&D or another pen and paper role-playing game (a “party”).

**The project is a work in progress.**

### Built With

- [Go](https://golang.org/)
- [Docker & docker-compose](https://golang.org/)

<!-- GETTING STARTED -->

## Getting Started

To get a local copy up and running follow these steps.

### Prerequisites

For local development, you should have Docker, docker-compose, and Go >= 1.11 installed.

If you are able to run `make`, you can use the included [Makefile](Makefile).

### Installation

1. Clone the repo

```sh
git clone https://github.com/sophiabrandt/go-party-finder.git
```

2. Spin up the docker containers (in detached mode):

```sh
docker-compose up -d
```

<!-- USAGE EXAMPLES -->

## Usage

1. Run database with Docker:

    ```sh
    docker-compose up
    ```

2. Start application locally:

    ```sh
    go run ./cmd/web
    ```


Navigate to [https://localhost:8000](https://localhost:8000).


## Tests

```sh
go test -v ./...
```

<!-- ROADMAP -->

## Roadmap

See the [open issues](https://github.com/sophiabrandt/go-party-finder/issues) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

Please update tests if necessary.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->

## License

Distributed under the Apache 2.0 License. See [`LICENSE`](LICENSE) for more information.

<!-- CONTACT -->

## Contact

Sophia Brandt - [@hisophiabrandt](https://twitter.com/@hisophiabrandt)

Project Link: [https://github.com/sophiabrandt/go-party-finder](https://github.com/sophiabrandt/go-party-finder)

<!-- ACKNOWLEDGEMENTS -->

## Acknowledgements

- [Let's Go!][letsgo] by Alex Edwards
- [Basic.css](https://github.com/vladocar/Basic.css) by Vladimir Carrer
- [Every Layout](https://every-layout.dev) by Heydon Pickering & Andy Bell
- [ardanlabs/service][service] by ArdanLabs
- [cloud native app](https://github.com/learning-cloud-native-go/myapp) by Dumindu Madunuwan
- [Building Modern Web Applications with Go (Golang)](https://www.udemy.com/course/building-modern-web-applications-with-go/) by Trevor Sawler

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[status-shield]: https://www.repostatus.org/badges/latest/wip.svg
[status-url]: https://www.repostatus.org/#wip
[letsgo]: https://lets-go.alexedwards.net/
[service]: https://github.com/ardanlabs/service
