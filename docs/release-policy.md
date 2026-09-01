# Release policy observation

On 2026-09-01, using the repository administrator's authenticated GitHub API
session, the repository setting endpoint returned:

{"enabled":true,"enforced_by_owner":false}

The setting was enabled with the same administrator session via the documented
repository endpoint, which returned HTTP 204. The existing v0.1.0 release was
then re-read and remains preserved with immutable=false; it was not modified,
deleted, or overwritten. GitHub's immutable-release documentation states that
the policy applies to future releases.

The endpoint requires administrator read access. Therefore a GITHUB_TOKEN 403
in a pull-request job is classified as:

stage=repository-policy-observation
reason=insufficient-administrator-read-credential
unknown_class=external-observability

It is not interpreted as enabled=false. The release workflow uses the ordinary
public release API after publication and fails closed unless the new release
reports immutable=true, its tag ref resolves to the exact annotated tag object
and commit, its asset count is exactly one, and the API asset digest matches the
locally computed digest.
