settings {
    logfile = "/Users/dennisschroder/private/lsyncd.log",
    statusFile = "/Users/dennisschroder/private/lsyncd.status",
    nodaemon = false
}

sync {
    default.rsync,
    source = "/Users/dennisschroder/private/ha-config/",
    target = "root@homeassistant.local:/config/",
    delay = 1,
    rsync = {
        binary = "/opt/homebrew/bin/rsync",
        archive = true,
        compress = true,
        _extra = {"-e", "ssh -p 22 -i /Users/dennisschroder/.ssh/id_rsa_ha -o StrictHostKeyChecking=no"}
    }
}
