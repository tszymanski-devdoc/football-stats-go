# Local Development with GCP Cloud SQL

## Option 1: Using Cloud SQL Proxy (Recommended)

### Step 1: Install Cloud SQL Proxy
Download from: https://cloud.google.com/sql/docs/postgres/connect-admin-proxy#install

Or via PowerShell:
```powershell
# Download Cloud SQL Proxy v2
Invoke-WebRequest -Uri "https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.15.0/cloud-sql-proxy.x64.exe" -OutFile "cloud-sql-proxy.exe"
```

### Step 2: Authenticate with GCP
```powershell
gcloud auth application-default login
```

### Step 3: Start Cloud SQL Proxy
```powershell
.\cloud-sql-proxy.exe balmy-apogee-484309-h2:europe-west1:football-stats-go-db --port 5432
```

Keep this terminal running!

### Step 4: In a new terminal, run your API
```powershell
cd c:\priv\football-stats-go\football-stats-go-api
go run cmd/api/main.go
```

The `.env` file is already configured to connect via localhost:5432 when using the proxy.

---

## Option 2: Direct Connection via Public IP

If your Cloud SQL instance has a public IP enabled:

1. Get the public IP from GCP Console or:
```powershell
gcloud sql instances describe football-stats-go-db --format="value(ipAddresses[0].ipAddress)"
```

2. Add your IP to authorized networks:
```powershell
gcloud sql instances patch football-stats-go-db --authorized-networks=YOUR_PUBLIC_IP
```

3. Update `.env`:
```env
DATABASE_URL=postgresql://postgres:DvVm$b6&kTMs3g?E@PUBLIC_IP_HERE:5432/postgres?sslmode=require
```

---

## Current Configuration

- **Database Instance:** `football-stats-go-db`
- **Project:** `balmy-apogee-484309-h2`
- **Region:** `europe-west1`
- **Database:** `postgres`
- **User:** `postgres`

## Test Connection

Once Cloud SQL Proxy is running:
```powershell
# Test with psql (if installed)
psql "postgresql://postgres:DvVm`$b6&kTMs3g?E@localhost:5432/postgres?sslmode=disable"

# Or run migrations
go run cmd/migrate/main.go -direction up
```
