# Homelab Cluster Dokumentation

Dies ist die zentrale Dokumentation für den lokalen Kubernetes-Cluster auf Basis von Talos Linux.

## Architektur & Hardware
- **OS**: Talos Linux (v1.12.5)
- **Kubernetes Version**: v1.35.2
- **CNI**: Flannel (Standard)
- **Nodes**:
  - `192.168.178.88` - Lenovo ThinkCentre M700 (AMD64) - **Control Plane**
  - `192.168.178.139` - Lenovo ThinkCentre M700 (AMD64) - **Worker**
  - *(Geplant)*: Ein weiterer Node (Lenovo oder Raspberry Pi 4 ARM64)

## Wichtige Konfigurationsdateien
Alle wichtigen Konfigurationsdateien liegen unter `~/private/homelab/talos/`:
- `talosconfig`: Die API-Konfiguration für `talosctl`.
- `kubeconfig`: Der Admin-Schlüssel für den Kubernetes-Zugriff via `kubectl`.
- `controlplane.yaml` / `worker.yaml`: Die generierten Talos-Maschinenkonfigurationen.

## Bekannte Probleme & Workarounds
**macOS Netzwerk-Problem (Tailscale/Firewall):**
Wenn Tailscale aktiv war, blockiert der macOS-Kernel oft native Go-Anwendungen (`talosctl`, `kubectl`) mit einem `no route to host` Fehler, selbst wenn das Gerät erreichbar ist.

*Lösung:* Die Nutzung von Docker umgeht dieses Problem zuverlässig.
Es wurde ein Alias für `kubectl` empfohlen, um nahtlos weiterarbeiten zu können:
```bash
alias hk='docker run --rm -it -v ~/private/homelab/talos:/talos -e KUBECONFIG=/talos/kubeconfig bitnami/kubectl:latest'
```

## Nächste Schritte
1. Hinzufügen des 3. Nodes (Lenovo oder Raspberry Pi).
2. Installation eines verteilten Speichersystems (z.B. Longhorn).
3. Deployment erster Anwendungen (paperless-ngx, Home Assistant).
