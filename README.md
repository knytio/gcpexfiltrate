# gcpexfiltrate

Ransomware data exfiltration to GCP bucket

The tool simulates attacker ransomware activity by copying files from local folder to GCP bucket.


# Validate Security Polices
Organizatins deploy network security tools to monitor DNS traffic and detect attackers C&C servers.

Attackers can download malware or exfiltrate to GCP storage buckets and bypass traditional network defenses such as Firewalls, IPS, DNS monitoring systems

# Running

gcpexfiltrate gcpbucket -h
Exfiltrate data to GCP storage bucket

Usage:
  gcpexfiltrate gcpbucket [flags]

Flags:
  -b, --bucket string    GCP bucket name to upload content
  -f, --folder string    Files in folder to upload
  -h, --help             help for gcpbucket
  -k, --keys string      GCP Keys
  -p, --project string   GCP Project Id
  
  
  Exfiltrate files from "folder" to GCP bucket.
  
  
