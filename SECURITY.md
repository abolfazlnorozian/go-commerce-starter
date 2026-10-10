# Security policy

## Supported versions

This project is an open-source starter kit. Security fixes are applied on a best-effort basis on the latest `main` branch.

| Version | Supported |
| --- | --- |
| `main` | Yes |
| Older tags / forks | Best effort |

## Reporting a vulnerability

Please report security issues **privately**. Do not open a public GitHub issue for vulnerabilities that could be exploited.

**Preferred contact:** email the maintainer at `abolfazlnorozian7@gmail.com` with:

- A short description of the issue
- Steps to reproduce (PoC if possible)
- Impact assessment (e.g. auth bypass, secret exposure, injection)
- Whether you plan to disclose publicly and on what timeline

You should receive an acknowledgment within a few days when possible. After a fix is available (or the risk is accepted and documented), coordinated disclosure is welcome.

## Secrets and credentials

- Never commit `.env`, production MongoDB URIs, JWT secrets, payment merchant IDs, or SMS API keys.
- If you discover a secret in git history, rotate it immediately and notify the maintainer.
- Use [`.env.example`](./.env.example) as the template for required variables.

## Scope notes

This backend is a modular e-commerce starter, not a hardened multi-tenant SaaS. Known gaps (test coverage, rate limiting, payment idempotency, admin auth consistency) are listed in the README. Reports that improve those areas are appreciated.
