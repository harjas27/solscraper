# solscraper
## Features
* create APIs from on-chain data which is not limited to just the generics (hash, fees etc)
* add dApp specific information to the APIs
* filter data using custom queries

## Design Details
* Scraper: create/save queries which will run as background jobs to fetch and parse on chain data
* Query: specify the contract, filters, and fields you would like to export
* Storage: Centralised storage for now in the form of MongoDB to store query metadata and the results from the queries
