#!/usr/bin/env bash
# health.sh: read-only host snapshot. Safe to run anywhere.
set -u
section() { printf '\n\033[1m== %s ==\033[0m\n' "$1"; }
have()    { command -v "$1" >/dev/null 2>&1; }

section "Host"
echo "hostname : $(hostname)"
echo "os       : $(. /etc/os-release 2>/dev/null && echo "$PRETTY_NAME")"
echo "kernel   : $(uname -r)"
echo "uptime   : $(uptime -p 2>/dev/null || uptime)"
echo "load     : $(cut -d' ' -f1-3 /proc/loadavg)  (cores: $(nproc))"

section "Memory"
free -h 2>/dev/null | sed -n '1,3p'

section "Disk (>80% flagged)"
df -hP -x tmpfs -x devtmpfs -x squashfs 2>/dev/null | awk '
  NR==1 {print; next}
  {use=$5; sub("%","",use); print $0 (use+0>80 ? "   <-- HIGH" : "")}'

section "Top 5 by CPU"
ps -eo pid,user,%cpu,%mem,comm --sort=-%cpu 2>/dev/null | head -6

section "Top 5 by memory"
ps -eo pid,user,%cpu,%mem,comm --sort=-%mem 2>/dev/null | head -6

section "Failed systemd units"
if have systemctl; then
  systemctl --failed --no-legend --no-pager 2>/dev/null | grep . || echo "none"
else
  echo "systemd not available"
fi

section "Listening ports"
if have ss; then
  ss -tulnH 2>/dev/null | awk '{print $1, $5}' | sort -u | head -25
fi

section "Logged-in users"
who 2>/dev/null | grep . || echo "none"

section "Recent failed SSH logins (last 24h)"
if have journalctl; then
  n=$(journalctl -u ssh -u sshd --since "24 hours ago" --no-pager 2>/dev/null \
      | grep -c "Failed password\|Invalid user")
  echo "${n:-0} failed attempts"
fi

section "Pending reboot / updates"
[ -f /var/run/reboot-required ] && echo "reboot required" || echo "no reboot flag"
if have checkupdates; then echo "updates: $(checkupdates 2>/dev/null | wc -l)"
elif have apt;        then echo "updates: $(apt list --upgradable 2>/dev/null | grep -vc Listing)"
fi

section "TLS cert / time sync"
have timedatectl && timedatectl 2>/dev/null | grep -i "synchronized"
