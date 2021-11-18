<p align="center">
  <h1 align="center">statica</h1>
  <p align="center">Simple CLI tool to deploy static websites to AWS S3 with Cloudflare DNS</p>
  <p align="center">
    <a href="https://github.com/vojtechmares/statica/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/vojtechmares/statica.svg?style=for-the-badge"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
    <a href="https://github.com/vojtechmares/statica/actions?workflow=Build"><img alt="GitHub Actions" src="https://img.shields.io/github/workflow/status/vojtechmares/statica/Build?style=for-the-badge"></a>
    <a href="https://goreportcard.com/report/github.com/vojtechmares/statica"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/vojtechmares/statica?style=for-the-badge"></a>
    <a href="https://github.com/goreleaser"><img alt="Powered By: GoReleaser" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge"></a>
  </p>
</p>

## Features

**Backendless deploy tool for static websites** to AWS S3 bucket with Cloudflare DNS

- Automatically creates and configures bucket if does not exist
- Automatically creates Cloudflare DNS record to AWS S3 website endpoint
- Uploads static files from current working directory (or give via second argument) to AWS S3 bucket

## Installation

### Homebrew

```bash
brew install vojtechmares/tap/statica
```

Supports macOS and linux.

### Docker

```bash
docker pull vojtechmares/statica
```

## Configuration

### Environment variables

statica is configured via environment variables

- `STATICA_AWS_ACCESS_KEY_ID` - AWS Access Key ID
- `STATICA_AWS_SECRET_KEY` - AWS Secret Key
- `STATICA_AWS_REGION` - AWS Region (region in which S3 bucket will be created)
- `STATICA_CF_API_TOKEN` - Cloudflare API Token

statica currently does **not** support configuration files

## Usage

### Example

```bash
statica example.com dist
```

- `domain` is your domain in Cloudflare, this argument is mandatory
- `directory` is directory from which to deploy files to S3, default is `.` (current working directory)

### Commands

#### statica

Deploys content from directory

Requires at least one argument (domain)

```bash
statica example.com
```

Second argument is optional and specifies source directory of files to upload

Default: `.` (current working directory)

##### Flags

- `bucket-name` - Overrides default bucket name (default bucket name is `domain` argument)
- `bucket-prefix` - Adds prefix in front of bucket name (does not include separator)
- `bucket-suffix` - Adds suffix behind bucket name (does not include separator)
- `no-dns` - Omits DNS record creation

#### statica destroy

Deletes AWS s3 bucket (including content) and Cloudflare DNS record

Requires exactly one argument (domain)

```bash
statica destroy example.com
```

##### Flags

- `bucket-name` - Overrides default bucket name (default bucket name is `domain` argument)
- `bucket-prefix` - Adds prefix in front of bucket name (does not include separator)
- `bucket-suffix` - Adds suffix behind bucket name (does not include separator)
- `no-dns` - Omits DNS record deletion

#### statica version

Prints version of statica.
