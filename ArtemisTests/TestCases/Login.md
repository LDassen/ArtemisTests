# Starting situation
- Broker setup contains 3 brokers. 
- User(s) are already created.
- That port 61619 is opened without SSL.

# Happy cases
- Log in with correct username and password. [1]

# Fault cases
- Log in with incorrect username and password. [2]
- Log in with incorrect username and correct password. [3]
- Log in with correct username and incorrect password. [4]
- Log in with only a username (correct). [5]
- Log in with only a username (incorrect). [6]
- Log in with only a password (correct). [7]
- Log in with only a password (incorrect). [8]
- Use strange/unusual characters. [9]
- Use spaces as characters. [10]
- Use caps lock key for letters. [11]
- Use lots of characters. [12]
- Log in with without entering any username or password. [13]

# Unknowns
- Password recovery
- User lockout
- Session login
- MFA
- Cross-browser/platform
- User  (CGI, guest, admin?)
- Network failure
- Caps lock sensitivity

# Test cases:
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Log in with the correct username and password.|Access is granted.||
|2|Log in with incorrect username and password.|Access is denied and no errors occur.||
|3|Log in with incorrect username and correct password.|Access is denied and no errors occur.||
|4|Log in with correct username and incorrect password.|Access is denied and no errors occur.||
|5|Log in with only a username (correct).|Access is denied and no errors occur.||
|6|Log in with only a username (incorrect).|Access is denied and no errors occur.||
|7|Log in with only a password (correct).|Access is denied and no errors occur.||
|8|Log in with only a password (incorrect).|Access is denied and no errors occur.||
|9|Log in with special characters/symbols.|Access is denied and no errors occur.||
|10|Log in using spaces as characters.|Access is denied and no errors occur.||
|11|Log in using caps lock for letters.|Access is denied and no errors occur.||
|12|Log in with lots of characters.|Access is denied and no errors occur.||
|13|Log in with without entering any username or password.|Access is denied and no errors occur.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of log in, including any configurable parameters and troubleshooting steps. | 