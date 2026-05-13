---
title: "Πώς έφτιαξα έναν Static Site Generator με Go"
description: "Η ιστορία πίσω από αυτό το project"
date: 2024-07-15
---

# Πώς έφτιαξα έναν Static Site Generator με Go

Ήθελα να καταλάβω πώς λειτουργεί ένας SSG από μέσα.
Αποφάσισα να φτιάξω τον δικό μου από το μηδέν.

## Η αρχιτεκτονική

Το project αποτελείται από δύο μέρη:

1. **TypeScript CLI** — η interface που βλέπει ο χρήστης
2. **Go Engine** — ο πυρήνας που κάνει τη δουλειά

Ο λόγος για αυτή τη διαχωρισμένη αρχιτεκτονική είναι ότι η Go
δίνει εξαιρετική performance για file processing, ενώ το TypeScript
ecosystem είναι ιδανικό για CLIs.

## Τεχνικές λεπτομέρειες

Το Go engine χρησιμοποιεί:
- **goldmark** για markdown parsing (CommonMark compliant)
- **html/template** για safe HTML rendering
- **gopkg.in/yaml.v3** για front matter parsing

Σύνολο dependencies: μόλις 2! Αυτό είναι η ομορφιά της Go.
