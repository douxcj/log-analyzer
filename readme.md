# Go Log Observability & Alerting Engine

A high-performance, real-time log monitoring tool built in Go. This project simulates a production observability agent that "tails" system logs, analyzes error rates using sliding time windows, and triggers alerts when thresholds are breached.



## 🚀 Features

* **Non-blocking Log Tailing:** Uses `bufio.Reader` to stream logs line-by-line, ensuring a constant memory footprint even with multi-gigabyte files.
* **Stateful Analysis:** Implemented logic to distinguish between "noise" (transient errors) and "incidents" (clustered errors) using a 10-second sliding window.
* **SRE-Focused Alerting:** Automated incident reporting based on configurable error thresholds (e.g., 3 errors within 10 seconds).
* **Graceful Recovery:** Designed to wait for file updates without crashing or consuming 100% CPU (busy-waiting).

## 🛠️ Technical Implementation

### Systems Thinking (The "SRE" Way)
Unlike a simple script that reads a file and quits, this engine acts as a **system daemon**:
1.  **File Pointers:** It uses `file.Seek(0, 2)` to jump to the end of logs, mimicking a cold-start of a monitoring agent in a live environment.
2.  **Concurrency-Ready:** The logic is partitioned to allow for future expansion into concurrent "worker" patterns.
3.  **Signal vs. Noise:** It specifically filters for `500` status codes and `ERROR` strings to minimize alert fatigue.

## 📋 Getting Started

### Prerequisites
* Go 1.21+
* Git

### Installation
```bash
git clone [https://github.com/douxcj/log-analyzer.git](https://github.com/douxcj/log-analyzer.git)
cd log-analyzer
go mod init [github.com/douxcj/log-analyzer](https://github.com/douxcj/log-analyzer)