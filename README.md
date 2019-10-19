# fs

This is `fs`. What does it stand for? Use your imagination. It's an extremely opinionated file sharing tool.

It's purpose-built for my own file sharing needs, but you can use it too. Basically, it:
- Takes a file either by name or via stdin
- Uploads the file to a Google Cloud Storage bucket
- Optionally shortens the link with [shorty](https://github.com/tjhorner/shorty)
- Spits out a public link so you can share the file

It's meant to be dead-simple. It will read the config file from `$HOME/.config/fs/config.json` and you will need to have a service account JSON file at `$HOME/.config/fs/service_account.json`.

Here is a sample `config.json`:

```json
{
  "projectId": "the-bin-256320", // GCP project ID
  "bucketName": "fs.horner.tj", // GCS bucket name
  "host": "https://fs.horner.tj", // Host the bucket is CNAMEd to
  "shorten": true, // Shorten by default (otherwise you will need to explicitly pass --shorten)
  "shorty": {
    "baseUrl": "https://horner.tj" // Base URL of your shorty instance if you want to use it
  }
}
```

## Usage

Some examples:

```shell
$ fs file.txt
https://horner.tj/iWvqZT

$ echo "hello" | fs --name file.txt
https://horner.tj/a2N95R

$ neofetch --stdout | fs --name neofetch.txt
https://horner.tj/UsaD6H
```

## License

MIT