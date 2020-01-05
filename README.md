# viyuelaeveryday_tw_analysis
 Analysis of interactions with tweets from @viyuelaeveryday

## How to use

1. Install golang.

2. Install dependencies: 
    ```shellsession
    $ go get gonum.org/v1/plot/
    $ go get golang.org/x/oauth2
    $ go get github.com/dghubble/go-twitter/twitter
    ```

3. Register as [Twitter developer](https://developer.twitter.com/) to get the necessary credentials to access their API.

4. Set the credentials as environment variables:
    ```shellsession
    $ export TWITTER_CONSUMER_SECRET="******"
    $ export TWITTER_CONSUMER_KEY="********"
    ```

5. Build
    ```shellsession
    $ go build
    ```

6. Run
    ```shellsession
    ./viyuelaeveryday_tw_analysis
    ```

A graph will be stored in the `interactions.png` file:

[interactions](interactions.png)
