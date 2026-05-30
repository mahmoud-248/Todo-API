# 1. Load and parse the .env file safely
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        # Matches key=value while ignoring lines starting with '#'
        if ($_ -match '^([^#][^=]+)=(.+)$') {
            $key = $matches[1].Trim()
            # Trim surrounding single or double quotes, and trailing whitespace
            $value = $matches[2].Trim().Trim('"', "'")
            
            Set-Item -Path "env:$key" -Value $value
        }
    }
} else {
    Write-Warning "No .env file found in the current directory."
}

# 2. Extract arguments
$command = $args[0]
$name = $args[1]

# 3. Route commands
switch ($command) {
    "up" { 
        if (-not $env:DATABASE_URL) { Write-Error "DATABASE_URL is not set."; exit 1 }
        migrate -path migrations -database $env:DATABASE_URL up 
    }
    "down" { 
        if (-not $env:DATABASE_URL) { Write-Error "DATABASE_URL is not set."; exit 1 }
        $count = if ($name) { $name } else { "1" }
        
        $confirm = Read-Host "Rolling back $count migration(s). Continue? [y/N]"
        if ($confirm -eq 'y' -or $confirm -eq 'Y') {
            migrate -path migrations -database $env:DATABASE_URL down $count
        } else {
            Write-Host "Migration rollback canceled." -ForegroundColor Yellow
        }
    }
    "create" { 
        if (-not $name) { 
            Write-Error "A name is required to create a migration. Usage: .\migrate.ps1 create <migration_name>"
            exit 1 
        }
        migrate create -ext sql -dir migrations -seq $name 
    }
    "force" { 
        if (-not $env:DATABASE_URL) { Write-Error "DATABASE_URL is not set."; exit 1 }
        if (-not $name) { 
            Write-Error "A version number is required for force. Usage: .\migrate.ps1 force <version>"
            exit 1 
        }
        migrate -path migrations -database $env:DATABASE_URL force $name 
    }
    Default {
        Write-Host "Usage: .\migrate.ps1 [command] [arguments]" -ForegroundColor Cyan
        Write-Host "Commands:"
        Write-Host "  up             Run all pending migrations"
        Write-Host "  down [count]   Roll back migrations (default count: 1)"
        Write-Host "  create [name]  Create a new sequential SQL migration pair"
        Write-Host "  force [ver]    Force a specific migration version state"
    }
}