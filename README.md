# CommaCloner - A 3Commas Deal Duplicator
[![Go Report Card](https://goreportcard.com/badge/github.com/jslowik/commacloner?style=flat-square)](https://goreportcard.com/report/github.com/jslowik/commacloner)
[![codecov](https://codecov.io/gh/JSlowik/commacloner/branch/main/graph/badge.svg?token=BKLAXDBDGT)](https://codecov.io/gh/JSlowik/CommaCloner)

## NOTICE 
This is a pre-alpha build that I am still testing and expanding on.  The mappings work paper trading to paper trading, so 
people are welcome to test out the configs and stability.  USE ON YOUR LIVE ACCOUNT AT YOUR OWN RISK
## 

CommaCloner leverages both the [3Commas Deals Websocket and Bots API](https://github.com/3commas-io/3commas-official-api-docs) 
in order to use built-in deal start conditions across exchanges.

By leveraging CommaCloner, you can set up a bot with any deal start conditions you like on Paper Trading, and open a 
deal on any live bot across exchanges.

## Getting Started
In order to utilize CommaCloner, you must first generate an [API key and secret](https://3commas.io/api_access_tokens) 
from 3commas.  When generating your API key and secret, ensure that you click the "Bots Read" and "Bots Write" 
checkboxes.

## Configuration
Configurations are loaded into CommaCloner via YAML.  [examples/config.yaml](examples/config.yaml) contains a basic 
template for setting up the application

When creating bots in 3commas, your source bot will have all the deal start conditions.  Once configured, copy your 
bot a second time into the exchange account of your choosing, and change the "Deal Start Condition" to "Manually/API".

Note that source and destination bots can be configured any way you wish, they don't HAVE to match in terms of base/safety
order sizes, max safety trades, etc. All that matters is that the destination bot has the same pairs available as the 
source bot (with exceptions, see "Overrides"). 

## Startup
Use the following command to startup
```bash
./commacloner serve examples/config.yaml
```


#### About Overrides
Overrides allow you to manipulate deals "on the fly" to account for different currencies (USD, USDT, USDC, etc), before
creating the deal on your destination bot.  Additionally, overrides also allow you to cancel deals on your source bot 
that are unavailable on the destination bot. 

NOTE:  
- `panicSellUnavailableDeals` is an extension of `cancelUnavailableDeals`.  if `cancelUnavailableDeals` is false, but 
  `panicSellUnavailableDeals` is true, the deal will NOT be cancelled or panic sold on the source bot.

#### Example Configuration
```yaml
# Options for controlling the logger.
logging:
  #logging level
  level: "debug"
  # "console" or "json" are the valid formats
  format: "console"
# The 3Commas API key and secret.
# NOTE you should not need to touch the websocket_url or rest_url.  This are only left as configuration items in the off
# chance 3commas changes their api endpoint
api:
  key: "qwertyu"
  secret: "asdfghjkl"
  websocket_url: "wss://ws.3commas.io/websocket"
  rest_url: "https://api.3commas.io/public/api"
#bot configurations
# this can be an array of 1 to n configurations.  there is no limit
bots:
  -
    #just a generic id for personal organization
    id: my_first_mapping
    source:
      #the id of the bot deals will be listened FROM
      bot_id: 1234
    dest:
      #the id of the bot deals will be sent TO
      bot_id: 5678
    overrides:
      quote_currency: "USD"
      base_currency: ""
      cancelUnavailableDeals: true
      panicSellUnavailableDeals: false
  -
    id: additional_bot
    source:
      bot_id: 4285
    dest:
      bot_id: 8675309
    overrides:
      quote_currency: ""
      base_currency: "USDC"
      cancelUnavailableDeals: true
      panicSellUnavailableDeals: true
```


## Getting help
- For feature requests and bugs, file an [issue](https://github.com/jslowik/CommaCloner/issues).

## Tip Jar
- BTC `3Ec8nc2kxKSZukwokG3L3s23EyiswWuQGQ`
- ETH `0xf1AD69127AE25f84F660B3A6C6cDdcd77716484F`
- DOGE `DHzEhGdTZSQQZDGwHmYMRiUR6ES4JPwazG`
- USDT `0x46B626d2A13e3F309F9b88278b4bF5a38a01061A`
