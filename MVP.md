# ✅ Read-Only Controller & View TODOs for MVP Kubernetes Resources

Track progress implementing read-only UI features in the `vigilant` project for common Kubernetes resource types.

---

## 🧱 Pods (`core/v1`)
- [x] List pods by namespace
- [x] Describe pod (status, IP, restarts, node, etc.)
- [x] View pod logs
- [ ] Show container specs (image, ports, limits, env)
- [ ] List related events
- [ ] Auto refresh

---

## 📦 Deployments (`apps/v1`)
- [ ] List deployments
- [ ] Describe deployment (replicas, strategy, selector)
- [ ] Show rollout status
- [ ] List underlying ReplicaSets
- [ ] List related events

---

## 🔁 ReplicaSets (`apps/v1`)
- [ ] List ReplicaSets
- [ ] Describe ReplicaSet
- [ ] List pods owned by the ReplicaSet
- [ ] Show labels and selectors

---

## 🧱 StatefulSets (`apps/v1`)
- [ ] List StatefulSets
- [ ] Describe StatefulSet
- [ ] List controlled pods
- [ ] Show volume claim templates
- [ ] List related events

---

## 🧃 DaemonSets (`apps/v1`)
- [ ] List DaemonSets
- [ ] Describe DaemonSet
- [ ] Show desired/current node counts
- [ ] List pods scheduled by the DaemonSet
- [ ] List related events

---

## 🧪 Jobs (`batch/v1`)
- [ ] List Jobs
- [ ] Describe Job (status, completions, backoff)
- [ ] List pods for the Job
- [ ] Show start/finish time
- [ ] List related events

---

## ⏰ CronJobs (`batch/v1`)
- [ ] List CronJobs
- [ ] Show schedule, last run, next run
- [ ] Describe CronJob (suspend, concurrency policy)
- [ ] List Jobs created by CronJob

---

## 🌐 Services (`core/v1`)
- [ ] List Services
- [ ] Describe Service (type, ports, selectors)
- [ ] Show cluster IP, external IPs
- [ ] List linked endpoints

---

## 🌍 Ingresses (`networking.k8s.io/v1`)
- [ ] List Ingresses
- [ ] Show host-to-service mappings
- [ ] Describe Ingress (rules, TLS)
- [ ] List related events

---

## ⚙️ ConfigMaps (`core/v1`)
- [ ] List ConfigMaps
- [ ] Inspect key-value contents
- [ ] Show labels and annotations

---

## 🔐 Secrets (`core/v1`)
- [ ] List Secrets
- [ ] Show metadata and type (Opaque, TLS, etc.)
- [ ] Show keys and sizes (no decoding)

---

## 🧭 Namespaces (`core/v1`)
- [ ] List Namespaces
- [ ] Show status (Active/Terminating)
- [ ] List resources in namespace
- [ ] List related events

---

## 🖥️ Nodes (`core/v1`)
- [ ] List Nodes
- [ ] Describe Node (conditions, capacity, taints)
- [ ] Show internal/external IPs
- [ ] List pods on node
- [ ] Show allocatable vs used CPU/memory

---

## 💾 PersistentVolumeClaims (PVCs) (`core/v1`)
- [ ] List PVCs
- [ ] Show status (Bound/Pending)
- [ ] Show capacity, access modes
- [ ] Link to bound PersistentVolume

---

## 💽 PersistentVolumes (PVs) (`core/v1`)
- [ ] List PVs
- [ ] Show capacity, access modes, reclaim policy
- [ ] Link to bound PVC
- [ ] Show storage class and backing details

---

## 🧪 Events (`core/v1`)
- [ ] List events in a namespace
- [ ] Show involved object
- [ ] Show reason, message, count, timestamps

---

## 📈 HorizontalPodAutoscalers (HPAs) (`autoscaling/v1`)
- [ ] List HPAs
- [ ] Show target resource
- [ ] Show current vs desired metrics
- [ ] Show scale history (if available)