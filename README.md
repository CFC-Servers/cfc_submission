# CFC Suggestions
AWS API for CFC Forms

## Project structure
### Types
##### Form
forms represent a type of form e.g. suggestion or ban appeal
- **validators:**
forms have validators these ensure the data is valid, checking submission fields for things like min and max length
if a submission is not valid the validator should return an error
- **destinations:**
forms have a slice  of destinations these contain a name and an object implementing the Sender interface
destinations handle sending a submission to its final destination, in this case discord

##### Submission
submissions represent a pending submission