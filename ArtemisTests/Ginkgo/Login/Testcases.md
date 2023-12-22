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
- Use caps lock key for letters. [case 11]
- Use lots of characters. [case 12]
- Log in with without entering any username or password. [case 13]

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
|[case1](case1_test.go)|Log in with the correct username and password.|Access is granted.||
|[case2](case2_test.go)|Log in with incorrect username and password.|Access is denied and no errors occur.||
|[case3](case3_test.go)|Log in with incorrect username and correct password.|Access is denied and no errors occur.||
|[case4](case4_test.go)|Log in with correct username and incorrect password.|Access is denied and no errors occur.||
|[case5](case5_test.go)|Log in with only a username (correct).|Access is denied and no errors occur.||
|[case6](case6_test.go)|Log in with only a username (incorrect).|Access is denied and no errors occur.||
|[case7](case7_test.go)|Log in with only a password (correct).|Access is denied and no errors occur.||
|[case8](case8_test.go)|Log in with only a password (incorrect).|Access is denied and no errors occur.||
|[case9](case9_test.go)|Log in with special characters/symbols.|Access is denied and no errors occur.||
|[case10](case10_test.go)|Log in using spaces as characters.|Access is denied and no errors occur.||
|[case11](case11_test.go)|Log in using caps lock for letters.|Access is denied and no errors occur.||
|[case12](case12_test.go)|Log in with lots of characters.|Access is denied and no errors occur.||
|[case13](case13_test.go)|Log in with without entering any username or password.|Access is denied and no errors occur.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of log in, including any configurable parameters and troubleshooting steps. | 