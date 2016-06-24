# alpine-bench

`alpine-bench` is useful for running Apache Bench in your cluster to gather
latency metrics from your services. It does this by:

1. Running Apache Bench using provided options
1. Parsing the results into a JSON data structure and printing them to stdout

This makes it useful for running from your automation (e.g. CI/CD) and making
decisions based on the results. For example, in a CD pipeline it can be useful
to deploy a new build to `dev`, measure its latency, and if the result is within
some threshold, proceed to `prod` rollout.

## Usage

Nearly all of `ab`'s options can simply be passed in using Docker's normal cmd
mechnaism, e.g.:

```bash
docker run -it --rm -t devth/alpine-bench -n3 http://google.com/
```

## Special options

Options which require reading from a file can be passed in via ENV instead.
Internally the contents will be written to a file and the corresponding `ab`
option will be used.

### Example

```
docker run -e 'POST=post body' -it --rm -t devth/alpine-bench -n3 http://google.com/
```

Sample output:

```json
{
  "connect": {
    "max": 154,
    "mean": 69,
    "median": 66,
    "min": 54,
    "stdDev": 18.60000038147
  },
  "percentiles": {
    "100": 340,
    "50": 258,
    "66": 264,
    "75": 269,
    "80": 283,
    "90": 291,
    "95": 296,
    "98": 340,
    "99": 340
  },
  "processing": {
    "max": 274,
    "mean": 196,
    "median": 196,
    "min": 142,
    "stdDev": 18.60000038147
  },
  "total": {
    "max": 340,
    "mean": 266,
    "median": 258,
    "min": 250,
    "stdDev": 18.799999237061
  },
  "waiting": {
    "max": 274,
    "mean": 196,
    "median": 196,
    "min": 141,
    "stdDev": 18.60000038147
  }
}
```

The output can be parsed by your CI tool to make decisions about whether to
rollout a deployment (or whatever else you might want to do).

### Options

| `ab` Option    | Env Var |
| -------------- | ------- |
| `-p POST-file` | POST    |
| `-u PUT-file`  | PUT     |

## Unsupported options

Options which cause ab to write to a file are not supported.

- `-e csv-file`
- `-g gnuplot-file`

They will work, but the resulting file will be lost inside the Docker container
when it terminates.

## Get a shell in the container

```bash
docker run -it --rm --entrypoint sh devth/alpine-bench
```

## Development

Build the image, obtain your docker machine IP, run an Nginx container to
measure, and run `alpine-bench` against it:

```bash
docker build -t devth/alpine-bench .
export DM_IP=$(dm ip dev)
docker run -p 8080:80 -it --rm nginx:alpine
docker run -it --rm devth/alpine-bench -n3 http://$DM_IP:8080/
```

Or run it outside of Docker:

```bash
POST='post body' go run main.go -n3 http://$DM_IP:8080/
```

In your nginx logs you'll see something like:

```
172.17.0.1 - - [14/Jun/2016:20:20:42 +0000] "GET / HTTP/1.0" 200 612 "-" "ApacheBench/2.3" "-"
172.17.0.1 - - [14/Jun/2016:20:20:42 +0000] "GET / HTTP/1.0" 200 612 "-" "ApacheBench/2.3" "-"
172.17.0.1 - - [14/Jun/2016:20:20:42 +0000] "GET / HTTP/1.0" 200 612 "-" "ApacheBench/2.3" "-"
```

## License

Copyright 2016 Trevor C. Hartman. Distributed under the
[MIT](https://opensource.org/licenses/MIT) license.
