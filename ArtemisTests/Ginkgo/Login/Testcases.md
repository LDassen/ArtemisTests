# Starting situation
- Broker setup contains 3 brokers. 
- User(s) are already created.
- That port 61619 is opened without SSL.

# Happy cases
- Log in with correct username and password. [case 1]

# Fault cases
- Log in with incorrect username and password. [case 2]
- Log in with incorrect username and correct password. [case 3]
- Log in with correct username and incorrect password. [case 4]
- Log in with only a username (correct). [case 5]
- Log in with only a username (incorrect). [case 6]
- Log in with only a password (correct). [case 7]
- Log in with only a password (incorrect). [case 8]
- Use strange/unusual characters. [case 9]
- Use spaces as characters. [case 10]
- Use caps lock key for some characters. [case 11]
- Use caps lock key for all characters. [case 12]
- Use lots of characters. [case 13]
- Log in with without entering any username or password. [case 14]

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
|[case1](case1_test.go)|Log in with the correct username and password and send a message to a queue.|Access is granted and message is sent to queue.||
|[case2](case2_test.go)|Log in with incorrect username and password and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case3](case3_test.go)|Log in with incorrect username and correct password and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case4](case4_test.go)|Log in with correct username and incorrect password and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case5](case5_test.go)|Log in with only a username (correct) and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case6](case6_test.go)|Log in with only a username (incorrect) and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case7](case7_test.go)|Log in with only a password (correct) and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case8](case8_test.go)|Log in with only a password (incorrect) and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case9](case9_test.go)|Log in with special characters/symbols and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case10](case10_test.go)|Log in using spaces as characters and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case11](case11_test.go)|Log in using caps lock for some characters and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case12](case12_test.go)|Log in using caps lock for all characters and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case13](case13_test.go)|Log in with lots of characters and send a message to a queue.|Access is denied and no message is sent to queue.||
|[case14](case14_test.go)|Log in with without entering any username or password and send a message to a queue.|Access is denied and no message is sent to queue.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of log in, including any configurable parameters and troubleshooting steps. | 