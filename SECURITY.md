# Security Policy

## Supported Versions

CCL is early-stage software. Security fixes are expected to target the latest released version and the current development branch.

Older releases may not receive separate backported fixes unless maintainers explicitly decide that the issue is severe enough and practical to backport.

## Reporting a Vulnerability

Please do not publish exploit details in a public issue.

If GitHub private vulnerability reporting is enabled for this repository, use it.

If private vulnerability reporting is not available, open a minimal public issue asking for a private security contact. Do not include exploit code, private data, or detailed reproduction steps in that issue.

## What to Report

Security reports are useful when they involve:

- arbitrary file writes or reads
- command execution
- path traversal
- unsafe generated code that can expose users to security risk
- parser or generator behavior that can be triggered by untrusted CCL input
- release, CI, or supply-chain risks

## What Usually Belongs in a Normal Issue

Use a normal bug report for:

- incorrect generated code without security impact
- confusing errors
- crashes caused by ordinary invalid input
- missing validation
- documentation mistakes

If you are unsure, report privately first.

## Expectations

Maintainers will try to acknowledge valid security reports promptly, investigate the issue, and publish a fix or mitigation when practical.

Because CCL is a small project, response times may vary.
