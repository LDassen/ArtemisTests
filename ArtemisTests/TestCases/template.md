# Test environment / Starting situation
Test cluster with Artemis deployment, Prometheus deployment, Grafana deployment, ElastisSearch deployment, .... (a bit more detail is preferred).

# Happy cases
- Def: Perfect situation (expected behavior for deployment).
- Example: Log in in with correct username and password. [*Test case link*]

# Fault cases
- Def: Wrong/unwanted situation coming from a correct deployment (unexpected behavior for deployment).
- Example: Log in in with incorrect username and password. [2]

# Unknowns
*Meaning: things you are not certain about, because it might or might not be implemented*
- Example: user lockout
- Example: MFA

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Example: Log in with the correct username and password.|Example: Access is granted.|Example: Acces is indeed granted.|
|2|Log in with incorrect username and password.|Access is denied.||
|3|...|...|...|

# Documentation review