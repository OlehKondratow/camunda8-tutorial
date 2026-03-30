# Security

## Reporting a vulnerability

Please **do not** open a public GitHub issue for security-sensitive reports.

If you believe you have found a security issue in this repository (e.g. unsafe defaults in sample code, credentials handling in tutorial snippets):

1. Open a [private security advisory](https://github.com/OlehKondratow/camunda8-tutorial/security/advisories/new) on GitHub, **or**
2. Contact the repository owner via the email associated with their GitHub profile.

Include: affected paths, reproduction steps, and impact assessment if possible.

## Scope and expectations

This project is **educational sample code**. It is **not** a hardened production distribution of Camunda 8. Tutorial workers and Compose stacks may use plaintext gRPC, demo credentials, or simplified TLS — by design for local learning. Do not deploy them as-is to the public internet without a security review.

For production deployments, follow [Camunda’s security documentation](https://docs.camunda.io/docs/self-managed/setup/overview/) for your platform and version.
