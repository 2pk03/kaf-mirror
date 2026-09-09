# Experimental ransomware protection

This is an **experiment**, not a guarantee. It is how we are trying to mitigate a failure we already saw in production.

**What happened:** the **source** cluster was ransomed. kaf-mirror did what it is built to do and **mirrored everything**, including the attack, onto the live replica. Data **already** written through kafscale to S3 could not be encrypted: those objects are read-only. New apply would still have written ransom payload as new objects.

The switch does not make Kafka immutable. It **stops further apply** so the replica and new S3 objects are not filled with poison. The salvage is the already-mirrored kafscale/S3 copy plus a frozen replica.

## Enable (admin)

Protection is off until an admin turns it on.

```bash
mirror-cli protection enable
# or POST /api/v1/protection/enable
```

YAML `protection.enabled: true` also enables it on process start.

Once enabled:

- Jobs **into** a cluster with `role: prod` are refused.
- Halt tripwires are armed.

Label clusters when you add them: `prod`, `dr`, or `other` (`role` on the cluster record).

## Halt

Any of these halt **all running jobs** (pause, do not auto-restart):

1. `mirror-cli protection halt "source compromised"`
2. `POST /api/v1/protection/halt` with `{"reason":"..."}`
3. Create the halt file (`data/HALT` by default, `protection.halt_file`)
4. Set `KAF_MIRROR_HALT=1` (`protection.halt_env`)
5. Auto-halt (experimental): tombstone storm, produce/auth error burst, or payload inflation while source lag collapses

Resume (`protection resume`) clears the API/DB halt only. **Remove the file and unset the env** or the process stays halted. Jobs do not restart by themselves; start them after you trust the source.

Out-of-band file/env still work if the API is unusable.

## What this does not do

- It does not encrypt-proof Kafka. kafscale/S3 already-written objects were safe because they are read-only.
- It does not delay apply by N minutes (no buffer).
- Cross-cluster high watermarks are **not** a proof that source and target match. Verify uses **source consumer-group lag** and topic presence.

## Auto-halt

Off unless protection is enabled. Heuristics are coarse. Prefer the admin halt or the HALT file when you know source is on fire.
