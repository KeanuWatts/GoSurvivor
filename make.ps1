param ($command)

$options = @("build", "test","")

if (!$options -contains $command) {
    Write-Host "Error: Use one of these options: $options"
} else {

    switch ($command) {
        ""{
            Write-Host "Go run main.go"
            go run main.go
        }
        "build" {
            go build -ldflags -H=windowsgui main.go
            Remove-Item GoSurvivor.exe
            Rename-Item main.exe GoSurvivor.exe
        }
        "test" {
            # Run test commands
        }
    }
}
