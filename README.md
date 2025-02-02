# url_shortener

URL Shortener service that use an http server as proxy and gRPC server for main logic.

## Internals

### Shortcode

#### Generator

For a string that contains `n` number of letters and digit, we want to have 36^n unique random values.

This is an important statement for us since we want to generate a unique code for each users URL.

Dictionnaire : `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
