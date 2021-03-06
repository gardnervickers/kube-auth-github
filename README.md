# github-kube-auth
Authenticate users to Kubernetes based on Github team membership.

## Usage
You can run the auth server in Docker by using the following command:

```
docker run --publish 8989:8989 --env GITHUB_ORG=CommercialTribe --env GITHUB_TEAM=Engineering CommercialTribe/github-kube-auth
```

### Configuration

The following environment variables are supported for configuration:

ENV Var | Description | Default
---:|---|---
`GITHUB_ORG`  | The GitHub organization of the team you wish to authenticate users against.
`GITHUB_TEAM` | The GitHub team that you wish to authenticate users against.
`PORT`        | The port to start the HTTP server on. | `8989`

### End User
The user should [create a personal access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use) with the `org:read` and `user:email` scopes. The user should then use the following command to setup kubectl.

```
kubectl config set-credentials [github-username] --token=[access-token]
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.
