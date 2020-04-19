# Synchronizer module
This module is responsible for synchronizing the local elevator with
the other network connected elevators and events. It became a pretty large
and somewhat complicated module. For acknolwedgement of orders we made 
a matrix of this form: 

Implicit acks | Button Up | Button Down
------------- | --------- | ----------
Floor 4       | :arrow_up: / :hash: / :hash: / :hash: | :arrow_down: / :hash: / :hash: / :hash:
Floor 3       | :arrow_up: / :hash: / :hash: / :hash: | :arrow_down: / :hash: / :hash: / :hash:
Floor 2       | :arrow_up: / :hash: / :hash: / :hash: | :arrow_down: / :hash: / :hash: / :hash:
Floor 1       | :arrow_up: / :hash: / :hash: / :hash: | :arrow_down: / :hash: / :hash: / :hash:

where :hash: represent's each online elevator's ack status. The status
can only change one way, not back and forth. The possible changes are:

`0` -> `1` -> `-1` -> `0`

- `0` - empty 
- `1` - acknowledged
- `-1` - finished

This works with the assumption that all other elevators tells the truth.
So if an elevator says that it has acknowledged an order up in floor `2`,
and you say that order is empty - now you also say that it's acknowledged.
It allows for efficient flow of information. 
