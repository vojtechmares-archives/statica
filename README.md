# statica

Simple CLI tool to deploy static websites to AWS S3 with Cloudflare DNS

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
docker pull -it vojtechmares/statica
```

## Configuration

### Environment variables

statica is configured via environment variables

- `STATICA_AWS_ACCESS_KEY_ID` - AWS Access Key ID
- `STATICA_AWS_SECRET_KEY` - AWS Secret Key
- `STATICA_AWS_REGION` - AWS Region (region in which S3 bucket will be created)
- `STATICA_CF_API_TOKEN` - Cloudflare API Token

statica currently does **not** support configuration files

### Flags

- `bucket-name` - Overrides default bucket name (default bucket name is `domain` argument)
- `bucket-prefix` - Adds prefix in front of bucket name (does not include separator)
- `bucket-suffix` - Adds suffix behind bucket name (does not include separator)

## Usage

```bash
statica <domain> <directory>
```

- `domain` is your domain in Cloudflare, this argument is mandatory
- `directory` is directory from which to deploy files to S3, default is `.` (current working directory)

### Example

```bash
statica statica.example.com
```

Uploads all files from current working directory to AWS S3 and creates DNS CNAME record pointing to the S3 bucket website endpoint.
