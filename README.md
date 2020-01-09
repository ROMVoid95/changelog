# Changelog
A changelog generator for Gitea

[![Build Status](https://drone.gitea.com/api/badges/gitea/changelog/status.svg)](https://drone.gitea.com/gitea/changelog)

## Purpose

This repo is currently part of Gitea. The purpose of it is to generate a changelog when writing release notes.  
This project was made for Gitea, so while you are free to use it in your own projects, it is subject to change with Gitea.  
This tool generates a changelog from PRs based on their milestone and labels.

## Installation

```
go get gitea.com/gitea/changelog
```

## Configuration

See the [changelog.example.yml](changelog.example.yml) example file.

## Usage

#### Changelog Entries
```
changelog -m=1.11.0 -c=/path/to/my_config_file generate
```

#### Contributors List
```
changelog -m=1.11.0 -c=/path/to/my_config_file contributors
```

## Building
```
go generate ./...
go build
```

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Maintainers](https://gitea.com/org/gitea/members)
* [Contributors](https://gitea.com/gitea/changelog/commits/branch/master)<!-- FIXME when contributors page works -->

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
